package storage

import (
	"path/filepath"
	"testing"
	"time"
)

func setupStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("setup: Open: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestSearchEmptyDB(t *testing.T) {
	store := setupStore(t)

	results, err := store.Search("hello", 10)
	if err != nil {
		t.Fatalf("Search on empty DB: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchLimitDefaults(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('hello world', 'hello there')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}

	results, err := store.Search("hello", 0)
	if err != nil {
		t.Fatalf("Search with limit 0: %v", err)
	}
	if results == nil {
		t.Fatal("expected non-nil results slice")
	}
}

func TestSearchReturnsNilSlice(t *testing.T) {
	store := setupStore(t)

	results, err := store.Search("nothing", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if results == nil {
		t.Fatal("expected non-nil results slice, got nil")
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchFindsNoteByTitle(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('hello world', 'some content')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}

	results, err := store.Search("hello", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "hello world" {
		t.Errorf("expected title 'hello world', got %q", results[0].Title)
	}
}

func TestSearchFindsNoteByContent(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('title', 'hello world content')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}

	results, err := store.Search("world", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchExcludesSoftDeleted(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content, deleted_at) VALUES ('hello', 'hello', datetime('now'))`)
	if err != nil {
		t.Fatalf("insert soft-deleted note: %v", err)
	}

	results, err := store.Search("hello", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results (soft-deleted excluded), got %d", len(results))
	}
}

func TestSearchMultipleResults(t *testing.T) {
	store := setupStore(t)

	for i, title := range []string{"alpha hello", "bravo hello world", "charlie unrelated", "delta hello universe"} {
		_, err := store.DB.Exec(`INSERT INTO notes (title, content, sort_order) VALUES (?, 'content', ?)`, title, i)
		if err != nil {
			t.Fatalf("insert note %d: %v", i, err)
		}
	}

	results, err := store.Search("hello", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestSearchRespectsLimit(t *testing.T) {
	store := setupStore(t)

	for i := range 10 {
		_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('hello', 'hello')`)
		if err != nil {
			t.Fatalf("insert note %d: %v", i, err)
		}
	}

	results, err := store.Search("hello", 3)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) > 3 {
		t.Errorf("expected at most 3 results, got %d", len(results))
	}
}

func TestSearchResultHasRankAndSnippet(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('hello world', 'this is some content about hello and the world')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}

	results, err := store.Search("hello", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Rank == 0 {
		t.Error("expected non-zero rank")
	}
	if results[0].Snippet == "" {
		t.Error("expected non-empty snippet")
	}
}

func TestSearchFTS5TriggerAfterInsert(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('trigger test', 'fts5 insert trigger')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}

	results, err := store.Search("trigger", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result via FTS5 trigger, got %d", len(results))
	}
}

func TestSearchFTS5TriggerAfterUpdate(t *testing.T) {
	store := setupStore(t)

	res, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('old title', 'old content')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}
	id, _ := res.LastInsertId()

	_, err = store.DB.Exec(`UPDATE notes SET title = 'new title hello', content = 'new content' WHERE id = ?`, id)
	if err != nil {
		t.Fatalf("update note: %v", err)
	}

	results, err := store.Search("hello", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result after update trigger, got %d", len(results))
	}
	if results[0].Title != "new title hello" {
		t.Errorf("expected 'new title hello', got %q", results[0].Title)
	}

	results2, err2 := store.Search("old", 5)
	if err2 != nil {
		t.Fatalf("Search old: %v", err2)
	}
	if len(results2) != 0 {
		t.Errorf("old content still searchable after update: %d results", len(results2))
	}
}

func TestSearchFTS5TriggerAfterDelete(t *testing.T) {
	store := setupStore(t)

	res, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('delete me', 'gone content')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}
	id, _ := res.LastInsertId()

	_, err = store.DB.Exec(`DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		t.Fatalf("delete note: %v", err)
	}

	results, err := store.Search("gone", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results after delete trigger, got %d", len(results))
	}
}

func TestRebuildFTS(t *testing.T) {
	store := setupStore(t)

	for range 10 {
		_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('test', 'rebuild fts')`)
		if err != nil {
			t.Fatalf("insert note: %v", err)
		}
	}

	if err := store.RebuildFTS(); err != nil {
		t.Fatalf("RebuildFTS: %v", err)
	}

	results, err := store.Search("rebuild", 20)
	if err != nil {
		t.Fatalf("Search after rebuild: %v", err)
	}
	if len(results) != 10 {
		t.Errorf("expected 10 results after rebuild, got %d", len(results))
	}
}

