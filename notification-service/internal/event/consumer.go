package event

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Pancreasz/Undead-Miles/notification/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	client *RabbitClient
	repo   *repository.Repository
}

func NewConsumer(client *RabbitClient, repo *repository.Repository) *Consumer {
	return &Consumer{
		client: client,
		repo:   repo,
	}
}

type MatchFoundEvent struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func (c *Consumer) Listen(queueName string) error {
	msgs, err := c.client.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received match_found: %s", d.Body)

			var event MatchFoundEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error parsing JSON: %v", err)
				continue
			}

			// Save to Database
			err := c.repo.SaveNotification(event.UserID, event.Message)
			if err != nil {
				log.Printf("Error saving to DB: %v", err)
			} else {
				log.Printf("Saved notification for user: %s", event.UserID)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages in %s...", queueName)
	<-forever
	return nil
}

// RabbitClient helper (same as other services)
type RabbitClient struct {
	Conn *amqp.Connection
	ch   *amqp.Channel
}

func Connect(url string) (*RabbitClient, error) {
	var conn *amqp.Connection
	var err error

	counts := 0
	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		counts++
		log.Printf("Failed to connect (Attempt %d/15)...", counts)
		if counts > 15 {
			return nil, err
		}
		time.Sleep(2 * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare the queue we want to listen to
	_, err = ch.QueueDeclare("match_found", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &RabbitClient{Conn: conn, ch: ch}, nil
}

func (rc *RabbitClient) Close() {
	rc.ch.Close()
	rc.Conn.Close()
}
