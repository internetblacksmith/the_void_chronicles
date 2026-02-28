// Copyright (C) 2024 Paolo Fabbri
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/keygen"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/getsentry/sentry-go"
	"github.com/muesli/termenv"
	"github.com/posthog/posthog-go"
)

var (
	// Build information (injected at build time via ldflags)
	buildTime = "unknown"
	gitCommit = "unknown"

	host          string
	httpPort      string
	httpsPort     string
	sshPort       string
	rateLimiter   *RateLimiter
	posthogClient posthog.Client
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := time.ParseDuration(value + "m"); err == nil {
			return int(intVal.Minutes())
		}
	}
	return defaultValue
}

func generateSSHKey(path string) error {
	// Use Charm's keygen to create proper SSH keys with write option
	_, err := keygen.New(path, keygen.WithKeyType(keygen.Ed25519), keygen.WithWrite())
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	log.Printf("Generated new SSH key at %s", path)
	return nil
}

// validatePassword checks if the provided password matches the required password
func validatePassword(password string) bool {
	requiredPassword := os.Getenv("SSH_PASSWORD")
	if requiredPassword == "" {
		log.Println("WARNING: SSH_PASSWORD not set, denying all connections")
		return false
	}
	return subtle.ConstantTimeCompare([]byte(password), []byte(requiredPassword)) == 1
}

// passwordHandler validates the password for SSH connections
func passwordHandler(ctx ssh.Context, password string) bool {
	addr := ctx.RemoteAddr()

	if !rateLimiter.AllowConnection(addr) {
		if posthogClient != nil {
			posthogClient.Enqueue(posthog.Capture{
				DistinctId: addr.String(),
				Event:      "ssh_rate_limited",
				Properties: posthog.NewProperties().
					Set("username", ctx.User()),
			})
		}
		return false
	}

	log.Printf("SSH authentication attempt from %s with user '%s'", addr, ctx.User())

	success := validatePassword(password)
	if success {
		log.Printf("SSH authentication successful for user '%s'", ctx.User())
		rateLimiter.RecordSuccessfulAuth(addr)

		if posthogClient != nil {
			posthogClient.Enqueue(posthog.Capture{
				DistinctId: ctx.User(),
				Event:      "ssh_login_success",
				Properties: posthog.NewProperties().
					Set("remote_addr", addr.String()),
			})
		}
	} else {
		log.Printf("SSH authentication failed for user '%s' (wrong password)", ctx.User())
		rateLimiter.RecordFailedAuth(addr)

		if posthogClient != nil {
			posthogClient.Enqueue(posthog.Capture{
				DistinctId: addr.String(),
				Event:      "ssh_login_failed",
				Properties: posthog.NewProperties().
					Set("username", ctx.User()),
			})
		}

		sentry.CaptureMessage("Failed SSH authentication attempt from " + addr.String())
	}

	return success
}

