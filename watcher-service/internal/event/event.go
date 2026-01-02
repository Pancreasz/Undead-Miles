package event

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// FIX: Capitalized 'Conn' so consumer.go can access it
	Conn *amqp.Connection
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

	// Declare the Queue we are listening to
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
		Conn: conn, // Assign to Capital field
		ch:   ch,
	}, nil
}

// Close cleans up connections
func (rc *RabbitClient) Close() {
	rc.ch.Close()
	rc.Conn.Close() // Use Capital field
}
