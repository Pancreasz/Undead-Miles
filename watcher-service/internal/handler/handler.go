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

type createWatcherRequest struct {
	UserEmail   string `json:"user_email" binding:"required"` // Gin validates this automatically
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

// Note: Function signature changes to take *gin.Context
func (h *Handler) CreateWatcher(c *gin.Context) {
	var params createWatcherRequest

	// 1. Parse & Validate JSON (One-liner in Gin)
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call Database
	watcher, err := h.DB.CreateWatcher(c.Request.Context(), database.CreateWatcherParams{
		UserEmail:   params.UserEmail,
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
