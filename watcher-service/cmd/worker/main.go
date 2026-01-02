package main

import (
	"log"

	// Import your internal event package
	"github.com/Pancreasz/Undead-Miles/watcher/internal/event"
)

func main() {
	// 1. Connect to RabbitMQ (Same URL as Marketplace)
	rabbitURL := "amqp://user:password@localhost:5672/"
	rabbitClient, err := event.Connect(rabbitURL)
	if err != nil {
		log.Fatal("Could not connect to RabbitMQ:", err)
	}
	defer rabbitClient.Close()

	// 2. Start the Consumer
	consumer := event.NewConsumer(rabbitClient)

	// 3. Listen to the "trip_created" queue
	// Note: This matches the queue name we declared in the Marketplace
	err = consumer.Listen("trip_created")
	if err != nil {
		log.Println("Error listening to queue:", err)
	}
}
