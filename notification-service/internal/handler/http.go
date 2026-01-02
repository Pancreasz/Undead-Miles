package handler

import (
	"net/http"

	"github.com/Pancreasz/Undead-Miles/notification/internal/models"
	repository "github.com/Pancreasz/Undead-Miles/notification/internal/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetNotifications(c *gin.Context) {
	userID := c.Param("user_id")

	notes, err := h.repo.GetNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	// Return empty list instead of null if none found
	if notes == nil {
		notes = []models.Notification{} // Empty array
	}

	c.JSON(http.StatusOK, notes)
}
