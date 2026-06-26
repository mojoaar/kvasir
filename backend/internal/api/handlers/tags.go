package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kvasir/internal/storage"
)

type createTagRequest struct {
	Name  string `json:"name"  binding:"required"`
	Color string `json:"color" binding:"required"`
}

type updateTagRequest struct {
	Name  string `json:"name"  binding:"required"`
	Color string `json:"color" binding:"required"`
}

type addTagRequest struct {
	TagID int64 `json:"tagId" binding:"required"`
}

// ListTags godoc
// @Summary      List tags
// @Description  Returns all tags ordered by name
// @Tags         tags
// @Produce      json
// @Success      200  {array}   storage.Tag
// @Failure      500  {object}  map[string]string
// @Router       /tags [get]
func (h *Handler) ListTags(c *gin.Context) {
	tags, err := h.Store.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tags"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// CreateTag godoc
// @Summary      Create a tag
// @Description  Creates a new tag with name and color
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        tag  body  createTagRequest  true  "Tag data"
// @Success      201  {object}  storage.Tag
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags [post]
func (h *Handler) CreateTag(c *gin.Context) {
	var req createTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := storage.Tag{
		Name:  req.Name,
		Color: req.Color,
	}
	if err := h.Store.CreateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tag"})
		return
	}
	c.JSON(http.StatusCreated, tag)
}

// GetTag godoc
// @Summary      Get a tag
// @Description  Returns a single tag by ID
// @Tags         tags
// @Produce      json
// @Param        id   path  int  true  "Tag ID"
// @Success      200  {object}  storage.Tag
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /tags/{id} [get]
func (h *Handler) GetTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	tag, err := h.Store.GetTag(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// UpdateTag godoc
// @Summary      Update a tag
// @Description  Updates an existing tag's name and color by ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path  int               true  "Tag ID"
// @Param        tag  body  updateTagRequest  true  "Updated tag data"
// @Success      200  {object}  storage.Tag
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags/{id} [put]
func (h *Handler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	var req updateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := storage.Tag{
		ID:    id,
		Name:  req.Name,
		Color: req.Color,
	}
	if err := h.Store.UpdateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tag"})
		return
	}
	c.JSON(http.StatusOK, tag)
}

// DeleteTag godoc
// @Summary      Delete a tag
// @Description  Deletes a tag and removes it from all notes
// @Tags         tags
// @Produce      json
// @Param        id   path  int  true  "Tag ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags/{id} [delete]
func (h *Handler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	if err := h.Store.DeleteTag(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// AddTagToNote godoc
// @Summary      Add tag to note
// @Description  Adds a tag to a note (idempotent — no error if already tagged)
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path  int           true  "Note ID"
// @Param        tag  body  addTagRequest true  "Tag ID to add"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id}/tags [post]
func (h *Handler) AddTagToNote(c *gin.Context) {
	noteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	var req addTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Store.AddTagToNote(noteID, req.TagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add tag to note"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "added"})
}

// RemoveTagFromNote godoc
// @Summary      Remove tag from note
// @Description  Removes a tag from a note
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path  int           true  "Note ID"
// @Param        tag  body  addTagRequest true  "Tag ID to remove"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id}/tags [delete]
func (h *Handler) RemoveTagFromNote(c *gin.Context) {
	noteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	var req addTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Store.RemoveTagFromNote(noteID, req.TagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove tag from note"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "removed"})
}

// GetNoteTags godoc
// @Summary      Get note tags
// @Description  Returns all tags assigned to a note, ordered by name
// @Tags         tags
// @Produce      json
// @Param        id   path  int  true  "Note ID"
// @Success      200  {array}   storage.Tag
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id}/tags [get]
func (h *Handler) GetNoteTags(c *gin.Context) {
	noteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	tags, err := h.Store.GetNoteTags(noteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get note tags"})
		return
	}
	c.JSON(http.StatusOK, tags)
}
