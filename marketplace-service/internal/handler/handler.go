package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Pancreasz/Undead-Miles/marketplace/internal/database"
	"github.com/Pancreasz/Undead-Miles/marketplace/internal/event"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// The Struct holds our dependencies
type Handler struct {
	DB     *database.Queries
	Rabbit *event.RabbitClient
}

// New creates a new Handler instance
func New(db *database.Queries, rabbit *event.RabbitClient) *Handler {
	return &Handler{
		DB:     db,
		Rabbit: rabbit,
	}
}

// --- Request Structs ---
type createTripRequest struct {
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DriverID      string    `json:"driver_id"`
	PriceThb      int32     `json:"price_thb"`
	DepartureTime time.Time `json:"departure_time"`
}

// --- Route Functions ---

func (h *Handler) ListTrips(w http.ResponseWriter, r *http.Request) {
	trips, err := h.DB.ListTrips(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch trips", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 200, trips)
}

func (h *Handler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON
	decoder := json.NewDecoder(r.Body)
	params := createTripRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithJSON(w, 400, map[string]string{"error": "Invalid JSON format"})
		return
	}

	driverUUID, err := uuid.Parse(params.DriverID)
	if err != nil {
		respondWithJSON(w, 400, map[string]string{"error": "Invalid Driver UUID"})
		return
	}

	// 2. Database Call
	trip, err := h.DB.CreateTrip(r.Context(), database.CreateTripParams{
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
		respondWithJSON(w, 500, map[string]string{"error": "Could not create trip"})
		return
	}

	// 3. RabbitMQ Publish
	err = h.Rabbit.Publish(r.Context(), trip)
	if err != nil {
		log.Println("WARNING: Failed to publish event:", err)
	}

	respondWithJSON(w, 201, trip)
}

// --- Helper ---
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
