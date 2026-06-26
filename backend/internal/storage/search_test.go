package storage

import (
	"path/filepath"
	"testing"
)

func setupSearchDB(t *testing.T) *Store {
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

func TestSearchByTag(t *testing.T) {
	st := setupSearchDB(t)

	note := &Note{Title: "Tagged Note", Content: "Test content"}
	if err := st.CreateNote(note); err != nil {
		t.Fatal(err)
	}

	tag := &Tag{Name: "test-tag", Color: "#ff0000"}
	if err := st.CreateTag(tag); err != nil {
		t.Fatal(err)
	}
	if err := st.AddTagToNote(note.ID, tag.ID); err != nil {
		t.Fatal(err)
	}

	results, err := st.SearchByTag("test-tag")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != note.ID {
		t.Errorf("expected note ID %d, got %d", note.ID, results[0].ID)
	}
}

func TestSearchByTag_EmptyResults(t *testing.T) {
	st := setupSearchDB(t)

	results, err := st.SearchByTag("nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchByTag_ExcludesSoftDeleted(t *testing.T) {
	st := setupSearchDB(t)

	note := &Note{Title: "Will Be Deleted", Content: "Test"}
	if err := st.CreateNote(note); err != nil {
		t.Fatal(err)
	}
	tag := &Tag{Name: "delete-test", Color: "#000"}
	if err := st.CreateTag(tag); err != nil {
		t.Fatal(err)
	}
	if err := st.AddTagToNote(note.ID, tag.ID); err != nil {
		t.Fatal(err)
	}
	if err := st.SoftDeleteNote(note.ID); err != nil {
		t.Fatal(err)
	}

	results, err := st.SearchByTag("delete-test")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for soft-deleted note, got %d", len(results))
	}
}

func TestSearchByTag_PartialMatch(t *testing.T) {
	st := setupSearchDB(t)

	note := &Note{Title: "Partial Match", Content: "Test"}
	if err := st.CreateNote(note); err != nil {
		t.Fatal(err)
	}
	tag := &Tag{Name: "cool-stuff", Color: "#00ff00"}
	if err := st.CreateTag(tag); err != nil {
		t.Fatal(err)
	}
	if err := st.AddTagToNote(note.ID, tag.ID); err != nil {
		t.Fatal(err)
	}

	results, err := st.SearchByTag("cool")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result for partial match, got %d", len(results))
	}
}