func main() {
	// Environment variables are provided by Doppler (via doppler run) or system
	// No .env file loading - use Doppler for all environments

	// Initialize Sentry for error tracking
	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			Environment:      getEnv("SENTRY_ENVIRONMENT", "production"),
			TracesSampleRate: 1.0,
		})
		if err != nil {
			log.Printf("Sentry initialization failed: %v", err)
		} else {
			log.Println("Sentry error tracking initialized")
			defer sentry.Flush(2 * time.Second)
		}
	} else {
		log.Println("SENTRY_DSN not set, skipping Sentry initialization")
	}

	// Initialize PostHog for analytics
	posthogKey := os.Getenv("POSTHOG_API_KEY")
	if posthogKey != "" {
		client, err := posthog.NewWithConfig(
			posthogKey,
			posthog.Config{
				Endpoint: getEnv("POSTHOG_HOST", "https://app.posthog.com"),
			},
		)
		if err != nil {
			log.Printf("PostHog initialization failed: %v", err)
			sentry.CaptureException(err)
		} else {
			posthogClient = client
			log.Println("PostHog analytics initialized")
			defer posthogClient.Close()
		}
	} else {
		log.Println("POSTHOG_API_KEY not set, skipping PostHog initialization")
	}

	// Configure ports
	host = getEnv("SSH_HOST", "0.0.0.0")

	// HTTP port: use HTTP_PORT env var or default
	httpPort = getEnv("HTTP_PORT", "8080")

	// HTTPS port: use HTTPS_PORT env var or default
	httpsPort = getEnv("HTTPS_PORT", "8443")

	// SSH port: use SSH_PORT env var or default
	sshPort = getEnv("SSH_PORT", "2222")

	// Initialize rate limiter
	rateLimiter = NewRateLimiter(5, 5*time.Minute, 15*time.Minute)

	// Log startup configuration
	log.Printf("=== Void Reader Starting ===")
	log.Printf("HTTP Port: %s", httpPort)
	log.Printf("HTTPS Port: %s", httpsPort)
	log.Printf("SSH Port: %s", sshPort)
	log.Printf("SSH Host: %s", host)
	log.Printf("Rate limiting: 5 attempts per 5 minutes, 15 minute block")
	log.Printf("Session timeout: %d minutes", getEnvInt("SSH_SESSION_TIMEOUT", 30))

	// Track application startup
	if posthogClient != nil {
		posthogClient.Enqueue(posthog.Capture{
			DistinctId: "system",
			Event:      "app_started",
			Properties: posthog.NewProperties().
				Set("environment", getEnv("SENTRY_ENVIRONMENT", "production")).
				Set("http_port", httpPort).
				Set("ssh_port", sshPort),
		})
	}

	// Ensure SSH key exists - use persistent volume in production
	var sshKeyPath string
	if _, err := os.Stat("/data"); err == nil {
		// Production: use persistent volume
		sshKeyPath = "/data/ssh/id_ed25519"
		os.MkdirAll("/data/ssh", 0700)
	} else {
		// Development: use local directory
		sshKeyPath = ".ssh/id_ed25519"
		os.MkdirAll(".ssh", 0700)
	}

	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		log.Println("SSH key not found, generating new key...")
		if err := generateSSHKey(sshKeyPath); err != nil {
			log.Fatalf("Failed to generate SSH key: %v", err)
		}
		log.Printf("Generated new SSH host key at %s", sshKeyPath)
	} else {
		log.Printf("Using existing SSH host key at %s", sshKeyPath)
	}

	// Check for port conflicts
	if httpPort == sshPort || httpsPort == sshPort || httpPort == httpsPort {
		log.Fatalf("ERROR: Port conflict! Ports must be unique. HTTP: %s, HTTPS: %s, SSH: %s", httpPort, httpsPort, sshPort)
	}

	// Start both HTTP and HTTPS servers
	go startHTTPServer()
	go startHTTPSServer()

	// Configure SSH server with more detailed logging
	wishMiddleware := []wish.Middleware{
		func(h ssh.Handler) ssh.Handler {
			return func(s ssh.Session) {
				sessionTimeout := time.Duration(getEnvInt("SSH_SESSION_TIMEOUT", 30)) * time.Minute
				log.Printf("SSH connection established from %s (user: %s)", s.RemoteAddr(), s.User())

				if posthogClient != nil {
					posthogClient.Enqueue(posthog.Capture{
						DistinctId: s.User(),
						Event:      "ssh_session_started",
						Properties: posthog.NewProperties().
							Set("remote_addr", s.RemoteAddr().String()),
					})
				}

				done := make(chan struct{})
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("SSH session panic recovered: %v", r)
							sentry.CaptureException(fmt.Errorf("SSH session panic: %v", r))
						}
						log.Printf("SSH connection closed for %s", s.RemoteAddr())

						if posthogClient != nil {
							posthogClient.Enqueue(posthog.Capture{
								DistinctId: s.User(),
								Event:      "ssh_session_ended",
								Properties: posthog.NewProperties().
									Set("remote_addr", s.RemoteAddr().String()),
							})
						}
						close(done)
					}()
					h(s)
				}()

				select {
				case <-done:
				case <-time.After(sessionTimeout):
					log.Printf("Session timeout for %s after %v", s.RemoteAddr(), sessionTimeout)

					if posthogClient != nil {
						posthogClient.Enqueue(posthog.Capture{
							DistinctId: s.User(),
							Event:      "ssh_session_timeout",
							Properties: posthog.NewProperties().
								Set("remote_addr", s.RemoteAddr().String()).
								Set("timeout_minutes", sessionTimeout.Minutes()),
						})
					}

					s.Close()
					<-done
				}
			}
		},
		logging.Middleware(),
		bubbletea.MiddlewareWithColorProfile(teaHandler, termenv.TrueColor),
	}

	serverOpts := []ssh.Option{
		wish.WithAddress(net.JoinHostPort(host, sshPort)),
		wish.WithHostKeyPath(sshKeyPath),
		wish.WithMiddleware(wishMiddleware...),
	}

	requirePassword := getEnv("SSH_REQUIRE_PASSWORD", "true")
	if strings.ToLower(requirePassword) == "true" || requirePassword == "1" {
		serverOpts = append(serverOpts, wish.WithPasswordAuth(passwordHandler))
		log.Println("Password authentication: ENABLED")
	} else {
		// When password is not required, accept both public key auth and password auth
		// This allows clients to connect with any key or even with a dummy password
		serverOpts = append(serverOpts,
			wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
				return true // Accept any public key
			}),
			wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
				return true // Accept any password
			}),
		)
		log.Println("Password authentication: DISABLED (allowing all connections)")
	}

	s, err := wish.NewServer(serverOpts...)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Log server configuration
	log.Printf("HTTP server listening on 0.0.0.0:%s", httpPort)
	log.Printf("HTTPS server listening on 0.0.0.0:%s (if certs available)", httpsPort)
	log.Printf("SSH server listening on %s:%s", host, sshPort)
	// Don't log the actual password for security
	if os.Getenv("SSH_PASSWORD") != "" {
		log.Printf("SSH Password: [configured via SSH_PASSWORD env var]")
	} else {
		log.Printf("SSH Password: [using default]")
	}

	// Start SSH server
	go func() {
		log.Println("Starting SSH server...")
		if err = s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Printf("SSH server error: %v", err)
			sentry.CaptureException(err)
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping servers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && err != ssh.ErrServerClosed {
		log.Fatalln(err)
	}
}

func startHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Bob's Personal Homepage</title>
    <style>
        body {
            background: #C0C0C0;
            font-family: "Times New Roman", serif;
            margin: 20px;
        }
        h1 {
            color: #000080;
            font-size: 36px;
            text-align: center;
        }
        .counter {
            text-align: center;
            margin: 20px;
        }
        .links {
            background: #FFFF00;
            border: 3px ridge #808080;
            padding: 10px;
            margin: 20px auto;
            width: 500px;
        }
        a {
            color: #0000FF;
            text-decoration: underline;
        }
        .construction {
            text-align: center;
            color: #FF0000;
            font-size: 18px;
            blink: true;
        }
        hr {
            border: none;
            border-top: 3px double #333;
            color: #333;
            overflow: visible;
            text-align: center;
            height: 5px;
        }
    </style>
</head>
<body>
    <h1>Welcome to Bob's Homepage!</h1>
    <hr>
    <marquee behavior="scroll" direction="left">🚧 Under Construction Since 1997! 🚧</marquee>
    
    <center>
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=" width="88" height="31" alt="Best viewed in Netscape Navigator">
    </center>
    
    <div class="construction">
        <p>🚧 This site is UNDER CONSTRUCTION 🚧</p>
        <p>Last updated: January 15, 1998</p>
    </div>
    
    <div class="links">
        <h2>My Favorite Links:</h2>
        <ul>
            <li><a href="#">My Resume</a> (Coming Soon!)</li>
            <li><a href="#">Pictures of my Cat</a> (Under Construction)</li>
            <li><a href="#">Cool MIDI Files</a> (Broken Link)</li>
            <li><a href="#">Guestbook</a> (Please Sign!)</li>
            <li><a href="#">WebRing</a> (Join my WebRing!)</li>
        </ul>
    </div>
    
    <center>
        <p>You are visitor number:</p>
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=" alt="00001337">
        <br><br>
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=" width="88" height="31" alt="Made with Notepad">
        <br>
        <p>Best viewed at 800x600 resolution</p>
        <p>© 1997-1998 Bob Smith. All rights reserved.</p>
    </center>
    
    <hr>
    <center>
        <p><i>Email me at: webmaster@bobshomepage.geocities.com</i></p>
    </center>
</body>
</html>
`

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		version := os.Getenv("KAMAL_VERSION")
		if version == "" {
			version = gitCommit
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	})

	// Storage monitoring endpoints
	pm := NewProgressManager()
	http.HandleFunc("/api/storage/stats", pm.StorageStatsHandler)
	http.HandleFunc("/api/storage/cleanup", pm.CleanupHandler)

	// Start cleanup scheduler
	pm.StartCleanupScheduler()

	log.Printf("Starting HTTP server on 0.0.0.0:%s", httpPort)
	if err := http.ListenAndServe("0.0.0.0:"+httpPort, nil); err != nil {
		log.Fatalf("FATAL: HTTP server failed to start: %v", err)
	}
}

func startHTTPSServer() {
	certFile := getEnv("TLS_CERT_PATH", "/data/ssl/cert.pem")
	keyFile := getEnv("TLS_KEY_PATH", "/data/ssl/key.pem")

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Printf("TLS certificate not found at %s, HTTPS server not started", certFile)
		return
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Printf("TLS key not found at %s, HTTPS server not started", keyFile)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Bob's Personal Homepage</title>
    <style>
        body {
            background: #C0C0C0;
            font-family: "Times New Roman", serif;
            margin: 20px;
        }
        h1 {
            color: #000080;
            font-size: 36px;
            text-align: center;
        }
        .counter {
            text-align: center;
            margin: 20px;
        }
        .links {
            background: #FFFF00;
            border: 3px ridge #808080;
            padding: 10px;
            margin: 20px auto;
            width: 500px;
        }
        a {
            color: #0000FF;
            text-decoration: underline;
        }
        .construction {
            text-align: center;
            color: #FF0000;
            font-size: 18px;
            blink: true;
        }
        hr {
            border: none;
            border-top: 3px double #333;
            color: #333;
            overflow: visible;
            text-align: center;
            height: 5px;
        }
    </style>
</head>
<body>
    <h1>Welcome to Bob's Homepage!</h1>
    <hr>
    <marquee behavior="scroll" direction="left">🚧 Under Construction Since 1997! 🚧</marquee>
    
    <center>
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=" width="88" height="31" alt="Best viewed in Netscape Navigator">
    </center>
    
    <div class="construction">
        <p>🚧 This site is UNDER CONSTRUCTION 🚧</p>
        <p>Last updated: January 15, 1998</p>
    </div>
    
    <div class="links">
        <h2>My Favorite Links:</h2>
        <ul>
            <li><a href="#">My Resume</a> (Coming Soon!)</li>
            <li><a href="#">Pictures of my Cat</a> (Under Construction)</li>
            <li><a href="#">Cool MIDI Files</a> (Broken Link)</li>
            <li><a href="#">Guestbook</a> (Please Sign!)</li>
            <li><a href="#">WebRing</a> (Join my WebRing!)</li>
        </ul>
    </div>
    
    <center>
        <p>You are visitor number:</p>
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=" alt="00001337">
        <br><br>
        <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=" width="88" height="31" alt="Made with Notepad">
        <br>
        <p>Best viewed at 800x600 resolution</p>
        <p>© 1997-1998 Bob Smith. All rights reserved.</p>
    </center>
    
    <hr>
    <center>
        <p><i>Email me at: webmaster@bobshomepage.geocities.com</i></p>
    </center>
</body>
</html>
`

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		version := os.Getenv("KAMAL_VERSION")
		if version == "" {
			version = gitCommit
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	})

	pm := NewProgressManager()
	mux.HandleFunc("/api/storage/stats", pm.StorageStatsHandler)
	mux.HandleFunc("/api/storage/cleanup", pm.CleanupHandler)

	pm.StartCleanupScheduler()

	log.Printf("Starting HTTPS server on 0.0.0.0:%s", httpsPort)
	server := &http.Server{
		Addr:    "0.0.0.0:" + httpsPort,
		Handler: mux,
	}

	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Printf("WARNING: HTTPS server failed to start: %v", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		fmt.Println("no active terminal, skipping")
		return nil, nil
	}

	log.Printf("PTY detected: term=%s, width=%d, height=%d",
		pty.Term, pty.Window.Width, pty.Window.Height)

	// Get username from SSH session
	username := s.User()
	if username == "" {
		username = "reader"
	}

	m := initialModelWithUser(pty.Window.Width, pty.Window.Height, username)
	return m, bubbletea.MakeOptions(s)
}

