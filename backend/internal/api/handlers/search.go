package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Search godoc
// @Summary      Full-text search
// @Description  Searches notes by title and content using FTS5. Returns ranked results with snippets.
// @Tags         search
// @Produce      json
// @Param        q      query  string  true   "Search query"
// @Param        limit  query  int     false  "Max results (default 20)"
// @Success      200    {array}   object
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /search [get]
func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}

	results, err := h.Store.Search(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// SearchByTag godoc
// @Summary      Search by tag
// @Description  Finds notes that have tags matching the query
// @Tags         search
// @Produce      json
// @Param        q      query  string  true  "Tag name search query"
// @Success      200    {array}   object
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /search/tags [get]
func (h *Handler) SearchByTag(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results, err := h.Store.SearchByTag(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
