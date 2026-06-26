package api

import (
	"github.com/gin-gonic/gin"
	"kvasir/internal/api/handlers"
	"kvasir/internal/api/middleware"
	"kvasir/internal/storage"
)

func RegisterRoutes(r *gin.Engine, store *storage.Store) {
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	h := handlers.New(store)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", h.Health)
	}
}