type state int

const (
	menuView state = iota
	chapterListView
	readingView
	aboutView
	progressView
	licenseView
)

type BookInfo struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Status    string `json:"status"`
	Summary   string `json:"summary"`
	Available bool   `json:"available"`
}

type SeriesInfo struct {
	Series string     `json:"series"`
	Books  []BookInfo `json:"books"`
}

type model struct {
	state           state
	book            *Book
	menuCursor      int
	chapterCursor   int
	currentChapter  int
	scrollOffset    int
	width           int
	height          int
	menuItems       []string
	quitting        bool
	progress        *UserProgress
	progressManager *ProgressManager
	username        string
	books           []BookInfo
	selectedBook    int
	licenseScroll   int
}

func getSeriesBooks() []BookInfo {
	// Try to load series info from JSON file
	data, err := os.ReadFile("series.json")
	if err != nil {
		// If file doesn't exist, try alternate paths
		data, err = os.ReadFile("ssh-reader/series.json")
		if err != nil {
			data, err = os.ReadFile("../series.json")
			if err != nil {
				log.Printf("Warning: Could not load series.json: %v", err)
				// Return a minimal fallback
				return []BookInfo{{
					Number:    1,
					Title:     "Void Reavers",
					Subtitle:  "A Tale of Space Pirates and Cosmic Plunder",
					Status:    "✓ Available",
					Available: true,
					Summary:   "Captain Zara leads her pirate crew through the lawless void.",
				}}
			}
		}
	}

	var series SeriesInfo
	if err := json.Unmarshal(data, &series); err != nil {
		log.Printf("Error parsing series.json: %v", err)
		return []BookInfo{{
			Number:    1,
			Title:     "Void Reavers",
			Subtitle:  "A Tale of Space Pirates and Cosmic Plunder",
			Status:    "✓ Available",
			Available: true,
			Summary:   "Captain Zara leads her pirate crew through the lawless void.",
		}}
	}

	return series.Books
}

func initialModelWithUser(width, height int, username string) model {
	// Load the book content
	book, err := LoadBook("book1_void_reavers_source")
	if err != nil {
		// Try alternate path when running from ssh-reader directory
		book, err = LoadBook("../book1_void_reavers_source")
		if err != nil {
			log.Printf("Error loading book: %v", err)
			book = &Book{
				Title:    "Error Loading Book",
				Chapters: []Chapter{{Title: "Error", Content: fmt.Sprintf("Could not load book: %v", err)}},
			}
		}
	}

	pm := NewProgressManager()
	progress, err := pm.LoadProgress(username)
	if err != nil {
		log.Printf("Error loading progress for %s: %v", username, err)
		progress = &UserProgress{
			Username:        username,
			CurrentChapter:  0,
			ScrollOffset:    0,
			ChapterProgress: make(map[int]bool),
			Bookmarks:       []Bookmark{},
		}
	}

	books := getSeriesBooks()

	// Build menu items from books
	menuItems := []string{}
	for _, book := range books {
		menuItems = append(menuItems, fmt.Sprintf("📚 Book %d: %s", book.Number, book.Title))
	}
	menuItems = append(menuItems, "", "ℹ️  About", "📄 License", "", "🚪 Exit")

	return model{
		state:           menuView,
		book:            book,
		width:           width,
		height:          height,
		menuItems:       menuItems,
		progress:        progress,
		progressManager: pm,
		username:        username,
		currentChapter:  progress.CurrentChapter,
		scrollOffset:    progress.ScrollOffset,
		books:           books,
		selectedBook:    0, // Start with first book selected
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case menuView:
			return m.updateMenu(msg)
		case chapterListView:
			return m.updateChapterList(msg)
		case readingView:
			return m.updateReading(msg)
		case aboutView:
			return m.updateAbout(msg)
		case progressView:
			return m.updateProgress(msg)
		case licenseView:
			return m.updateLicense(msg)
		}
	}

	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	case "up", "k":
		if m.menuCursor > 0 {
			m.menuCursor--
			// Skip separators
			if m.menuItems[m.menuCursor] == "" && m.menuCursor > 0 {
				m.menuCursor--
			}
			// Update selected book when cursor moves
			if m.menuCursor < len(m.books) {
				m.selectedBook = m.menuCursor
			}
		}
	case "down", "j":
		if m.menuCursor < len(m.menuItems)-1 {
			m.menuCursor++
			// Skip separators
			if m.menuItems[m.menuCursor] == "" && m.menuCursor < len(m.menuItems)-1 {
				m.menuCursor++
			}
			// Update selected book when cursor moves
			if m.menuCursor < len(m.books) {
				m.selectedBook = m.menuCursor
			}
		}
	case "enter", " ", "r":
		item := m.menuItems[m.menuCursor]
		if item == "🚪 Exit" {
			// Save progress before quitting
			m.progress.CurrentChapter = m.currentChapter
			m.progress.ScrollOffset = m.scrollOffset
			m.progressManager.SaveProgress(m.progress)
			m.quitting = true
			return m, tea.Quit
		} else if item == "ℹ️  About" {
			m.state = aboutView
		} else if item == "📄 License" {
			m.licenseScroll = 0
			m.state = licenseView
		} else if m.menuCursor < len(m.books) && m.books[m.menuCursor].Available {
			// Start reading the selected book
			m.state = chapterListView
			m.chapterCursor = 0
		}
	case "c":
		// Continue reading from saved position
		if m.menuCursor < len(m.books) && m.books[m.menuCursor].Available && m.progress.CurrentChapter > 0 {
			m.state = readingView
			m.currentChapter = m.progress.CurrentChapter
			m.scrollOffset = m.progress.ScrollOffset
		}
	}
	return m, nil
}

