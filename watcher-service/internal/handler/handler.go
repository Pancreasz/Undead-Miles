package handler

import (
	"net/http"

	"github.com/Pancreasz/Undead-Miles/watcher/internal/database"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *database.Queries
}

func New(db *database.Queries) *Handler {
	return &Handler{DB: db}
}

// FIX: Changed UserEmail to UserID and updated JSON tag to "user_id"
type createWatcherRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

func (h *Handler) CreateWatcher(c *gin.Context) {
	var params createWatcherRequest

	// 1. Parse & Validate JSON
	// This will now look for "user_id" in the input instead of "user_email"
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call Database
	// Note: We map params.UserID to the DB column 'UserEmail'.
	// (We do this to avoid having to wipe the database and regenerate SQL code)
	watcher, err := h.DB.CreateWatcher(c.Request.Context(), database.CreateWatcherParams{
		UserEmail:   params.UserID,
		Origin:      params.Origin,
		Destination: params.Destination,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create watcher"})
		return
	}

	// 3. Respond
	c.JSON(http.StatusCreated, watcher)
}
