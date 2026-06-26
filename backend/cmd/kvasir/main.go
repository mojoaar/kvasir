package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"kvasir/internal/api"
)

func main() {
	port := os.Getenv("KVASIR_PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	api.RegisterRoutes(r)

	log.Printf("Kvasir server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