func (m model) updateChapterList(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.state = menuView
	case "up", "k":
		if m.chapterCursor > 0 {
			m.chapterCursor--
		}
	case "down", "j":
		if m.chapterCursor < len(m.book.Chapters)-1 {
			m.chapterCursor++
		}
	case "enter", " ":
		m.state = readingView
		m.currentChapter = m.chapterCursor
		m.scrollOffset = 0
	}
	return m, nil
}

func (m model) updateReading(msg tea.KeyMsg) (model, tea.Cmd) {
	contentHeight := m.height - 4 // Leave space for header and footer

	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.state = menuView
	case "up", "k":
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
	case "down", "j":
		maxScroll := m.getMaxScroll()
		if m.scrollOffset < maxScroll {
			m.scrollOffset++
		}
	case "page up", "ctrl+b":
		m.scrollOffset -= contentHeight
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	case "page down", "ctrl+f", " ":
		maxScroll := m.getMaxScroll()
		m.scrollOffset += contentHeight
		if m.scrollOffset > maxScroll {
			m.scrollOffset = maxScroll
		}
	case "home", "g":
		m.scrollOffset = 0
	case "end", "G":
		m.scrollOffset = m.getMaxScroll()
	case "left", "h", "p":
		if m.currentChapter > 0 {
			// Save progress before changing chapter
			m.progress.CurrentChapter = m.currentChapter
			m.progress.ScrollOffset = m.scrollOffset
			m.progressManager.SaveProgress(m.progress)

			m.currentChapter--
			m.scrollOffset = 0
		}
	case "right", "l", "n":
		if m.currentChapter < len(m.book.Chapters)-1 {
			// Mark current chapter as completed and save progress
			m.progress.MarkChapterComplete(m.currentChapter)
			m.progress.CurrentChapter = m.currentChapter + 1
			m.progress.ScrollOffset = 0
			m.progressManager.SaveProgress(m.progress)

			m.currentChapter++
			m.scrollOffset = 0
		}
	case "b":
		// Add bookmark
		m.progress.AddBookmark(m.currentChapter, m.scrollOffset, "")
		m.progressManager.SaveProgress(m.progress)
	}
	return m, nil
}

func (m model) updateAbout(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "enter", " ":
		m.state = menuView
	}
	return m, nil
}

func (m model) updateProgress(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc", "enter", " ":
		m.state = menuView
	}
	return m, nil
}

func (m model) updateLicense(msg tea.KeyMsg) (model, tea.Cmd) {
	contentHeight := m.height - 4

	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.state = menuView
		m.licenseScroll = 0
	case "up", "k":
		if m.licenseScroll > 0 {
			m.licenseScroll--
		}
	case "down", "j":
		maxScroll := m.getLicenseMaxScroll()
		if m.licenseScroll < maxScroll {
			m.licenseScroll++
		}
	case "page up", "ctrl+b":
		m.licenseScroll -= contentHeight
		if m.licenseScroll < 0 {
			m.licenseScroll = 0
		}
	case "page down", "ctrl+f", " ":
		maxScroll := m.getLicenseMaxScroll()
		m.licenseScroll += contentHeight
		if m.licenseScroll > maxScroll {
			m.licenseScroll = maxScroll
		}
	case "home", "g":
		m.licenseScroll = 0
	case "end", "G":
		m.licenseScroll = m.getLicenseMaxScroll()
	}
	return m, nil
}

