package storage

import (
	"os"
	"testing"
)

func TestSeedIfEmpty_FreshDB(t *testing.T) {
	path := "testdata/seed_fresh.db"
	os.Remove(path)
	os.Remove(path + "-wal")
	os.Remove(path + "-shm")

	store, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()
	defer os.Remove(path)

	if err := store.SeedIfEmpty(); err != nil {
		t.Fatalf("SeedIfEmpty: %v", err)
	}

	var count int
	if err := store.DB.Get(&count, `SELECT COUNT(*) FROM notes WHERE deleted_at IS NULL`); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 note, got %d", count)
	}

	var title string
	if err := store.DB.Get(&title, `SELECT title FROM notes WHERE deleted_at IS NULL LIMIT 1`); err != nil {
		t.Fatalf("get title: %v", err)
	}
	if title != "Welcome to Kvasir" {
		t.Errorf("expected title 'Welcome to Kvasir', got %q", title)
	}
}

func TestSeedIfEmpty_Idempotent(t *testing.T) {
	path := "testdata/seed_idem.db"
	os.Remove(path)
	os.Remove(path + "-wal")
	os.Remove(path + "-shm")

	store, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()
	defer os.Remove(path)

	if err := store.SeedIfEmpty(); err != nil {
		t.Fatalf("first SeedIfEmpty: %v", err)
	}
	if err := store.SeedIfEmpty(); err != nil {
		t.Fatalf("second SeedIfEmpty: %v", err)
	}

	var count int
	if err := store.DB.Get(&count, `SELECT COUNT(*) FROM notes WHERE deleted_at IS NULL`); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 note after idempotent seed, got %d", count)
	}
}
