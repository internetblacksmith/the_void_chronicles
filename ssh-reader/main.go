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
)

var (
	host     string
	httpPort string
	sshPort  string
)

func init() {
	host = getEnv("SSH_HOST", "0.0.0.0")
	
	// Get all port-related environment variables
	portEnv := os.Getenv("PORT")
	httpPortEnv := os.Getenv("HTTP_PORT") 
	sshPortEnv := os.Getenv("SSH_PORT")
	railwayTcpPort := os.Getenv("RAILWAY_TCP_APPLICATION_PORT")
	
	// Railway dual port configuration:
	// - PORT: for HTTP domain (web traffic)
	// - RAILWAY_TCP_APPLICATION_PORT: for TCP proxy (SSH traffic)
	// These MUST be different ports for Railway to work correctly
	
	if portEnv != "" {
		httpPort = portEnv
		log.Printf("Using Railway-provided PORT for HTTP: %s", httpPort)
	} else {
		httpPort = "8080"
		log.Printf("No Railway PORT found, using 8080 for local development")
	}
	
	// SSH always uses port 2222
	// Railway's TCP proxy handles the external routing to a different URL
	sshPort = "2222"
	if railwayTcpPort != "" {
		log.Printf("Railway TCP proxy configured for port %s, but SSH stays on 2222", railwayTcpPort)
	} else {
		log.Printf("SSH server using standard port: 2222")
	}
	
	log.Printf("Port resolution: PORT=%s, HTTP_PORT=%s, SSH_PORT=%s, RAILWAY_TCP=%s -> Using HTTP=%s, SSH=%s",
		portEnv, httpPortEnv, sshPortEnv, railwayTcpPort, httpPort, sshPort)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
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

// passwordHandler validates the password for SSH connections
func passwordHandler(ctx ssh.Context, password string) bool {
	// Get password from environment variable, or use default
	requiredPassword := os.Getenv("SSH_PASSWORD")
	if requiredPassword == "" {
		requiredPassword = "Amigos4Life!"
	}
	
	// Check if the password matches
	return password == requiredPassword
}

func main() {
	// Log final port configuration
	log.Printf("Starting servers - HTTP on 0.0.0.0:%s, SSH on %s:%s", httpPort, host, sshPort)
	
	// Ensure SSH key exists
	sshKeyPath := "../.ssh/id_ed25519"
	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		log.Println("SSH key not found, generating new key...")
		os.MkdirAll("../.ssh", 0700)
		if err := generateSSHKey(sshKeyPath); err != nil {
			log.Fatalf("Failed to generate SSH key: %v", err)
		}
	}
	
	// Handle port conflicts - Railway will provide different external URLs
	if httpPort == sshPort {
		log.Printf("Port conflict detected: both HTTP and SSH want port %s", httpPort)
		if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
			log.Printf("On Railway: HTTP gets domain access, SSH gets TCP proxy URL - no conflict")
		} else {
			// Local development conflict - adjust SSH port
			sshPort = "2223"
			log.Printf("Local development: moving SSH to port 2223 to avoid conflict")
		}
	}
	
	// Start both HTTP and SSH servers
	go startHTTPServer()
	
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, sshPort)),
		wish.WithHostKeyPath(sshKeyPath),
		wish.WithPasswordAuth(passwordHandler),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("HTTP server listening on port %s", httpPort)
	log.Printf("SSH server listening on internal port %s", sshPort)
	log.Printf("SSH Password: Amigos4Life!")
	log.Printf("Note: For Railway, configure TCP proxy to forward to port %s", sshPort)
	go func() {
		if err = s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
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
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	
	log.Printf("Starting HTTP server on 0.0.0.0:%s", httpPort)
	if err := http.ListenAndServe("0.0.0.0:"+httpPort, nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}


func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		fmt.Println("no active terminal, skipping")
		return nil, nil
	}
	
	// Get username from SSH session
	username := s.User()
	if username == "" {
		username = "reader"
	}
	
	m := initialModelWithUser(pty.Window.Width, pty.Window.Height, username)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type state int

const (
	menuView state = iota
	chapterListView
	readingView
	aboutView
	progressView
)

type BookInfo struct {
	Number    int
	Title     string
	Subtitle  string
	Status    string
	Summary   string
	Available bool
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
}