func (m model) getMaxScroll() int {
	chapter := m.book.Chapters[m.currentChapter]
	lines := len(wrapText(chapter.Content, m.width-4))
	contentHeight := m.height - 4
	maxScroll := lines - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

func (m model) getLicenseText() string {
	return `📄 BOOK CONTENT LICENSE

The Void Chronicles book series is licensed under:

Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📖 YOU ARE FREE TO:

• Share — copy and redistribute the material in any medium or format
• Adapt — remix, transform, and build upon the material

The licensor cannot revoke these freedoms as long as you follow the license terms.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 UNDER THE FOLLOWING TERMS:

Attribution (BY):
You must give appropriate credit, provide a link to the license, and indicate if changes were made. You may do so in any reasonable manner, but not in any way that suggests the licensor endorses you or your use.

NonCommercial (NC):
You may not use the material for commercial purposes. Commercial use requires separate permission from the author.

ShareAlike (SA):
If you remix, transform, or build upon the material, you must distribute your contributions under the same license as the original.

No additional restrictions:
You may not apply legal terms or technological measures that legally restrict others from doing anything the license permits.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

ℹ️  NOTICES:

You do not have to comply with the license for elements of the material in the public domain or where your use is permitted by an applicable exception or limitation.

No warranties are given. The license may not give you all of the permissions necessary for your intended use. For example, other rights such as publicity, privacy, or moral rights may limit how you use the material.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💻 SSH READER LICENSE

The SSH reader application is licensed under:

GNU Affero General Public License v3.0 (AGPL-3.0)

This is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

Source code: github.com/internetblacksmith/void-chronicles

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔗 FULL LICENSE TEXTS:

CC BY-NC-SA 4.0: https://creativecommons.org/licenses/by-nc-sa/4.0/
AGPL-3.0: https://www.gnu.org/licenses/agpl-3.0.html

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📧 QUESTIONS OR COMMERCIAL USE?

For commercial licensing inquiries or questions about the license, please contact the author.`
}

func (m model) getLicenseMaxScroll() int {
	licenseText := m.getLicenseText()
	lines := len(wrapText(licenseText, m.width-4))
	contentHeight := m.height - 4
	maxScroll := lines - contentHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	switch m.state {
	case menuView:
		return m.viewMenu()
	case chapterListView:
		return m.viewChapterList()
	case readingView:
		return m.viewReading()
	case aboutView:
		return m.viewAbout()
	case progressView:
		return m.viewProgress()
	case licenseView:
		return m.viewLicense()
	}
	return ""
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Align(lipgloss.Center)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("86")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	contentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center)
)

