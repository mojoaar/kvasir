package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kvasir/internal/storage"
)

type Handler struct {
	Store *storage.Store
}

func New(store *storage.Store) *Handler {
	return &Handler{Store: store}
}

// Health godoc
// @Summary      Health check
// @Description  Returns the service status and version
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "0.1.0",
	})
}
