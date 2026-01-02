-- name: CreateTrip :one
INSERT INTO trips (origin, destination, driver_id, price_thb, departure_time)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTripById :one
SELECT * FROM trips
WHERE id = $1 LIMIT 1;

-- name: ListTrips :many
SELECT * FROM trips
ORDER BY created_at DESC;

-- name: SearchTrips :many
SELECT * FROM trips
WHERE origin = $1 
AND destination = $2 
AND status = 'OPEN'
ORDER BY price_thb ASC;

-- name: UpdateTripStatus :exec
UPDATE trips
SET status = $2, updated_at = NOW()
WHERE id = $1;