package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"kvasir/internal/storage"
)

func setupSearchHandler(t *testing.T) (*Handler, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	store, err := storage.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	h := New(store)
	r := gin.New()
	r.GET("/api/v1/search", h.Search)
	r.GET("/api/v1/search/tags", h.SearchByTag)

	return h, r
}

func TestSearchHandlerMissingQuery(t *testing.T) {
	_, r := setupSearchHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSearchHandlerEmptyResults(t *testing.T) {
	_, r := setupSearchHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search?q=nonexistent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []storage.SearchResult
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestSearchHandlerFindsNote(t *testing.T) {
	h, r := setupSearchHandler(t)

	note := storage.Note{
		Title:   "Kvasir Documentation",
		Content: "How to use the markdown editor",
	}
	if err := h.Store.CreateNote(&note); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search?q=Kvasir", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []storage.SearchResult
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Kvasir Documentation" {
		t.Errorf("expected title 'Kvasir Documentation', got %q", results[0].Title)
	}
}

func TestSearchHandlerFindsByContent(t *testing.T) {
	h, r := setupSearchHandler(t)

	note := storage.Note{
		Title:   "Untitled",
		Content: "This document explains the Nordic theme system",
	}
	if err := h.Store.CreateNote(&note); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search?q=Nordic", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []storage.SearchResult
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestSearchHandlerRespectsLimit(t *testing.T) {
	h, r := setupSearchHandler(t)

	for i := 0; i < 5; i++ {
		note := storage.Note{
			Title:   "Note about search",
			Content: "Testing search limit",
		}
		if err := h.Store.CreateNote(&note); err != nil {
			t.Fatalf("failed to create note: %v", err)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search?q=search&limit=3", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []storage.SearchResult
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(results) > 3 {
		t.Errorf("expected at most 3 results, got %d", len(results))
	}
}

func TestSearchHandlerExcludesSoftDeleted(t *testing.T) {
	h, r := setupSearchHandler(t)

	deleted := storage.Note{
		Title:   "Deleted Note",
		Content: "This should not appear in search",
	}
	if err := h.Store.CreateNote(&deleted); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}
	if err := h.Store.SoftDeleteNote(deleted.ID); err != nil {
		t.Fatalf("failed to delete note: %v", err)
	}

	active := storage.Note{
		Title:   "Active Note",
		Content: "This should appear in search",
	}
	if err := h.Store.CreateNote(&active); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search?q=appear", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []storage.SearchResult
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Active Note" {
		t.Errorf("expected 'Active Note', got %q", results[0].Title)
	}
}

func TestSearchHandlerHasSnippet(t *testing.T) {
	h, r := setupSearchHandler(t)

	note := storage.Note{
		Title:   "My Note",
		Content: "This is a note about project management with agile methodology",
	}
	if err := h.Store.CreateNote(&note); err != nil {
		t.Fatalf("failed to create note: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search?q=agile", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []storage.SearchResult
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(results) < 1 {
		t.Fatal("expected at least 1 result")
	}
	if results[0].Snippet == "" {
		t.Error("expected non-empty snippet")
	}
}

func TestSearchByTagHandlerMissingQuery(t *testing.T) {
	_, r := setupSearchHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search/tags", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSearchByTagHandlerEmptyResults(t *testing.T) {
	_, r := setupSearchHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/search/tags?q=nonexistent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var results []storage.Note
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
