package service

import (
	"context"
	"errors"

	"github.com/Pancreasz/Undead-Miles/marketplace/internal/database"
)

// 1. The Repository Interface
// This tells the code: "I need a database that can do these things."
// We define this interface so we can swap the real DB with a Mock DB during tests.
type TripRepository interface {
	CreateTrip(ctx context.Context, arg database.CreateTripParams) (database.Trip, error)
}

// 2. The Service Struct
// This holds your dependencies (like the database and the queue).
type TripService struct {
	repo  TripRepository
	queue any // We use 'any' for now to keep it simple since we aren't using the queue in this test
}

// 3. The Constructor (NewTripService)
// This is the function your test was looking for!
func NewTripService(repo TripRepository, queue any) *TripService {
	return &TripService{
		repo:  repo,
		queue: queue,
	}
}

// 4. The Logic (CreateTrip)
// This is where your business rules live.
func (s *TripService) CreateTrip(ctx context.Context, arg database.CreateTripParams) (database.Trip, error) {
	// --- VALIDATION CHECK ---
	if arg.PriceThb <= 0 {
		return database.Trip{}, errors.New("price must be greater than zero")
	}
	// ------------------------

	// If validation passes, save to database
	return s.repo.CreateTrip(ctx, arg)
}
