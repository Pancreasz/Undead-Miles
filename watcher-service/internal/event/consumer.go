package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Pancreasz/Undead-Miles/watcher/internal/database"
	// No amqp import needed here, we use the helper
)

type Consumer struct {
	client *RabbitClient
	db     *database.Queries
}

func NewConsumer(client *RabbitClient, db *database.Queries) *Consumer {
	return &Consumer{
		client: client,
		db:     db,
	}
}

type tripMessage struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

func (c *Consumer) Listen(queueName string) error {
	// We use c.client.ch (internal channel) for listening
	msgs, err := c.client.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			var trip tripMessage
			if err := json.Unmarshal(d.Body, &trip); err != nil {
				log.Printf("Error parsing JSON: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			matches, err := c.db.FindMatches(ctx, database.FindMatchesParams{
				Origin:      trip.Origin,
				Destination: trip.Destination,
			})
			cancel()

			if err != nil {
				log.Printf("Error querying database: %v", err)
				continue
			}

			if len(matches) > 0 {
				log.Printf("🔥 FOUND %d MATCHES! Sending events...", len(matches))

				for _, watcher := range matches {
					payload := map[string]interface{}{
						"user_id": watcher.UserEmail,
						"message": "We found a trip from " + trip.Origin + " to " + trip.Destination,
					}

					// FIX: Use .Conn (Capital C) from event.go
					// FIX: Call Publish() which lives in producer.go
					err := Publish(c.client.Conn, "match_found", payload)

					if err != nil {
						log.Printf("Failed to publish match_found: %v", err)
					} else {
						log.Printf(" -> Event sent for user: %s", watcher.UserEmail)
					}
				}
			} else {
				log.Println(" -> No watchers found for this route.")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages in %s...", queueName)
	<-forever
	return nil
}
