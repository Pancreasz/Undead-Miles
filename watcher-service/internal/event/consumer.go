package event

import (
	"encoding/json"
	"log"
)

// Consumer is a struct that listens for messages
type Consumer struct {
	client *RabbitClient
}

func NewConsumer(client *RabbitClient) *Consumer {
	return &Consumer{
		client: client,
	}
}

// Listen starts a loop that reads from the queue
func (c *Consumer) Listen(queueName string) error {
	msgs, err := c.client.ch.Consume(
		queueName, // queue
		"",        // consumer tag
		true,      // auto-ack (true means we tell RabbitMQ "Got it!" immediately)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// This is a "Blocking Channel" - it keeps the program alive
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// This code runs every time a message arrives!
			log.Printf("Received a message: %s", d.Body)

			// TODO: Add logic here to match against Database Alerts
			var tripData map[string]interface{}
			if err := json.Unmarshal(d.Body, &tripData); err == nil {
				log.Printf("New Trip Detected from %s to %s", tripData["Origin"], tripData["Destination"])
			}
		}
	}()

	log.Printf(" [*] Waiting for messages in %s. To exit press CTRL+C", queueName)
	<-forever // Wait here forever

	return nil
}
