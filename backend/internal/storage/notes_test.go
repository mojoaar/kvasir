package storage

import (
	"path/filepath"
	"testing"
)

func setupNotesDB(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestCreateNote(t *testing.T) {
	s := setupNotesDB(t)

	note := &Note{Title: "Hello", Content: "World", IsFolder: false}
	if err := s.CreateNote(note); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	if note.ID == 0 {
		t.Error("expected non-zero ID after create")
	}
	if note.CreatedAt.IsZero() {
		t.Error("expected created_at to be set")
	}
}

func TestCreateNoteFolder(t *testing.T) {
	s := setupNotesDB(t)

	folder := &Note{Title: "My Folder", IsFolder: true}
	if err := s.CreateNote(folder); err != nil {
		t.Fatalf("CreateNote folder: %v", err)
	}
	if !folder.IsFolder {
		t.Error("expected IsFolder to be true")
	}
}

func TestCreateNoteWithParent(t *testing.T) {
	s := setupNotesDB(t)

	folder := &Note{Title: "Parent", IsFolder: true}
	if err := s.CreateNote(folder); err != nil {
		t.Fatalf("CreateNote parent: %v", err)
	}

	child := &Note{Title: "Child", ParentID: &folder.ID, IsFolder: false}
	if err := s.CreateNote(child); err != nil {
		t.Fatalf("CreateNote child: %v", err)
	}
	if child.ParentID == nil || *child.ParentID != folder.ID {
		t.Errorf("expected parent_id=%d, got %v", folder.ID, child.ParentID)
	}
}

func TestGetNote(t *testing.T) {
	s := setupNotesDB(t)

	created := &Note{Title: "Test", Content: "Body", IsFolder: false}
	if err := s.CreateNote(created); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}

	got, err := s.GetNote(created.ID)
	if err != nil {
		t.Fatalf("GetNote: %v", err)
	}
	if got.Title != "Test" {
		t.Errorf("expected title 'Test', got %q", got.Title)
	}
	if got.Content != "Body" {
		t.Errorf("expected content 'Body', got %q", got.Content)
	}
}

func TestGetNoteNotFound(t *testing.T) {
	s := setupNotesDB(t)

	_, err := s.GetNote(999)
	if err == nil {
		t.Error("expected error for non-existent note")
	}
}

func TestGetNoteSoftDeleted(t *testing.T) {
	s := setupNotesDB(t)

	note := &Note{Title: "Gone", IsFolder: false}
	if err := s.CreateNote(note); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	if err := s.SoftDeleteNote(note.ID); err != nil {
		t.Fatalf("SoftDeleteNote: %v", err)
	}

	_, err := s.GetNote(note.ID)
	if err == nil {
		t.Error("expected error for soft-deleted note")
	}
}

func TestUpdateNote(t *testing.T) {
	s := setupNotesDB(t)

	note := &Note{Title: "Original", Content: "Old", IsFolder: false}
	if err := s.CreateNote(note); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}

	note.Title = "Updated"
	note.Content = "New"
	if err := s.UpdateNote(note); err != nil {
		t.Fatalf("UpdateNote: %v", err)
	}

	got, err := s.GetNote(note.ID)
	if err != nil {
		t.Fatalf("GetNote: %v", err)
	}
	if got.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %q", got.Title)
	}
	if got.Content != "New" {
		t.Errorf("expected content 'New', got %q", got.Content)
	}
}

func TestUpdateNoteMoveToFolder(t *testing.T) {
	s := setupNotesDB(t)

	folder := &Note{Title: "Folder", IsFolder: true}
	if err := s.CreateNote(folder); err != nil {
		t.Fatalf("CreateNote folder: %v", err)
	}

	note := &Note{Title: "Movable", IsFolder: false}
	if err := s.CreateNote(note); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}

	note.ParentID = &folder.ID
	if err := s.UpdateNote(note); err != nil {
		t.Fatalf("UpdateNote move: %v", err)
	}
	if note.ParentID == nil || *note.ParentID != folder.ID {
		t.Errorf("expected parent_id=%d, got %v", folder.ID, note.ParentID)
	}
}

func TestSoftDeleteNote(t *testing.T) {
	s := setupNotesDB(t)

	note := &Note{Title: "To Delete", IsFolder: false}
	if err := s.CreateNote(note); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}

	if err := s.SoftDeleteNote(note.ID); err != nil {
		t.Fatalf("SoftDeleteNote: %v", err)
	}

	got, err := s.GetNote(note.ID)
	if err == nil {
		t.Errorf("expected note to be gone, got %+v", got)
	}
}

