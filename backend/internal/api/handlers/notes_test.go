package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"kvasir/internal/storage"
)

func setupNotesHandler(t *testing.T) (*Handler, *gin.Engine) {
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
	r.GET("/api/v1/notes", h.ListNotes)
	r.POST("/api/v1/notes", h.CreateNote)
	r.GET("/api/v1/notes/:id", h.GetNote)
	r.PUT("/api/v1/notes/:id", h.UpdateNote)
	r.DELETE("/api/v1/notes/:id", h.DeleteNote)

	return h, r
}

func TestListNotesHandlerEmpty(t *testing.T) {
	_, r := setupNotesHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestCreateNoteHandler(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]interface{}{
		"title":   "Test Note",
		"content": "Hello world",
	}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var note storage.Note
	if err := json.Unmarshal(w.Body.Bytes(), &note); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if note.Title != "Test Note" {
		t.Errorf("expected title 'Test Note', got %q", note.Title)
	}
	if note.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestCreateNoteHandlerInvalidJSON(t *testing.T) {
	_, r := setupNotesHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateNoteHandlerMissingTitle(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]interface{}{"content": "no title"}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateNoteHandlerFolder(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]interface{}{
		"title":    "My Folder",
		"isFolder": true,
	}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var note storage.Note
	json.Unmarshal(w.Body.Bytes(), &note)
	if !note.IsFolder {
		t.Error("expected IsFolder to be true")
	}
}

func TestGetNoteHandler(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]interface{}{"title": "Get Me"}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var created storage.Note
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/notes/%d", created.ID), nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w2.Code)
	}

	var fetched storage.Note
	json.Unmarshal(w2.Body.Bytes(), &fetched)
	if fetched.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, fetched.ID)
	}
}

func TestGetNoteHandlerNotFound(t *testing.T) {
	_, r := setupNotesHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestGetNoteHandlerInvalidID(t *testing.T) {
	_, r := setupNotesHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes/abc", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateNoteHandler(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]interface{}{"title": "Original"}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var created storage.Note
	json.Unmarshal(w.Body.Bytes(), &created)

	update := map[string]interface{}{
		"title":   "Updated",
		"content": "New content",
	}
	ub, _ := json.Marshal(update)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/notes/%d", created.ID), bytes.NewReader(ub))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w2.Code, w2.Body.String())
	}

	var updated storage.Note
	json.Unmarshal(w2.Body.Bytes(), &updated)
	if updated.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %q", updated.Title)
	}
}

func TestUpdateNoteHandlerNotFound(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]interface{}{"title": "Ghost"}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notes/999", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestDeleteNoteHandler(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]interface{}{"title": "To Delete"}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var created storage.Note
	json.Unmarshal(w.Body.Bytes(), &created)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/notes/%d", created.ID), nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w2.Code, w2.Body.String())
	}
}

func TestDeleteNoteHandlerNotFound(t *testing.T) {
	_, r := setupNotesHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notes/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestListNotesHandlerWithParentFilter(t *testing.T) {
	_, r := setupNotesHandler(t)

	folder := map[string]interface{}{"title": "Folder", "isFolder": true}
	fb, _ := json.Marshal(folder)

	wf := httptest.NewRecorder()
	reqf, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(fb))
	reqf.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wf, reqf)

	var f storage.Note
	json.Unmarshal(wf.Body.Bytes(), &f)

	child := map[string]interface{}{"title": "Child", "parentId": f.ID}
	cb, _ := json.Marshal(child)

	wc := httptest.NewRecorder()
	reqc, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(cb))
	reqc.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wc, reqc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/notes?parent_id=%d", f.ID), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}


func TestListNotesHandlerPagination(t *testing.T) {
	_, r := setupNotesHandler(t)

	for i := 1; i <= 5; i++ {
		body := map[string]interface{}{"title": fmt.Sprintf("Note %d", i)}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes?offset=0&limit=3", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var notes []storage.Note
	json.NewDecoder(w.Body).Decode(&notes)
	if len(notes) != 3 {
		t.Errorf("expected 3 notes with limit=3, got %d", len(notes))
	}
}

func TestListNotesHandlerOffset(t *testing.T) {
	_, r := setupNotesHandler(t)

	for i := 1; i <= 5; i++ {
		body := map[string]interface{}{"title": fmt.Sprintf("Note %d", i)}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/notes?offset=2&limit=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var notes []storage.Note
	json.NewDecoder(w.Body).Decode(&notes)
	if len(notes) != 3 {
		t.Errorf("expected 3 notes after offset=2 from 5 total, got %d", len(notes))
	}
}

func TestUpdateNoteHandlerInvalidID(t *testing.T) {
	_, r := setupNotesHandler(t)

	body := map[string]string{"title": "test"}
	b, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/notes/invalid", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeleteNoteHandlerInvalidID(t *testing.T) {
	_, r := setupNotesHandler(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/notes/invalid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
