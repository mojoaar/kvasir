package api

import (
	"github.com/gin-gonic/gin"
	"kvasir/internal/api/handlers"
	"kvasir/internal/api/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", handlers.Health)
	}
}
