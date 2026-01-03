package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Pancreasz/Undead-Miles/notification/internal/event"
	"github.com/Pancreasz/Undead-Miles/notification/internal/handler"
	"github.com/Pancreasz/Undead-Miles/notification/internal/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 1. Database Connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:cpre888@localhost:5432/undeadmiles?sslmode=disable"
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer dbPool.Close()

	// Initialize Repository
	repo := repository.New(dbPool)

	// 2. RabbitMQ Connection
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://user:password@localhost:5672/"
	}

	rabbitClient, err := event.Connect(rabbitURL)
	if err != nil {
		log.Fatal("Could not connect to RabbitMQ:", err)
	}
	defer rabbitClient.Close()

	// 3. Start Consumer (Background Worker)
	// This listens for "match_found" events and saves them using the repo
	consumer := event.NewConsumer(rabbitClient, repo)
	go func() {
		log.Println("Notification Consumer starting...")
		if err := consumer.Listen("match_found"); err != nil {
			log.Println("Consumer failed:", err)
		}
	}()

	// 4. HTTP Server (Gin)
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(cors.Default())

	h := handler.New(repo)

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Notification Service is OK!")
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// API to get notifications for a user (Frontend will use this)
	r.GET("/notifications/:user_id", h.GetNotifications)

	port := "8082"
	log.Printf("Notification Service running on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
