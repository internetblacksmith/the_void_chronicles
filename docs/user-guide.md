# User Guide

Complete guide to using the Void Reavers SSH Reader - from basic navigation to advanced features.

## 🎯 Getting Started

### Your First Connection

Connect to the SSH reader:
```bash
ssh localhost -p 23234
```

You'll be greeted with the main menu:

```
🚀 VOID REAVERS 🚀
A Tale of Space Pirates and Cosmic Plunder

▶ 📖 Continue Reading
  📚 Chapter List  
  📊 Progress
  ℹ️  About
  🚪 Exit
```

## 🎮 Navigation Reference

### Universal Controls
- `q`: Quit application (from any screen)
- `Ctrl+C`: Force quit
- `Esc`: Return to previous screen/main menu
- `↑/↓` or `k/j`: Navigate up/down (vi-style)
- `Enter` or `Space`: Select/confirm

### Main Menu Navigation
- **📖 Continue Reading**: Resume from last position
- **📚 Chapter List**: Browse all chapters
- **📊 Progress**: View reading statistics and bookmarks
- **ℹ️ About**: Information about the book and application
- **🚪 Exit**: Quit the application

## 📖 Reading Experience

### Reading View Controls

#### Basic Navigation
- `↑/↓` or `k/j`: Scroll line by line
- `Page Up/Down`: Scroll by full page
- `Space`: Page down (same as Page Down)
- `Ctrl+B/F`: Page up/down (emacs-style)

#### Chapter Navigation
- `h/←` or `p`: Previous chapter
- `l/→` or `n`: Next chapter
- Numbers `1-9`: Quick jump to chapters 1-9
- `0`: Jump to chapter 10

#### Position Controls
- `g` or `Home`: Go to beginning of chapter
- `G` or `End`: Go to end of chapter
- `/`: Search within current chapter (if implemented)

#### Bookmarking
- `b`: Add bookmark at current position
- `B`: View bookmark list for current chapter

### Reading Features

#### Auto-Save Progress
Your reading position is automatically saved when you:
- Move to a different chapter
- Exit the application
- Return to the main menu

#### Chapter Completion
When you move to the next chapter (`l`, `→`, or `n`), the current chapter is automatically marked as completed in your progress.

#### Text Formatting
The reader supports various text formatting:
- **Bold text** displays in bright colors
- *Italic text* shows in emphasized styling
- Regular paragraphs with proper word wrapping
- Automatic text wrapping based on terminal width

## 📚 Chapter List

### Navigation
- `↑/↓` or `k/j`: Browse through chapters
- `Enter`: Jump directly to selected chapter
- `Esc`: Return to main menu

### Chapter Status Indicators
- **✅**: Completed chapters
- **📍**: Current chapter (where you left off)
- **⭕**: Unread chapters

### Chapter Information
Each chapter shows:
- Chapter number (1-20)
- Chapter title
- Completion status

## 📊 Progress Tracking

### Progress View

Access your reading statistics via **📊 Progress** from the main menu:

```
📊 READING PROGRESS

👤 Reader: your-username

📈 Overall Progress: 45.2% Complete
📖 Current Chapter: 9 of 20
⏱️  Total Reading Time: 2h 34m
📅 Last Read: Jan 15, 2024 14:30
🔖 Bookmarks: 7

📚 Chapter Progress:
  ✅ Chapter 1: The Void Between Stars
  ✅ Chapter 2: The Crimson Nebula
  ...
  📍 Chapter 9: Shadows of the Asteroid Belt
  ⭕ Chapter 10: The Corporate Armada
  ...
```

### Progress Features

#### Automatic Tracking
- **Reading Time**: Tracks total time spent reading
- **Session Time**: Time spent in current session
- **Chapter Completion**: Automatically marks chapters as done
- **Last Position**: Remembers exact scroll position

#### Statistics
- **Completion Percentage**: Overall book progress
- **Chapter Status**: Visual indicators for each chapter
- **Reading Speed**: Estimated based on time and progress
- **Session History**: When you last read

## 🔖 Bookmark System

### Creating Bookmarks

#### Quick Bookmark
Press `b` while reading to create a bookmark at your current position:
- Saves current chapter and scroll position
- Timestamp is automatically added
- No interruption to reading flow

#### Bookmark with Note (Future Feature)
Advanced bookmarking with custom notes for important passages.

### Managing Bookmarks

#### Viewing Bookmarks
From the Progress screen, see your recent bookmarks:
```
🔖 Recent Bookmarks:
  • Chapter 3 - Jan 14 09:15
  • Chapter 5 - Jan 14 14:22
  • Chapter 7 - Jan 15 10:45
```

#### Navigating to Bookmarks
- Select bookmark from progress screen
- Jump directly to bookmarked location
- Resume reading from exact position

### Bookmark Features

