package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin" // Import Gin
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Pancreasz/Undead-Miles/watcher/internal/database"
	"github.com/Pancreasz/Undead-Miles/watcher/internal/event"
	"github.com/Pancreasz/Undead-Miles/watcher/internal/handler"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// 1. Database Setup
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:cpre888@localhost:5555/undeadmiles?sslmode=disable")
	connPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer connPool.Close()
	db := database.New(connPool)

	// 2. RabbitMQ Setup
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://user:password@localhost:5672/")
	rabbitClient, err := event.Connect(rabbitURL)
	if err != nil {
		log.Fatal("Could not connect to RabbitMQ:", err)
	}
	defer rabbitClient.Close()

	// ---------------------------------------------------------
	// TASK A: Start the HTTP Server (Gin)
	// ---------------------------------------------------------
	h := handler.New(db)

	// Set Gin to Release mode to quiet down logs (Optional)
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default() // Creates a router with default middleware (logger, recovery)
	r.POST("/watchers", h.CreateWatcher)

	// Run Gin in a Goroutine (Port 8081)
	go func() {
		port := "8081"
		log.Printf("Watcher API (Gin) running on port %s...", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start Gin server:", err)
		}
	}()

	// ---------------------------------------------------------
	// TASK B: Start the RabbitMQ Consumer
	// ---------------------------------------------------------
	consumer := event.NewConsumer(rabbitClient, db)
	log.Println("Watcher Consumer starting...")

	err = consumer.Listen("trip_created")
	if err != nil {
		log.Println("Error listening to queue:", err)
	}
}
