package service

import (
	"context"
	"testing"
	"time"

	"github.com/Pancreasz/Undead-Miles/marketplace/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 2. The Mock Repository
type MockTripRepository struct {
	mock.Mock
}

func (m *MockTripRepository) CreateTrip(ctx context.Context, arg database.CreateTripParams) (database.Trip, error) {
	args := m.Called(ctx, arg)
	// Return the object at index 0 (as database.Trip) and error at index 1
	return args.Get(0).(database.Trip), args.Error(1)
}

// --- TEST 1: Validation Failure (Invalid Price) ---
func TestCreateTrip_InvalidPrice(t *testing.T) {
	// Setup
	mockRepo := new(MockTripRepository)
	service := NewTripService(mockRepo, nil) // Assuming 2nd arg is Queue/Events

	// Prepare invalid data (Price is negative)
	params := database.CreateTripParams{
		DriverID:    uuid.New(),
		Origin:      "Bangkok",
		Destination: "Pattaya",
		PriceThb:    -500, // INVALID!
		DepartureTime: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}

	// Execute
	_, err := service.CreateTrip(context.Background(), params)

	// Verify
	assert.Error(t, err)
	assert.Equal(t, "price must be greater than zero", err.Error())

	// Ensure the repo was NEVER called because validation failed first
	mockRepo.AssertNotCalled(t, "CreateTrip")
}

// --- TEST 2: Success Case ---
func TestCreateTrip_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockTripRepository)
	service := NewTripService(mockRepo, nil)

	// Valid Data
	validID := uuid.New()
	params := database.CreateTripParams{
		DriverID:    validID,
		Origin:      "Bangkok",
		Destination: "Chiang Mai",
		PriceThb:    1500,
		DepartureTime: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}

	// What we expect the DB to return
	expectedTrip := database.Trip{
		ID:            uuid.New(),
		DriverID:      validID,
		Origin:        "Bangkok",
		Destination:   "Chiang Mai",
		PriceThb:      1500,
		Status:        "OPEN",
		DepartureTime: params.DepartureTime,
	}

	// Mock Expectation: "When CreateTrip is called with 'params', return 'expectedTrip'"
	mockRepo.On("CreateTrip", mock.Anything, params).Return(expectedTrip, nil)

	// Execute
	result, err := service.CreateTrip(context.Background(), params)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedTrip.ID, result.ID)
	assert.Equal(t, "OPEN", result.Status)

	mockRepo.AssertExpectations(t)
}