#### Automatic Organization
- Sorted by creation date (newest first)
- Grouped by chapter
- Duplicate protection (won't bookmark same position twice)

#### Persistence
- Bookmarks survive application restarts
- Unique per user
- Backed up with progress data

## 👥 Multi-User Support

### User Identification
Each SSH connection creates a separate user session:
- Anonymous users get default tracking
- Named users maintain individual progress
- No authentication required for reading

### Individual Progress
Each user maintains:
- **Separate reading position**: Your chapter and scroll location
- **Individual bookmarks**: Your personal bookmark collection
- **Personal statistics**: Your reading time and completion
- **Independent history**: Your reading sessions

### Data Isolation
- User data stored in separate files
- No cross-user data sharing
- Privacy protection built-in

## 🎨 Interface Customization

### Terminal Compatibility

#### Optimal Experience
- **Terminal Size**: 80x24 minimum, 120x40+ recommended
- **Color Support**: 256-color terminals for best styling
- **Font**: Monospace font with Unicode support
- **SSH Client**: Modern SSH client with proper terminal emulation

#### Responsive Design
The interface automatically adapts to:
- **Terminal Width**: Text wraps to fit your terminal
- **Terminal Height**: Page scrolling adjusts to window size
- **Color Support**: Graceful degradation on limited color terminals

### Visual Elements

#### Styling
- **Headers**: Colored borders and titles
- **Selection**: Highlighted current option
- **Status**: Progress bars and indicators
- **Text**: Proper formatting preservation

#### Emojis and Icons
- Intuitive emoji icons for menu items
- Status indicators (✅, 📍, ⭕)
- Visual feedback for actions

## 🔧 Advanced Features

### Search (Future Feature)
- Search within current chapter
- Search across entire book
- Bookmark search results

### Export (Future Feature)
- Export reading progress
- Export bookmarks
- Generate reading reports

### Themes (Future Feature)
- Color scheme customization
- High contrast mode
- Custom styling options

## 💡 Tips and Tricks

### Efficiency Tips

#### Keyboard Navigation
- Learn vi-style keys (`hjkl`) for faster navigation
- Use `g` and `G` for quick chapter beginning/end jumps
- Master page scrolling with `Space` and `Page Up/Down`

#### Reading Strategy
- Use bookmarks liberally for important passages
- Check progress regularly to track completion
- Take advantage of auto-save - no need to manually save

#### Terminal Setup
- Use a comfortable terminal size (120x40 is great)
- Enable color support for better visual experience
- Consider using a dedicated terminal profile for reading

### Productivity Features

#### Session Management
- Multiple SSH connections allow multiple reading positions
- Leave connections open to maintain reading context
- Use different terminals for different reading sessions

#### Progress Tracking
- Set reading goals based on completion percentage
- Use reading time tracking to monitor habits
- Bookmark interesting quotes or key plot points

## 🐛 Common Issues and Solutions

### Reading Issues

#### Text Too Small/Large
- Adjust your terminal font size
- Resize terminal window for better text wrapping
- Use zoom features in your terminal emulator

#### Formatting Problems
- Ensure terminal supports UTF-8 encoding
- Check color support in your SSH client
- Try a different terminal emulator if issues persist

#### Navigation Not Working
- Verify key bindings in your terminal
- Check if SSH client is capturing key combinations
- Try alternative navigation keys (arrows vs. hjkl)

### Connection Issues

#### SSH Connection Fails
- Verify server is running: `ps aux | grep void-reader`
- Check port availability: `netstat -tlnp | grep 23234`
- Test local connection first: `ssh localhost -p 23234`

#### Disconnection Problems
- Check network stability
- Verify SSH client timeout settings
- Look for server logs if repeatedly disconnecting

### Data Issues

#### Progress Not Saving
- Check file permissions in `.void_reader_data/`
- Verify disk space availability
- Look for error messages in server logs

#### Bookmarks Not Working
- Ensure you're using `b` key while reading (not in menu)
- Check progress screen to verify bookmarks are saved
- Try creating bookmark at different position

## 🆘 Getting Help

### Self-Help Resources
1. **About Screen**: Basic information about the application
2. **Progress Screen**: Shows your current status and any issues
3. **Server Logs**: Check console output where server is running

### Advanced Troubleshooting
1. **Restart Application**: Exit and restart the SSH reader
2. **Check File Permissions**: Ensure user data directory is writable
3. **Test Different Terminal**: Try connecting from different SSH client
4. **Review Configuration**: Check server configuration for issues

### Documentation
- [Installation Guide](installation.md): Setup and configuration issues
- [Configuration Guide](configuration.md): Advanced customization
- [Troubleshooting Guide](troubleshooting.md): Detailed problem solving
- [Development Guide](development.md): Technical details and modifications

---

**Happy Reading!** 📚✨

*"The best way to read a book is the way that works for you - and in the void between stars, every reader finds their own path."*

---

**Next Steps:**
- Master the interface? Check out [Configuration Guide](configuration.md) for advanced options
- Having issues? See [Troubleshooting Guide](troubleshooting.md)
- Want to contribute? Read [Development Guide](development.md)