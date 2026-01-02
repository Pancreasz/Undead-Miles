package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish sends a message to a specific queue
func Publish(conn *amqp.Connection, queueName string, body interface{}) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 1. Declare the queue (idempotent - ensures it exists)
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// 2. Encode the body to JSON
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// 3. Publish the message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})

	if err != nil {
		log.Printf("[Producer] Failed to publish event: %v", err)
		return err
	}

	log.Printf("[Producer] Sent event to '%s'", queueName)
	return nil
}