func getSeriesBooks() []BookInfo {
	return []BookInfo{
		{
			Number:    1,
			Title:     "Void Reavers",
			Subtitle:  "A Tale of Space Pirates and Cosmic Plunder",
			Status:    "✓ Available",
			Available: true,
			Summary: `Captain Zara "Bloodhawk" Vega leads her pirate crew through the lawless void between solar systems. When humanity attracts the attention of ancient alien Architects, pirates become unlikely diplomats in a test that will determine if humans deserve a place among the stars. A thrilling space adventure where rogues and outlaws must save civilization itself.`,
		},
		{
			Number:    2,
			Title:     "Shadow Dancers",
			Subtitle:  "Echoes from Beyond",
			Status:    "Coming 2025",
			Available: false,
			Summary: `Dr. Elena Vasquez transcends physical reality to explore interdimensional realms. In spaces between spaces, she discovers ruins of civilizations that predate even the Architects and uncovers evidence of the universe's greatest threat: the reality-eating Devourers.`,
		},
		{
			Number:    3,
			Title:     "The Quantum Academy",
			Subtitle:  "Children of Two Realities",
			Status:    "Coming 2025",
			Available: false,
			Summary: `A new generation born with quantum abilities threatens the divide between enhanced and "pure" humans. Chen Wei must build academies to train these gifted children while preventing a civil war that could destroy humanity's cosmic probation.`,
		},
		{
			Number:    4,
			Title:     "Empire of Stars",
			Subtitle:  "The Corporate Renaissance",
			Status:    "Coming 2025",
			Available: false,
			Summary: `Mega-corporations form the Stellar Consortium for legitimate expansion, but some executives secretly fund neo-pirates to eliminate competition. Diana Marsh must prevent a new corporate war from destroying humanity's hard-won stability.`,
		},
		{
			Number:    5,
			Title:     "The Hegemony War",
			Subtitle:  "When Architects Sleep",
			Status:    "Coming 2026",
			Available: false,
			Summary: `The galaxy faces the Devourers—beings that can unmake existence itself. Admiral Lisa Park must choose: keep humanity safe in their small sector, or join a desperate alliance against entities from the universe's first epoch.`,
		},
		{
			Number:    6,
			Title:     "Ghosts of Morrison",
			Subtitle:  "The True Heirs' Revenge",
			Status:    "Coming 2026",
			Available: false,
			Summary: `Rex Morrison's descendants return from the galactic rim with a shadow empire and alien allies. Their twisted vision of humanity's expansion threatens everything the species has built under the Architects' guidance.`,
		},
		{
			Number:    7,
			Title:     "The Eternal Gambit",
			Subtitle:  "First Contact Protocol",
			Status:    "Coming 2026",
			Available: false,
			Summary: `Humanity joins an interdimensional council that shapes reality's fundamental rules. Elder Zara Vega leads quantum diplomats in proving humanity won't repeat their expansion mistakes on a universal scale.`,
		},
		{
			Number:    8,
			Title:     "Pirates of the Quantum Sea",
			Subtitle:  "The New Frontier",
			Status:    "Coming 2027",
			Available: false,
			Summary: `Descendants of the original pirates become Quantum Salvagers, rescuing civilizations from collapsed timelines. They discover the Devourers aren't destroyers—they're trying to return the universe to its original state.`,
		},
		{
			Number:    9,
			Title:     "The Architect's Dilemma",
			Subtitle:  "Guardians' Choice",
			Status:    "Coming 2027",
			Available: false,
			Summary: `The story from the Architects' perspective reveals their struggle with guiding younger species. Humanity's chaos challenged their orderly philosophy and changed them in ways they never expected.`,
		},
		{
			Number:    10,
			Title:     "New Horizons",
			Subtitle:  "Children of the Void",
			Status:    "Coming 2027",
			Available: false,
			Summary: `Centuries later, humanity has become gardeners of reality. When a new chaotic species emerges, humans must decide: guide them as the Architects did, or forge a new path that honors both order and chaos.`,
		},
	}
}

