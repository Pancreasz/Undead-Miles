package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Pancreasz/Undead-Miles/watcher/internal/database" // Update with your actual path
)

type Consumer struct {
	client *RabbitClient
	db     *database.Queries // <--- NEW: Access to the DB
}

// Update constructor to accept DB
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
	msgs, err := c.client.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			// 1. Parse the Trip Data
			var trip tripMessage
			if err := json.Unmarshal(d.Body, &trip); err != nil {
				log.Printf("Error parsing JSON: %v", err)
				continue
			}

			// 2. Check Database for Matches
			// We use a background context here
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			matches, err := c.db.FindMatches(ctx, database.FindMatchesParams{
				Origin:      trip.Origin,
				Destination: trip.Destination,
			})
			cancel() // Clean up context

			if err != nil {
				log.Printf("Error querying database: %v", err)
				continue
			}

			// 3. Log the Result (Simulating sending emails)
			if len(matches) > 0 {
				log.Printf("🔥 FOUND %d MATCHES! Sending alerts...", len(matches))
				for _, watcher := range matches {
					log.Printf("   -> Emailing user: %s", watcher.UserEmail)
				}
			} else {
				log.Println("   -> No watchers found for this route.")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages in %s...", queueName)
	<-forever
	return nil
}