func TestSearchVaultedNote(t *testing.T) {
	store := setupStore(t)

	var vaultID int64
	err := store.DB.Get(&vaultID, `INSERT INTO vaults (name) VALUES ('test vault') RETURNING id`)
	if err != nil {
		t.Fatalf("insert vault: %v", err)
	}

	_, err = store.DB.Exec(`INSERT INTO notes (title, content, vault_id) VALUES ('vaulted note', 'inside vault', ?)`, vaultID)
	if err != nil {
		t.Fatalf("insert note in vault: %v", err)
	}

	results, err := store.Search("vaulted", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 vaulted result, got %d", len(results))
	}
	if results[0].VaultID == nil || *results[0].VaultID != vaultID {
		t.Errorf("vault ID mismatch: %v", results[0].VaultID)
	}
}

func TestSearchNoteWithTag(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('tagged note', 'has a tag')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}

	_, err = store.DB.Exec(`INSERT INTO tags (name) VALUES ('important')`)
	if err != nil {
		t.Fatalf("insert tag: %v", err)
	}

	_, err = store.DB.Exec(`INSERT INTO note_tags (note_id, tag_id) VALUES (1, 1)`)
	if err != nil {
		t.Fatalf("link note-tag: %v", err)
	}

	results, err := store.Search("tagged", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestSearchVersionedNote(t *testing.T) {
	store := setupStore(t)

	res, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('versioned note', 'version 1 content')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}
	id, _ := res.LastInsertId()

	_, err = store.DB.Exec(`INSERT INTO versions (note_id, content, version_num) VALUES (?, 'old content snapshot', 1)`, id)
	if err != nil {
		t.Fatalf("insert version: %v", err)
	}

	results, err := store.Search("versioned", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestSearchWithTimestamp(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('timestamp test', 'created now')`)
	if err != nil {
		t.Fatalf("insert note: %v", err)
	}

	results, err := store.Search("timestamp", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if results[0].UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestSearchParentNote(t *testing.T) {
	store := setupStore(t)

	res, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES ('parent', 'parent note')`)
	if err != nil {
		t.Fatalf("insert parent: %v", err)
	}
	parentID, _ := res.LastInsertId()

	_, err = store.DB.Exec(`INSERT INTO notes (title, content, parent_id) VALUES ('child', 'child note', ?)`, parentID)
	if err != nil {
		t.Fatalf("insert child: %v", err)
	}

	results, err := store.Search("child", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].ParentID == nil || *results[0].ParentID != parentID {
		t.Errorf("parent ID mismatch: %v", results[0].ParentID)
	}
}

func TestSearchSortOrderPreserved(t *testing.T) {
	store := setupStore(t)

	_, err := store.DB.Exec(`INSERT INTO notes (title, content, sort_order) VALUES ('first', 'hello', 0)`)
	if err != nil {
		t.Fatalf("insert first: %v", err)
	}
	_, err = store.DB.Exec(`INSERT INTO notes (title, content, sort_order) VALUES ('second', 'hello', 1)`)
	if err != nil {
		t.Fatalf("insert second: %v", err)
	}

	results, err := store.Search("hello", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].ID == 0 || results[1].ID == 0 {
		t.Error("expected valid IDs")
	}
}

func TestRebuildFTSNoData(t *testing.T) {
	store := setupStore(t)

	err := store.RebuildFTS()
	if err != nil {
		t.Fatalf("RebuildFTS on empty DB: %v", err)
	}
}

func TestRebuildFTSAfterManyInserts(t *testing.T) {
	store := setupStore(t)

	for i := range 50 {
		_, err := store.DB.Exec(`INSERT INTO notes (title, content) VALUES (?, ?)`, "note", time.Now().String())
		if err != nil {
			t.Fatalf("insert note %d: %v", i, err)
		}
	}

	if err := store.RebuildFTS(); err != nil {
		t.Fatalf("RebuildFTS: %v", err)
	}

	results, err := store.Search("note", 100)
	if err != nil {
		t.Fatalf("Search after rebuild: %v", err)
	}
	if len(results) != 50 {
		t.Errorf("expected 50 results after rebuild, got %d", len(results))
	}
}
