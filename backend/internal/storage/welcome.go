package storage

const welcomeNoteContent = `# Welcome to Kvasir

**Kvasir** is your beautiful, techy, Nordic-inspired markdown knowledge base. This is your first note — a quick tour to get you started.

---

## The Editor

Kvasir uses a rich markdown editor with a toolbar for common formatting:

- **Bold**, *Italic*, ~~Strikethrough~~, ==Highlight==
- Headings (H1, H2, H3)
- Bullet lists, ordered lists, and task lists
- Tables, links, and code blocks
- Math formulas (KaTeX) and diagrams (Mermaid)

### Try It

Click anywhere in this note to start editing. The toolbar above gives you quick access to formatting. Press **Cmd+S** to save manually, or let the 2-second auto-save handle it.

---

## Themes

Kvasir ships with 5 themes, each with dark and light mode:

| Theme     | Dark Palette | Light Palette |
| --------- | ------------ | ------------- |
| Kvasir    | Nordic dark  | Nordic snow   |
| Dracula   | Purple dark  | Warm light    |
| Nord      | Polar night  | Snow storm    |
| GitHub    | Dark dimmed  | Default light |
| Cyberpunk | Neon black   | Neon white    |

Switch themes from the sidebar header or press **Cmd+K** and type "theme". Toggle dark/light mode with **Cmd+Shift+T** or the sun/moon button in the sidebar.

---

## Search

Kvasir uses full-text search powered by SQLite FTS5. Search is instant, even with thousands of notes.

- Type in the sidebar search box to search all notes
- Results show highlighted matches
- You can also search by tag

**Pro tip:** Use **Cmd+K** to open the command palette and search from anywhere.

---

## Keyboard Shortcuts

| Shortcut         | Action              |
| ---------------- | ------------------- |
| Cmd+K            | Command palette     |
| Cmd+S            | Save note           |
| Cmd+Shift+T      | Toggle dark/light   |
| Cmd+Z / Cmd+Shift+Z | Undo / Redo     |
| Enter            | New note (sidebar)  |

---

## Organizing Notes

Use the sidebar to organize your notes:

- Create folders to group related notes
- Drag and drop notes to reorganize
- Add tags for flexible categorization
- Double-click any note or folder in the sidebar to rename it

---

## Next Steps

1. Create a new note with **Cmd+K** → "Create Note" or click the + button in the sidebar
2. Try the different themes
3. Write some markdown — tables, code, math, and diagrams
4. Search across your notes

Welcome to the kernel of your mind. Happy note-taking!

`

func (s *Store) SeedIfEmpty() error {
	var count int
	if err := s.DB.Get(&count, `SELECT COUNT(*) FROM notes WHERE deleted_at IS NULL`); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	note := Note{
		Title:   "Welcome to Kvasir",
		Content: welcomeNoteContent,
	}
	return s.CreateNote(&note)
}
