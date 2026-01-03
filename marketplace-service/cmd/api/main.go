package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Pancreasz/Undead-Miles/marketplace/internal/database"
	"github.com/Pancreasz/Undead-Miles/marketplace/internal/event"
	"github.com/Pancreasz/Undead-Miles/marketplace/internal/handler"
	"github.com/gin-contrib/cors" // Gin specific CORS
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// 1. Database
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:cpre888@localhost:5556/undeadmiles?sslmode=disable")
	connPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer connPool.Close()
	db := database.New(connPool)

	// 2. RabbitMQ
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://user:password@localhost:5672/")
	rabbitClient, err := event.Connect(rabbitURL)
	if err != nil {
		log.Fatal("Could not connect to RabbitMQ:", err)
	}
	defer rabbitClient.Close()

	// 3. Initialize Handler
	h := handler.New(db, rabbitClient)

	// 4. Router (Gin)
	r := gin.Default()

	// Basic CORS for Gin
	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Marketplace Service is OK!")
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/trips", h.ListTrips)
	r.POST("/trips", h.CreateTrip)

	// 5. Start Server
	port := "8080"
	fmt.Printf("Marketplace Service running on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
