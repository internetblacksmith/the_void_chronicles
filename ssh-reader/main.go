package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

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
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath("../.ssh/id_ed25519"),
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
	log.Printf("Starting SSH server on %s", net.JoinHostPort(host, port))
	log.Printf("Password authentication enabled")
	go func() {
		if err = s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && err != ssh.ErrServerClosed {
		log.Fatalln(err)
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
}

func initialModelWithUser(width, height int, username string) model {
	book, err := LoadBook("../book1_void_reavers")
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

	return model{
		state:           menuView,
		book:            book,
		width:           width,
		height:          height,
		menuItems:       []string{"📖 Continue Reading", "📚 Chapter List", "📊 Progress", "ℹ️  About", "🚪 Exit"},
		progress:        progress,
		progressManager: pm,
		username:        username,
		currentChapter:  progress.CurrentChapter,
		scrollOffset:    progress.ScrollOffset,
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
		}
	case "down", "j":
		if m.menuCursor < len(m.menuItems)-1 {
			m.menuCursor++
		}
	case "enter", " ":
		switch m.menuCursor {
		case 0: // Continue Reading
			m.state = readingView
			// Use saved progress
			m.currentChapter = m.progress.CurrentChapter
			m.scrollOffset = m.progress.ScrollOffset
		case 1: // Chapter List
			m.state = chapterListView
			m.chapterCursor = 0
		case 2: // Progress
			m.state = progressView
		case 3: // About
			m.state = aboutView
		case 4: // Exit
			// Save progress before quitting
			m.progress.CurrentChapter = m.currentChapter
			m.progress.ScrollOffset = m.scrollOffset
			m.progressManager.SaveProgress(m.progress)
			m.quitting = true
			return m, tea.Quit
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
	title := titleStyle.Width(m.width).Render("🚀 VOID REAVERS 🚀")
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Width(m.width).
		Render("A Tale of Space Pirates and Cosmic Plunder")

	var menuItems []string
	for i, item := range m.menuItems {
		if i == m.menuCursor {
			menuItems = append(menuItems, selectedStyle.Render("▶ "+item))
		} else {
			menuItems = append(menuItems, normalStyle.Render("  "+item))
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		"",
		lipgloss.JoinVertical(lipgloss.Left, menuItems...),
	)

	footer := footerStyle.Width(m.width).Render("↑/↓: navigate • enter: select • q: quit")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.Place(m.width, m.height-2, lipgloss.Center, lipgloss.Center, content),
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