-- name: CreateWatcher :one
INSERT INTO watchers (user_email, origin, destination)
VALUES ($1, $2, $3)
RETURNING *;

-- name: FindMatches :many
SELECT * FROM watchers 
WHERE origin = $1 
AND destination = $2;