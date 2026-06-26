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

		v1.GET("/notes", h.ListNotes)
		v1.POST("/notes", h.CreateNote)
		v1.GET("/notes/:id", h.GetNote)
		v1.PUT("/notes/:id", h.UpdateNote)
		v1.DELETE("/notes/:id", h.DeleteNote)

		v1.GET("/search", h.Search)
		v1.GET("/search/tags", h.SearchByTag)
	}
}
