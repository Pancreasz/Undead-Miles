package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// Connect establishes the TCP connection to RabbitMQ
func Connect(url string) (*RabbitClient, error) {
	var conn *amqp.Connection
	var err error

	// Try to connect 15 times (30 seconds total)
	counts := 0
	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("Connected to RabbitMQ!")
			break
		}

		counts++
		log.Printf("Failed to connect to RabbitMQ (Attempt %d/15). Retrying in 2s...", counts)

		if counts > 15 {
			return nil, fmt.Errorf("could not connect to RabbitMQ after multiple retries: %v", err)
		}

		time.Sleep(2 * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare the Queue
	_, err = ch.QueueDeclare(
		"trip_created", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

// Close cleans up connections
func (rc *RabbitClient) Close() {
	rc.ch.Close()
	rc.conn.Close()
}

// Publish sends a JSON message to the queue
func (rc *RabbitClient) Publish(ctx context.Context, tripData interface{}) error {
	// 1. Convert struct to JSON bytes
	body, err := json.Marshal(tripData)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 2. Publish to the queue
	// Context with timeout to prevent hanging forever
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = rc.ch.PublishWithContext(ctx,
		"",             // exchange
		"trip_created", // routing key (queue name)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish to RabbitMQ: %w", err)
	}

	log.Println("Successfully published event to RabbitMQ!")
	return nil
}
