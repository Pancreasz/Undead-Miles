package repository

import (
	"context"
	"fmt"

	"github.com/Pancreasz/Undead-Miles/notification/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// SaveNotification saves a new match to the database
func (r *Repository) SaveNotification(userID, message string) error {
	query := `INSERT INTO notifications (user_id, message) VALUES ($1, $2)`
	_, err := r.db.Exec(context.Background(), query, userID, message)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}
	return nil
}

// GetNotifications fetches all alerts for a specific user
func (r *Repository) GetNotifications(userID string) ([]models.Notification, error) {
	query := `SELECT id, user_id, message, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		// We scan ID as string (Postgres UUID -> Go string)
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}