func TestSoftDeleteNoteAlreadyDeleted(t *testing.T) {
	s := setupNotesDB(t)

	note := &Note{Title: "Double Delete", IsFolder: false}
	if err := s.CreateNote(note); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	if err := s.SoftDeleteNote(note.ID); err != nil {
		t.Fatalf("first SoftDeleteNote: %v", err)
	}

	err := s.SoftDeleteNote(note.ID)
	if err == nil {
		t.Error("expected error for already deleted note")
	}
}

func TestListNotesEmpty(t *testing.T) {
	s := setupNotesDB(t)

	notes, err := s.ListNotes(nil, nil, 0, 50)
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("expected 0 notes, got %d", len(notes))
	}
}

func TestListNotes(t *testing.T) {
	s := setupNotesDB(t)

	for _, title := range []string{"A", "B", "C"} {
		if err := s.CreateNote(&Note{Title: title, IsFolder: false}); err != nil {
			t.Fatalf("CreateNote %q: %v", title, err)
		}
	}

	notes, err := s.ListNotes(nil, nil, 0, 50)
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(notes) != 3 {
		t.Errorf("expected 3 notes, got %d", len(notes))
	}
}

func TestListNotesSoftDeleteExcluded(t *testing.T) {
	s := setupNotesDB(t)

	keep := &Note{Title: "Keep", IsFolder: false}
	if err := s.CreateNote(keep); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}

	del := &Note{Title: "Delete", IsFolder: false}
	if err := s.CreateNote(del); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	if err := s.SoftDeleteNote(del.ID); err != nil {
		t.Fatalf("SoftDeleteNote: %v", err)
	}

	notes, err := s.ListNotes(nil, nil, 0, 50)
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(notes))
	}
	if notes[0].ID != keep.ID {
		t.Errorf("expected note ID %d, got %d", keep.ID, notes[0].ID)
	}
}

func TestListNotesByParent(t *testing.T) {
	s := setupNotesDB(t)

	folder := &Note{Title: "Folder", IsFolder: true}
	if err := s.CreateNote(folder); err != nil {
		t.Fatalf("CreateNote folder: %v", err)
	}

	child := &Note{Title: "Child", ParentID: &folder.ID, IsFolder: false}
	if err := s.CreateNote(child); err != nil {
		t.Fatalf("CreateNote child: %v", err)
	}

	root := &Note{Title: "Root", IsFolder: false}
	if err := s.CreateNote(root); err != nil {
		t.Fatalf("CreateNote root: %v", err)
	}

	children, err := s.ListNotes(nil, &folder.ID, 0, 50)
	if err != nil {
		t.Fatalf("ListNotes by parent: %v", err)
	}
	if len(children) != 1 {
		t.Errorf("expected 1 child, got %d", len(children))
	}
	if children[0].ID != child.ID {
		t.Errorf("expected child ID %d, got %d", child.ID, children[0].ID)
	}
}

func TestListNotesPagination(t *testing.T) {
	s := setupNotesDB(t)

	for i := 0; i < 5; i++ {
		if err := s.CreateNote(&Note{Title: "N", IsFolder: false}); err != nil {
			t.Fatalf("CreateNote: %v", err)
		}
	}

	page, err := s.ListNotes(nil, nil, 0, 2)
	if err != nil {
		t.Fatalf("ListNotes page 1: %v", err)
	}
	if len(page) != 2 {
		t.Errorf("expected 2 notes on page 1, got %d", len(page))
	}

	page2, err := s.ListNotes(nil, nil, 2, 2)
	if err != nil {
		t.Fatalf("ListNotes page 2: %v", err)
	}
	if len(page2) != 2 {
		t.Errorf("expected 2 notes on page 2, got %d", len(page2))
	}
}

func TestListNotesFoldersFirst(t *testing.T) {
	s := setupNotesDB(t)

	if err := s.CreateNote(&Note{Title: "Note A", IsFolder: false}); err != nil {
		t.Fatalf("CreateNote: %v", err)
	}
	if err := s.CreateNote(&Note{Title: "Folder", IsFolder: true}); err != nil {
		t.Fatalf("CreateNote folder: %v", err)
	}

	notes, err := s.ListNotes(nil, nil, 0, 50)
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(notes) < 2 {
		t.Fatal("expected at least 2 notes")
	}
	if !notes[0].IsFolder {
		t.Error("expected folder first")
	}
}

func TestUpdateNoteNotFound(t *testing.T) {
	s := setupNotesDB(t)

	note := &Note{ID: 999, Title: "Ghost", IsFolder: false}
	err := s.UpdateNote(note)
	if err == nil {
		t.Error("expected error updating non-existent note")
	}
}

func TestDeleteNoteNotFound(t *testing.T) {
	s := setupNotesDB(t)

	err := s.SoftDeleteNote(999)
	if err == nil {
		t.Error("expected error deleting non-existent note")
	}
}
