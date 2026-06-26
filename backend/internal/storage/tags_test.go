package storage

import (
	"path/filepath"
	"testing"
)

func setupTagsStore(t *testing.T) *Store {
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

func TestListTagsEmpty(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	tags, err := s.ListTags()
	if err != nil {
		t.Fatalf("ListTags failed: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestCreateTag(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	tag := Tag{Name: "important", Color: "#ef4444"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}
	if tag.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if tag.Name != "important" {
		t.Errorf("expected name 'important', got %q", tag.Name)
	}
	if tag.Color != "#ef4444" {
		t.Errorf("expected color '#ef4444', got %q", tag.Color)
	}
}

func TestGetTagFound(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	tag := Tag{Name: "bug", Color: "#f97316"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	got, err := s.GetTag(tag.ID)
	if err != nil {
		t.Fatalf("GetTag failed: %v", err)
	}
	if got.Name != "bug" {
		t.Errorf("expected name 'bug', got %q", got.Name)
	}
}

func TestGetTagNotFound(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	_, err := s.GetTag(99999)
	if err == nil {
		t.Error("expected error for non-existent tag")
	}
}

func TestUpdateTag(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	tag := Tag{Name: "old", Color: "#000000"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	tag.Name = "new"
	tag.Color = "#ffffff"
	if err := s.UpdateTag(&tag); err != nil {
		t.Fatalf("UpdateTag failed: %v", err)
	}

	got, _ := s.GetTag(tag.ID)
	if got.Name != "new" {
		t.Errorf("expected name 'new', got %q", got.Name)
	}
	if got.Color != "#ffffff" {
		t.Errorf("expected color '#ffffff', got %q", got.Color)
	}
}

func TestUpdateTagNotFound(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	tag := Tag{ID: 99999, Name: "ghost", Color: "#000"}
	err := s.UpdateTag(&tag)
	if err == nil {
		t.Error("expected error for non-existent tag update")
	}
}

func TestDeleteTag(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	tag := Tag{Name: "temp", Color: "#aaa"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	if err := s.DeleteTag(tag.ID); err != nil {
		t.Fatalf("DeleteTag failed: %v", err)
	}

	_, err := s.GetTag(tag.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestDeleteTagNotFound(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	err := s.DeleteTag(99999)
	if err == nil {
		t.Error("expected error for deleting non-existent tag")
	}
}

func TestListTagsWithData(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	tags := []Tag{
		{Name: "alpha", Color: "#111"},
		{Name: "beta", Color: "#222"},
		{Name: "gamma", Color: "#333"},
	}
	for i := range tags {
		if err := s.CreateTag(&tags[i]); err != nil {
			t.Fatalf("CreateTag failed: %v", err)
		}
	}

	all, err := s.ListTags()
	if err != nil {
		t.Fatalf("ListTags failed: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 tags, got %d", len(all))
	}
	if all[0].Name != "alpha" {
		t.Errorf("expected first tag 'alpha', got %q", all[0].Name)
	}
}

func TestAddTagToNote(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	note := Note{Title: "test note"}
	if err := s.CreateNote(&note); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	tag := Tag{Name: "feature", Color: "#3b82f6"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	if err := s.AddTagToNote(note.ID, tag.ID); err != nil {
		t.Fatalf("AddTagToNote failed: %v", err)
	}

	noteTags, err := s.GetNoteTags(note.ID)
	if err != nil {
		t.Fatalf("GetNoteTags failed: %v", err)
	}
	if len(noteTags) != 1 {
		t.Fatalf("expected 1 tag on note, got %d", len(noteTags))
	}
	if noteTags[0].Name != "feature" {
		t.Errorf("expected tag 'feature', got %q", noteTags[0].Name)
	}
}

func TestRemoveTagFromNote(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	note := Note{Title: "test note"}
	if err := s.CreateNote(&note); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	tag := Tag{Name: "feature", Color: "#3b82f6"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	_ = s.AddTagToNote(note.ID, tag.ID)

	if err := s.RemoveTagFromNote(note.ID, tag.ID); err != nil {
		t.Fatalf("RemoveTagFromNote failed: %v", err)
	}

	noteTags, _ := s.GetNoteTags(note.ID)
	if len(noteTags) != 0 {
		t.Errorf("expected 0 tags after removal, got %d", len(noteTags))
	}
}

func TestGetNoteTagsEmpty(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	note := Note{Title: "untagged"}
	if err := s.CreateNote(&note); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	tags, err := s.GetNoteTags(note.ID)
	if err != nil {
		t.Fatalf("GetNoteTags failed: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestGetNoteTagsMultiple(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	note := Note{Title: "tagged"}
	if err := s.CreateNote(&note); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	tags := []Tag{
		{Name: "urgent", Color: "#ef4444"},
		{Name: "backend", Color: "#3b82f6"},
		{Name: "docs", Color: "#10b981"},
	}
	for i := range tags {
		if err := s.CreateTag(&tags[i]); err != nil {
			t.Fatalf("CreateTag failed: %v", err)
		}
		if err := s.AddTagToNote(note.ID, tags[i].ID); err != nil {
			t.Fatalf("AddTagToNote failed: %v", err)
		}
	}

	noteTags, err := s.GetNoteTags(note.ID)
	if err != nil {
		t.Fatalf("GetNoteTags failed: %v", err)
	}
	if len(noteTags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(noteTags))
	}
}

func TestDeleteTagCascadesNoteTags(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	note := Note{Title: "note"}
	if err := s.CreateNote(&note); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	tag := Tag{Name: "stale", Color: "#666"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}
	_ = s.AddTagToNote(note.ID, tag.ID)

	if err := s.DeleteTag(tag.ID); err != nil {
		t.Fatalf("DeleteTag failed: %v", err)
	}

	noteTags, _ := s.GetNoteTags(note.ID)
	if len(noteTags) != 0 {
		t.Errorf("expected 0 tags after tag deletion, got %d", len(noteTags))
	}
}

func TestAddTagDuplicate(t *testing.T) {
	s := setupTagsStore(t)
	defer s.Close()

	note := Note{Title: "note"}
	if err := s.CreateNote(&note); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}
	tag := Tag{Name: "dup", Color: "#999"}
	if err := s.CreateTag(&tag); err != nil {
		t.Fatalf("CreateTag failed: %v", err)
	}

	_ = s.AddTagToNote(note.ID, tag.ID)
	if err := s.AddTagToNote(note.ID, tag.ID); err != nil {
		t.Fatalf("AddTagToNote duplicate should not error: %v", err)
	}

	noteTags, _ := s.GetNoteTags(note.ID)
	if len(noteTags) != 1 {
		t.Errorf("expected 1 tag after duplicate add, got %d", len(noteTags))
	}
}
