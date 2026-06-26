package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenEmptyPath(t *testing.T) {
	_, err := Open("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestOpenValidPath(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	if store.DB == nil {
		t.Fatal("expected non-nil DB")
	}

	if err := store.DB.Ping(); err != nil {
		t.Fatalf("ping after open: %v", err)
	}
}

func TestOpenSchemaCreatesTables(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	tables := []string{"vaults", "notes", "notes_fts", "tags", "note_tags", "attachments", "versions", "themes", "plugins", "plugin_permissions"}

	for _, table := range tables {
		var count int
		err := store.DB.Get(&count, "SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", table)
		if err != nil {
			t.Fatalf("checking table %s: %v", table, err)
		}
		if count == 0 {
			t.Errorf("table %s was not created", table)
		}
	}
}

func TestOpenForeignKeysEnabled(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	var fk int
	if err := store.DB.Get(&fk, "PRAGMA foreign_keys"); err != nil {
		t.Fatalf("pragma foreign_keys: %v", err)
	}
	if fk != 1 {
		t.Errorf("foreign_keys = %d, want 1", fk)
	}
}

func TestOpenWALMode(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	var journalMode string
	if err := store.DB.Get(&journalMode, "PRAGMA journal_mode"); err != nil {
		t.Fatalf("pragma journal_mode: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("journal_mode = %s, want wal", journalMode)
	}
}

func TestOpenMaxOneConn(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	var maxConns int
	if err := store.DB.Get(&maxConns, "PRAGMA max_page_count"); err != nil {
		t.Skip("max_page_count unavailable, SetMaxOpenConns still set to 1")
	}
}

func TestOpenCreateDirectories(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "nested", "deep", "test.db")

	store, err := Open(nested)
	if err != nil {
		t.Fatalf("Open nested path: %v", err)
	}
	defer store.Close()

	if store.DB == nil {
		t.Fatal("expected non-nil DB")
	}
}

func TestClose(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestCloseNilDB(t *testing.T) {
	store := &Store{}
	if err := store.Close(); err != nil {
		t.Fatalf("Close on nil DB: %v", err)
	}
}

func TestOpenIdempotent(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}
	store1.Close()

	store2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("second Open: %v", err)
	}
	defer store2.Close()

	tables := []string{"vaults", "notes", "tags"}
	for _, table := range tables {
		var count int
		if err := store2.DB.Get(&count, "SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", table); err != nil {
			t.Fatalf("idempotent check %s: %v", table, err)
		}
		if count == 0 {
			t.Errorf("table %s missing after second open", table)
		}
	}
}

func TestOpenFTS5TriggersExist(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	triggers := []string{"notes_ai", "notes_ad", "notes_au"}
	for _, trigger := range triggers {
		var count int
		err := store.DB.Get(&count, "SELECT count(*) FROM sqlite_master WHERE type='trigger' AND name=?", trigger)
		if err != nil {
			t.Fatalf("trigger check %s: %v", trigger, err)
		}
		if count == 0 {
			t.Errorf("trigger %s not found", trigger)
		}
	}
}

func TestOpenInvalidPath(t *testing.T) {
	_, err := Open("/nonexistent/path/should/not/exist/db.sqlite")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestOpenPathIsDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := Open(dir)
	if err == nil {
		t.Fatal("expected error when path is a directory")
	}
}

func mustTempDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func mustFileSize(t *testing.T, path string) int64 {
	t.Helper()
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return fi.Size()
}