func (m model) viewMenu() string {
	// Title bar
	title := titleStyle.Width(m.width).Render("📚 THE VOID CHRONICLES 📚")
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Width(m.width).
		Render("An AI-Generated Space Opera Series")

	// Calculate split widths
	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth - 3 // Account for borders

	// Left panel - Book list
	leftPanelStyle := lipgloss.NewStyle().
		MaxWidth(leftWidth).
		Height(m.height-8).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)

	var menuItems []string
	maxMenuWidth := leftWidth - 6
	menuItems = append(menuItems, lipgloss.NewStyle().Bold(true).MaxWidth(maxMenuWidth).Render("BOOK LIBRARY"))
	menuItems = append(menuItems, "")

	maxTitleLen := maxMenuWidth - 12
	for i, book := range m.books {
		title := book.Title
		if len(title) > maxTitleLen {
			title = title[:maxTitleLen-1] + "…"
		}
		item := fmt.Sprintf("Book %d: %s", book.Number, title)
		if i == m.menuCursor {
			if book.Available {
				menuItems = append(menuItems, selectedStyle.MaxWidth(maxMenuWidth).Render("▶ "+item+" ✓"))
			} else {
				menuItems = append(menuItems, selectedStyle.MaxWidth(maxMenuWidth).Render("▶ "+item))
			}
		} else {
			if book.Available {
				menuItems = append(menuItems, normalStyle.MaxWidth(maxMenuWidth).Render("  "+item+" ✓"))
			} else {
				statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MaxWidth(maxMenuWidth)
				menuItems = append(menuItems, statusStyle.Render("  "+item))
			}
		}
	}

	// Add About, License, and Exit options
	menuItems = append(menuItems, "")
	aboutIndex := len(m.books)
	if m.menuCursor == aboutIndex+1 {
		menuItems = append(menuItems, selectedStyle.MaxWidth(maxMenuWidth).Render("▶ ℹ️  About"))
	} else {
		menuItems = append(menuItems, normalStyle.MaxWidth(maxMenuWidth).Render("  ℹ️  About"))
	}

	licenseIndex := aboutIndex + 2
	if m.menuCursor == licenseIndex {
		menuItems = append(menuItems, selectedStyle.MaxWidth(maxMenuWidth).Render("▶ 📄 License"))
	} else {
		menuItems = append(menuItems, normalStyle.MaxWidth(maxMenuWidth).Render("  📄 License"))
	}

	menuItems = append(menuItems, "")
	if m.menuCursor == len(m.menuItems)-1 {
		menuItems = append(menuItems, selectedStyle.MaxWidth(maxMenuWidth).Render("▶ 🚪 Exit"))
	} else {
		menuItems = append(menuItems, normalStyle.MaxWidth(maxMenuWidth).Render("  🚪 Exit"))
	}

	leftPanel := leftPanelStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, menuItems...),
	)

	// Right panel - Book details
	rightPanelStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height-8).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)

	var rightContent string
	if m.menuCursor < len(m.books) {
		book := m.books[m.menuCursor]

		// Book title
		bookTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Width(rightWidth - 6).
			Render(fmt.Sprintf("BOOK %d: %s", book.Number, strings.ToUpper(book.Title)))

		// Subtitle
		bookSubtitle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(rightWidth - 6).
			Render(book.Subtitle)

		// Status
		var statusText string
		statusStyle := lipgloss.NewStyle().Width(rightWidth - 6)
		if book.Available {
			statusStyle = statusStyle.Foreground(lipgloss.Color("82"))
			statusText = "✓ Available to Read"
		} else {
			statusStyle = statusStyle.Foreground(lipgloss.Color("214"))
			statusText = fmt.Sprintf("📅 %s", book.Status)
		}
		status := statusStyle.Render(statusText)

		// Summary
		summaryTitle := lipgloss.NewStyle().
			Bold(true).
			Margin(1, 0, 0, 0).
			Render("Synopsis:")

		summaryText := lipgloss.NewStyle().
			Width(rightWidth - 6).
			Render(book.Summary)

		// Options
		var options string
		if book.Available {
			optionsStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Margin(1, 0, 0, 0)

			optionsList := []string{
				"[Enter] Start Reading",
			}
			if m.progress.CurrentChapter > 0 {
				optionsList = append(optionsList, fmt.Sprintf("[C] Continue Chapter %d", m.progress.CurrentChapter+1))
			}
			options = optionsStyle.Render(strings.Join(optionsList, "\n"))
		}

		rightContent = lipgloss.JoinVertical(
			lipgloss.Left,
			bookTitle,
			bookSubtitle,
			"",
			status,
			"",
			summaryTitle,
			summaryText,
			"",
			options,
		)
	} else if m.menuItems[m.menuCursor] == "ℹ️  About" {
		// About selected
		rightContent = lipgloss.NewStyle().
			Width(rightWidth - 6).
			Render(`ℹ️  ABOUT THE VOID CHRONICLES

This is an experimental SSH-based book reader for reading science fiction novels directly in your terminal.

🚀 The Project:
An open-source reading platform that combines classic terminal aesthetics with modern TUI frameworks.

🌌 The Series:
The Void Chronicles is a 10-book epic following humanity's evolution from chaotic space pirates to cosmic gardeners.

✨ Built With:
• Go programming language
• Charm's Bubbletea TUI framework
• Wish SSH server library
• Lipgloss for styling

🔧 Features:
• Read books over SSH
• Progress tracking
• Bookmarks
• Chapter navigation
• No installation required

📡 Connect: ssh vc.internetblacksmith.dev
🐙 Source: github.com/internetblacksmith

[Enter] Learn More`)
	} else if m.menuItems[m.menuCursor] == "📄 License" {
		// License selected
		rightContent = lipgloss.NewStyle().
			Width(rightWidth - 6).
			Render(`📄 LICENSE INFORMATION

The Void Chronicles books are released under:

📖 Creative Commons BY-NC-SA 4.0

This means you can:
• Share and remix the books
• Create derivative works
• Distribute freely

As long as you:
• Give appropriate credit
• Don't use commercially
• Share derivatives under same license

The SSH reader application is:

💻 GNU AGPL-3.0 (open source)

[Enter] View Full License`)
	} else {
		// Exit selected
		rightContent = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(rightWidth - 6).
			Height(m.height - 12).
			Render("\n\n\n🚪 Exit SSH Reader\n\nThank you for exploring\nThe Void Chronicles!")
	}

	rightPanel := rightPanelStyle.Render(rightContent)

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		"  ",
		rightPanel,
	)

	// Footer
	footer := footerStyle.Width(m.width).Render("↑/↓: navigate • enter: read • c: continue • q: quit")

	// Combine all elements
	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		"",
		panels,
		"",
		footer,
	)
}

func (m model) viewChapterList() string {
	contentWidth := m.width - 6
	contentHeight := m.height - 8

	header := headerStyle.Width(contentWidth - 2).Render("📚 CHAPTERS")

	maxItemWidth := contentWidth - 6
	var items []string
	for i, chapter := range m.book.Chapters {
		prefix := fmt.Sprintf("%2d. ", i+1)
		if i == m.chapterCursor {
			items = append(items, selectedStyle.MaxWidth(maxItemWidth).Render("▶ "+prefix+chapter.Title))
		} else {
			items = append(items, normalStyle.MaxWidth(maxItemWidth).Render("  "+prefix+chapter.Title))
		}
	}

	chapterStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)

	content := chapterStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))

	footer := footerStyle.Width(m.width).Render("↑/↓: navigate • enter: read • esc: back • q: quit")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		content,
		"",
		footer,
	)
}

func (m model) viewReading() string {
	chapter := m.book.Chapters[m.currentChapter]
	progress := fmt.Sprintf("Chapter %d of %d", m.currentChapter+1, len(m.book.Chapters))

	contentWidth := m.width - 6
	contentHeight := m.height - 8

	header := headerStyle.Width(contentWidth - 2).Render(fmt.Sprintf("📖 %s (%s)", chapter.Title, progress))

	lines := wrapText(chapter.Content, contentWidth-6)

	visibleHeight := contentHeight - 4
	startLine := m.scrollOffset
	endLine := startLine + visibleHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	var visibleLines []string
	if startLine < len(lines) {
		visibleLines = lines[startLine:endLine]
	}

	readingStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)

	content := readingStyle.Render(lipgloss.JoinVertical(lipgloss.Left, visibleLines...))

	navHelp := "h/←: prev chapter • l/→: next chapter • ↑/↓: scroll • b: bookmark • space: page down • esc: menu"
	footer := footerStyle.Width(m.width).Render(navHelp)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		content,
		footer,
	)
}

