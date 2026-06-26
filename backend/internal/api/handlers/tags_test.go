package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"kvasir/internal/storage"
)

func setupTagsRouter(t *testing.T) (*gin.Engine, *storage.Store) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	store, err := storage.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	r := gin.New()
	h := New(store)
	r.GET("/api/v1/tags", h.ListTags)
	r.POST("/api/v1/tags", h.CreateTag)
	r.GET("/api/v1/tags/:id", h.GetTag)
	r.PUT("/api/v1/tags/:id", h.UpdateTag)
	r.DELETE("/api/v1/tags/:id", h.DeleteTag)
	r.GET("/api/v1/notes/:id/tags", h.GetNoteTags)
	r.POST("/api/v1/notes/:id/tags", h.AddTagToNote)
	r.DELETE("/api/v1/notes/:id/tags", h.RemoveTagFromNote)
	return r, store
}

func TestListTagsHandlerEmpty(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tags", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var tags []storage.Tag
	json.NewDecoder(w.Body).Decode(&tags)
	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestCreateTagHandlerValid(t *testing.T) {
	r, _ := setupTagsRouter(t)

	body := map[string]string{"name": "important", "color": "#ef4444"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tags", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var tag storage.Tag
	json.NewDecoder(w.Body).Decode(&tag)
	if tag.Name != "important" {
		t.Errorf("expected name 'important', got %q", tag.Name)
	}
}

func TestCreateTagHandlerInvalidJSON(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tags", bytes.NewBuffer([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTagHandlerMissingName(t *testing.T) {
	r, _ := setupTagsRouter(t)

	body := map[string]string{"color": "#ef4444"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tags", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetTagHandlerFound(t *testing.T) {
	r, store := setupTagsRouter(t)

	tag := storage.Tag{Name: "bug", Color: "#f97316"}
	_ = store.CreateTag(&tag)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tags/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetTagHandlerNotFound(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tags/99999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetTagHandlerInvalidID(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tags/abc", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTagHandlerValid(t *testing.T) {
	r, store := setupTagsRouter(t)

	tag := storage.Tag{Name: "old", Color: "#000"}
	_ = store.CreateTag(&tag)

	body := map[string]string{"name": "new", "color": "#fff"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tags/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var got storage.Tag
	json.NewDecoder(w.Body).Decode(&got)
	if got.Name != "new" {
		t.Errorf("expected name 'new', got %q", got.Name)
	}
}

func TestUpdateTagHandlerNotFound(t *testing.T) {
	r, _ := setupTagsRouter(t)

	body := map[string]string{"name": "ghost", "color": "#fff"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tags/99999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestDeleteTagHandlerValid(t *testing.T) {
	r, store := setupTagsRouter(t)

	tag := storage.Tag{Name: "temp", Color: "#aaa"}
	_ = store.CreateTag(&tag)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tags/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAddTagToNoteHandlerValid(t *testing.T) {
	r, store := setupTagsRouter(t)

	note := storage.Note{Title: "note"}
	_ = store.CreateNote(&note)
	tag := storage.Tag{Name: "feature", Color: "#3b82f6"}
	_ = store.CreateTag(&tag)

	body := map[string]int64{"tagId": tag.ID}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes/1/tags", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAddTagToNoteHandlerInvalidJSON(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes/1/tags", bytes.NewBuffer([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRemoveTagFromNoteHandlerValid(t *testing.T) {
	r, store := setupTagsRouter(t)

	note := storage.Note{Title: "note"}
	_ = store.CreateNote(&note)
	tag := storage.Tag{Name: "feature", Color: "#3b82f6"}
	_ = store.CreateTag(&tag)
	_ = store.AddTagToNote(note.ID, tag.ID)

	body := map[string]int64{"tagId": tag.ID}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notes/1/tags", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetNoteTagsHandlerValid(t *testing.T) {
	r, store := setupTagsRouter(t)

	note := storage.Note{Title: "note"}
	_ = store.CreateNote(&note)
	tag := storage.Tag{Name: "urgent", Color: "#ef4444"}
	_ = store.CreateTag(&tag)
	_ = store.AddTagToNote(note.ID, tag.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes/1/tags", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var tags []storage.Tag
	json.NewDecoder(w.Body).Decode(&tags)
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
}

func TestGetNoteTagsHandlerEmpty(t *testing.T) {
	r, store := setupTagsRouter(t)

	note := storage.Note{Title: "untagged"}
	_ = store.CreateNote(&note)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes/1/tags", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var tags []storage.Tag
	json.NewDecoder(w.Body).Decode(&tags)
	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestGetNoteTagsHandlerInvalidID(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes/invalid/tags", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeleteTagHandlerNotFound(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/tags/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestRemoveTagFromNoteHandlerInvalidID(t *testing.T) {
	r, _ := setupTagsRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notes/invalid/tags", bytes.NewBuffer([]byte(`{"tagId":1}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAddTagToNoteHandlerInvalidNoteID(t *testing.T) {
	r, _ := setupTagsRouter(t)

	body := map[string]int64{"tagId": 1}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes/invalid/tags", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTagHandlerInvalidID(t *testing.T) {
	r, _ := setupTagsRouter(t)

	body := map[string]string{"name": "test", "color": "#fff"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/tags/invalid", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
