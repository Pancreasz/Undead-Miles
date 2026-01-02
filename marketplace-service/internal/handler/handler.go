package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/Pancreasz/Undead-Miles/marketplace/internal/database"
	"github.com/Pancreasz/Undead-Miles/marketplace/internal/event"
	"github.com/gin-gonic/gin" // Import Gin
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	DB     *database.Queries
	Rabbit *event.RabbitClient
}

func New(db *database.Queries, rabbit *event.RabbitClient) *Handler {
	return &Handler{
		DB:     db,
		Rabbit: rabbit,
	}
}

type createTripRequest struct {
	Origin        string    `json:"origin" binding:"required"`
	Destination   string    `json:"destination" binding:"required"`
	DriverID      string    `json:"driver_id" binding:"required"`
	PriceThb      int32     `json:"price_thb" binding:"required"`
	DepartureTime time.Time `json:"departure_time" binding:"required"`
}

// --- Route Functions ---

func (h *Handler) ListTrips(c *gin.Context) {
	trips, err := h.DB.ListTrips(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trips"})
		return
	}
	c.JSON(http.StatusOK, trips)
}

func (h *Handler) CreateTrip(c *gin.Context) {
	var params createTripRequest

	// 1. Parse JSON (Gin style)
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	driverUUID, err := uuid.Parse(params.DriverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Driver UUID"})
		return
	}

	// 2. Database Call
	trip, err := h.DB.CreateTrip(c.Request.Context(), database.CreateTripParams{
		Origin:      params.Origin,
		Destination: params.Destination,
		DriverID:    driverUUID,
		PriceThb:    params.PriceThb,
		DepartureTime: pgtype.Timestamp{
			Time:  params.DepartureTime,
			Valid: true,
		},
	})

	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create trip"})
		return
	}

	// 3. RabbitMQ Publish
	err = h.Rabbit.Publish(c.Request.Context(), trip)
	if err != nil {
		log.Println("WARNING: Failed to publish event:", err)
	}

	c.JSON(http.StatusCreated, trip)
}