func initialModelWithUser(width, height int, username string) model {
	book, err := LoadBook("../book1_void_reavers_source")
	if err != nil {
		log.Printf("Error loading book: %v", err)
		book = &Book{
			Title: "Error Loading Book",
			Chapters: []Chapter{{Title: "Error", Content: "Could not load book content"}},
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
	menuItems = append(menuItems, "", "🚪 Exit")

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

func initialModel(width, height int) model {
	return initialModelWithUser(width, height, "anonymous")
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
		Width(leftWidth).
		Height(m.height - 8).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86"))
	
	var menuItems []string
	menuItems = append(menuItems, lipgloss.NewStyle().Bold(true).Render("  BOOK LIBRARY"))
	menuItems = append(menuItems, "")
	
	for i, book := range m.books {
		item := fmt.Sprintf("Book %d: %s", book.Number, book.Title)
		if i == m.menuCursor {
			if book.Available {
				menuItems = append(menuItems, selectedStyle.Render("▶ "+item+" ✓"))
			} else {
				menuItems = append(menuItems, selectedStyle.Render("▶ "+item))
			}
		} else {
			if book.Available {
				menuItems = append(menuItems, normalStyle.Render("  "+item+" ✓"))
			} else {
				statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				menuItems = append(menuItems, statusStyle.Render("  "+item))
			}
		}
	}
	
	// Add Exit option
	menuItems = append(menuItems, "")
	if m.menuCursor == len(m.menuItems)-1 {
		menuItems = append(menuItems, selectedStyle.Render("▶ 🚪 Exit"))
	} else {
		menuItems = append(menuItems, normalStyle.Render("  🚪 Exit"))
	}
	
	leftPanel := leftPanelStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, menuItems...),
	)
	
	// Right panel - Book details
	rightPanelStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height - 8).
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
	header := headerStyle.Width(m.width-2).Render("📚 CHAPTERS")

	var items []string
	for i, chapter := range m.book.Chapters {
		prefix := fmt.Sprintf("%2d. ", i+1)
		if i == m.chapterCursor {
			items = append(items, selectedStyle.Render("▶ "+prefix+chapter.Title))
		} else {
			items = append(items, normalStyle.Render("  "+prefix+chapter.Title))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)
	
	footer := footerStyle.Width(m.width).Render("↑/↓: navigate • enter: read • esc: back • q: quit")

	return lipgloss.JoinVertical(
		lipgloss.Top,
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
	
	header := headerStyle.Width(m.width-2).Render(fmt.Sprintf("📖 %s (%s)", chapter.Title, progress))

	contentHeight := m.height - 4 // Header + 2 empty lines + footer
	lines := wrapText(chapter.Content, m.width-4)
	
	startLine := m.scrollOffset
	endLine := startLine + contentHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	var visibleLines []string
	if startLine < len(lines) {
		visibleLines = lines[startLine:endLine]
	}

	content := contentStyle.Render(lipgloss.JoinVertical(lipgloss.Left, visibleLines...))

	navHelp := "h/←: prev chapter • l/→: next chapter • ↑/↓: scroll • b: bookmark • space: page down • esc: menu"
	footer := footerStyle.Width(m.width).Render(navHelp)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		"",
		content,
		footer,
	)
}

func (m model) viewAbout() string {
	header := headerStyle.Width(m.width-2).Render("ℹ️  ABOUT VOID REAVERS")

	aboutText := `🚀 Welcome to the Void Chronicles Universe! 🚀

Void Reavers is the first book in an epic 10-book science fiction series exploring humanity's evolution from chaotic space pirates to cosmic gardeners.

📖 The Story:
Follow Captain Zara "Bloodhawk" Vega's fifty-year journey as she transforms from a young pirate forced into Rex Morrison's brutal crew to humanity's ambassador to alien civilizations.

🌌 The Universe:
Set in a cosmos where quantum physics can tear reality apart and ancient alien Architects judge humanity's every move, pirates must evolve from raiders to protectors.

✨ Themes:
• Personal transformation mirrors species evolution
• The balance between order and chaos
• Earning cosmic citizenship through wisdom
• Honor among thieves in the vastness of space

🎭 Author: Captain J. Starwind
📅 Series: The Void Chronicles (Book 1 of 10)
🔧 Reader: Built with Go, Bubbletea, and Wish`

	content := contentStyle.Render(aboutText)
	footer := footerStyle.Width(m.width).Render("press any key to return to menu")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		"",
		content,
		"",
		footer,
	)
}

func (m model) viewProgress() string {
	header := headerStyle.Width(m.width-2).Render("📊 READING PROGRESS")

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

	content := contentStyle.Render(progressText)
	footer := footerStyle.Width(m.width).Render("press any key to return to menu")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		"",
		content,
		"",
		footer,
	)
}