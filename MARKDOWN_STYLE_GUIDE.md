# Markdown Style Guide for Void Reavers

This guide ensures consistent formatting when editing the book in Markdown format.

## File Structure

```
book1_void_reavers_source/
├── book.md                 # Main book file with metadata
├── metadata.yaml          # Book metadata (title, author, ISBN, etc.)
└── chapters/
    ├── chapter-01.md      # Individual chapter files
    ├── chapter-02.md
    └── ...
```

## Chapter Format

Each chapter file should follow this structure:

```markdown
# Chapter N: Chapter Title

First paragraph of the chapter...

Second paragraph...
```

## Text Formatting

### Emphasis
- **Italics**: Use single asterisks for emphasis, especially ship names
  - `the *Crimson Revenge*` → the *Crimson Revenge*
- **Bold**: Use double asterisks for strong emphasis
  - `**"Stop right there!"**` → **"Stop right there!"**

### Dialogue
- Use standard double quotes: `"This is dialogue."`
- For quotes within dialogue: `"He said 'stop' but I didn't listen."`
- Em-dashes in dialogue: `"I thought—no, I knew—it was over."`

### Special Characters
- Em-dash: `—` (not `--` or `---`)
- En-dash: `–` (for ranges)
- Ellipsis: `...` (three periods)

## Paragraphs
- Separate paragraphs with a blank line
- No indentation needed (handled by output formatting)
- One sentence per line is acceptable for version control

## British English Spelling
Use British English throughout:
- honour (not honor)
- colour (not color)  
- realised (not realized)
- centre (not center)
- defence (not defense)

## Scene Breaks
For scene breaks within chapters, use three asterisks centered:

```markdown
The ship vanished into hyperspace.

* * *

Three days later, on Titan Station...
```

## Character Names
- First mention in a chapter: Full name
- Subsequent mentions: Last name or nickname
- Be consistent with nicknames: "Bloodhawk" or Zara, not both randomly

## Technical Terms
- Ship names: Always italicized (*Crimson Revenge*)
- Station names: Not italicized (Deadman's Port)
- Planet names: Not italicized (Mars, Kepler-442b)

## Common Mistakes to Avoid
1. Don't use LaTeX commands (`\textit{}`, `\chapter{}`)
2. Don't use HTML tags (`<em>`, `<strong>`)
3. Don't use underscores for emphasis (conflicts with Markdown)
4. Keep line lengths reasonable for readability

## Editing Workflow

1. Edit the Markdown source files directly
2. Use any text editor (VS Code, Sublime, even Notepad)
3. Preview changes with a Markdown viewer
4. Generate output formats using the provided scripts:
   ```bash
   ./markdown_to_epub.rb book1_void_reavers_source
   ./markdown_to_kdp_pdf.rb book1_void_reavers_source
   ```

## Version Control
- Commit changes frequently
- Use meaningful commit messages
- One chapter per commit is ideal for tracking changes

## Adding New Chapters
1. Create new file: `chapters/chapter-21.md`
2. Add chapter header: `# Chapter 21: New Beginning`
3. Update `metadata.yaml` to include the new chapter
4. Update `book.md` table of contents

## Metadata File (metadata.yaml)

```yaml
title: "Void Reavers"
subtitle: "A Tale of Space Pirates and Cosmic Plunder"
author: "Captain J. Starwind"
language: en-GB
rights: © 2024 Captain J. Starwind
isbn: "YOUR-ISBN-HERE"
publisher: "Self Published"
publication_date: "2024"
chapters:
  - number: 1
    title: "The Void Between Stars"
    file: "chapters/chapter-01.md"
  # ... etc
```

## Quality Checklist
- [ ] British English spelling throughout
- [ ] Consistent character names
- [ ] Ship names italicized
- [ ] Proper paragraph breaks
- [ ] No LaTeX artifacts
- [ ] Em-dashes used correctly
- [ ] Dialogue punctuation correct