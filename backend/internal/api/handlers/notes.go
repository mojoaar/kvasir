package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kvasir/internal/storage"
)

type createNoteRequest struct {
	Title     string `json:"title"     binding:"required"`
	Content   string `json:"content"`
	VaultID   *int64 `json:"vaultId"`
	ParentID  *int64 `json:"parentId"`
	IsFolder  bool   `json:"isFolder"`
	SortOrder int    `json:"sortOrder"`
}

type updateNoteRequest struct {
	Title     string `json:"title"     binding:"required"`
	Content   string `json:"content"`
	VaultID   *int64 `json:"vaultId"`
	ParentID  *int64 `json:"parentId"`
	IsFolder  bool   `json:"isFolder"`
	SortOrder int    `json:"sortOrder"`
}

// ListNotes godoc
// @Summary      List notes
// @Description  Returns a paginated list of notes, optionally filtered by vault or parent. Folders are listed first.
// @Tags         notes
// @Produce      json
// @Param        offset     query     int   false  "Pagination offset"
// @Param        limit      query     int   false  "Page size (default 50)"
// @Param        vault_id   query     int   false  "Filter by vault ID"
// @Param        parent_id  query     int   false  "Filter by parent folder ID"
// @Success      200  {array}   storage.Note
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes [get]
func (h *Handler) ListNotes(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	var vaultID *int64
	var parentID *int64

	if v := c.Query("vault_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vault_id"})
			return
		}
		vaultID = &id
	}

	if p := c.Query("parent_id"); p != "" {
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_id"})
			return
		}
		parentID = &id
	}

	notes, err := h.Store.ListNotes(vaultID, parentID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// CreateNote godoc
// @Summary      Create a note
// @Description  Creates a new note or folder
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        note  body  createNoteRequest  true  "Note data"
// @Success      201   {object}  storage.Note
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes [post]
func (h *Handler) CreateNote(c *gin.Context) {
	var req createNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := storage.Note{
		Title:     req.Title,
		Content:   req.Content,
		VaultID:   req.VaultID,
		ParentID:  req.ParentID,
		IsFolder:  req.IsFolder,
		SortOrder: req.SortOrder,
	}

	if err := h.Store.CreateNote(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetNote godoc
// @Summary      Get a note
// @Description  Returns a single note by ID
// @Tags         notes
// @Produce      json
// @Param        id   path  int  true  "Note ID"
// @Success      200  {object}  storage.Note
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /notes/{id} [get]
func (h *Handler) GetNote(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	note, err := h.Store.GetNote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// UpdateNote godoc
// @Summary      Update a note
// @Description  Updates an existing note by ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id    path  int                true  "Note ID"
// @Param        note  body  updateNoteRequest  true  "Updated note data"
// @Success      200   {object}  storage.Note
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes/{id} [put]
func (h *Handler) UpdateNote(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	var req updateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := storage.Note{
		ID:        id,
		Title:     req.Title,
		Content:   req.Content,
		VaultID:   req.VaultID,
		ParentID:  req.ParentID,
		IsFolder:  req.IsFolder,
		SortOrder: req.SortOrder,
	}

	if err := h.Store.UpdateNote(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote godoc
// @Summary      Delete a note
// @Description  Soft-deletes a note by ID
// @Tags         notes
// @Produce      json
// @Param        id   path  int  true  "Note ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (h *Handler) DeleteNote(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	if err := h.Store.SoftDeleteNote(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
