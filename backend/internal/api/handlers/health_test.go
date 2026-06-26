package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"kvasir/internal/storage"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestHealthReturnsOK(t *testing.T) {
	r := setupRouter()
	h := New(nil)
	r.GET("/health", h.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", body["status"])
	}

	if body["version"] == "" {
		t.Error("expected non-empty version")
	}
}

func TestHealthReturnsJSONHeader(t *testing.T) {
	r := setupRouter()
	h := New(nil)
	r.GET("/health", h.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected json content-type, got %q", contentType)
	}
}

func TestHealthWithStoreSet(t *testing.T) {
	r := setupRouter()
	h := New(&storage.Store{})
	r.GET("/health", h.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNewNilStore(t *testing.T) {
	h := New(nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.Store != nil {
		t.Error("expected nil Store")
	}
}

func TestNewWithStore(t *testing.T) {
	s := &storage.Store{}
	h := New(s)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
	if h.Store != s {
		t.Error("Store pointer mismatch")
	}
}
