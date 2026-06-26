// @title           Kvasir API
// @version         0.1.0
// @description     Beautiful, techy, Nordic-inspired markdown knowledge base. Sync-first, API-first, plugin-extensible.
// @host            localhost:8080
// @BasePath        /api/v1

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"kvasir/internal/api"
	embed "kvasir/internal/embed"
	"kvasir/internal/storage"
)

func main() {
	port := os.Getenv("KVASIR_PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("KVASIR_DB_PATH")
	if dbPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get home directory: %v", err)
		}
		dbPath = filepath.Join(home, ".kvasir", "kvasir.db")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	store, err := storage.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer store.Close()

	if err := store.SeedIfEmpty(); err != nil {
		log.Printf("welcome note seed failed (non-fatal): %v", err)
	}

	r := gin.Default()
	api.RegisterRoutes(r, store)

	r.NoRoute(gin.WrapH(http.FileServer(embed.DistFS)))

	log.Printf("Kvasir server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
