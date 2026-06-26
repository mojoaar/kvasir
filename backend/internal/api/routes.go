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

		v1.GET("/tags", h.ListTags)
		v1.POST("/tags", h.CreateTag)
		v1.GET("/tags/:id", h.GetTag)
		v1.PUT("/tags/:id", h.UpdateTag)
		v1.DELETE("/tags/:id", h.DeleteTag)

		v1.GET("/notes/:id/tags", h.GetNoteTags)
		v1.POST("/notes/:id/tags", h.AddTagToNote)
		v1.DELETE("/notes/:id/tags", h.RemoveTagFromNote)
	}
}