func (m model) viewAbout() string {
	contentWidth := m.width - 6
	contentHeight := m.height - 8

	header := headerStyle.Width(contentWidth - 2).Render("ℹ️  ABOUT THE VOID CHRONICLES")

	aboutText := `🚀 Welcome to the Void Chronicles Universe! 🚀

Void Reavers is the first book in an epic 10-book science fiction series exploring 
humanity's evolution from chaotic space pirates to cosmic gardeners.

📖 The Story:
Follow Captain Zara "Bloodhawk" Vega's fifty-year journey as she transforms from 
a young pirate forced into Rex Morrison's brutal crew to humanity's ambassador to 
alien civilizations.

🌌 The Universe:
Set in a cosmos where quantum physics can tear reality apart and ancient alien 
Architects judge humanity's every move, pirates must evolve from raiders to 
protectors.

✨ Themes:
• Personal transformation mirrors species evolution
• The balance between order and chaos
• Earning cosmic citizenship through wisdom
• Honor among thieves in the vastness of space

📡 The SSH Reader:
This terminal-based reading experience is an experiment in making books accessible 
through SSH. No installation required - just connect and read!

Features:
• Progress tracking across sessions
• Bookmarks for favorite passages
• Chapter navigation
• Cross-platform (works anywhere SSH works)

🔧 Technical Details:
• Built with Go, Bubbletea, and Wish
• Open source (AGPL-3.0)
• Deployed with Kamal
• Secrets managed with Doppler

📦 Deployment Info:
• Build Time: ` + buildTime + `
• Git Commit: ` + gitCommit + `

🎭 Author: Captain J. Starwind
📅 Series: The Void Chronicles (Book 1 of 10)
🐙 Source: github.com/internetblacksmith/void-chronicles
📡 Connect: ssh vc.internetblacksmith.dev`

	wrappedLines := wrapText(aboutText, contentWidth-6)
	wrappedText := strings.Join(wrappedLines, "\n")

	aboutStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)

	content := aboutStyle.Render(wrappedText)
	footer := footerStyle.Width(m.width).Render("press any key to return to menu")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		content,
		"",
		footer,
	)
}

func (m model) viewProgress() string {
	contentWidth := m.width - 6
	contentHeight := m.height - 8

	header := headerStyle.Width(contentWidth - 2).Render("📊 READING PROGRESS")

	completion := m.progress.GetCompletionPercentage(len(m.book.Chapters))

	progressText := fmt.Sprintf(`👤 Reader: %s

📈 Overall Progress: %.1f%% Complete
📖 Current Chapter: %d of %d
⏱️  Total Reading Time: %v
📅 Last Read: %s
🔖 Bookmarks: %d

📚 Chapter Progress:`,
		m.progress.Username,
		completion,
		m.progress.CurrentChapter+1,
		len(m.book.Chapters),
		m.progress.ReadingTime.Truncate(time.Minute),
		m.progress.LastRead.Format("Jan 2, 2006 15:04"),
		len(m.progress.Bookmarks),
	)

	// Add chapter completion status
	for i, chapter := range m.book.Chapters {
		status := "⭕"
		if m.progress.IsChapterComplete(i) {
			status = "✅"
		} else if i == m.progress.CurrentChapter {
			status = "📍"
		}
		progressText += fmt.Sprintf("\n  %s Chapter %d: %s", status, i+1, chapter.Title)
	}

	if len(m.progress.Bookmarks) > 0 {
		progressText += "\n\n🔖 Recent Bookmarks:"
		for i, bookmark := range m.progress.Bookmarks {
			if i >= 5 { // Show only first 5 bookmarks
				break
			}
			progressText += fmt.Sprintf("\n  • Chapter %d - %s",
				bookmark.Chapter+1,
				bookmark.Created.Format("Jan 2 15:04"))
		}
	}

	wrappedLines := wrapText(progressText, contentWidth-6)
	wrappedProgress := strings.Join(wrappedLines, "\n")

	progressStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)

	content := progressStyle.Render(wrappedProgress)
	footer := footerStyle.Width(m.width).Render("press any key to return to menu")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		content,
		"",
		footer,
	)
}

func (m model) viewLicense() string {
	contentWidth := m.width - 6
	contentHeight := m.height - 8

	header := headerStyle.Width(contentWidth - 2).Render("📄 LICENSE INFORMATION")

	licenseText := m.getLicenseText()
	wrappedLines := wrapText(licenseText, contentWidth-6)

	visibleHeight := contentHeight - 4

	visibleLines := wrappedLines
	if len(wrappedLines) > visibleHeight {
		end := m.licenseScroll + visibleHeight
		if end > len(wrappedLines) {
			end = len(wrappedLines)
		}
		visibleLines = wrappedLines[m.licenseScroll:end]
	}

	licenseStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2)

	content := licenseStyle.Render(strings.Join(visibleLines, "\n"))

	scrollInfo := ""
	if len(wrappedLines) > visibleHeight {
		scrollInfo = fmt.Sprintf(" (line %d/%d)", m.licenseScroll+1, len(wrappedLines)-visibleHeight+1)
	}

	footer := footerStyle.Width(m.width).Render("↑/↓: scroll • space/pgdn: page down • pgup: page up • esc: back" + scrollInfo)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		content,
		"",
		footer,
	)
}
