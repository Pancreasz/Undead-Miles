package main

import (
	"context"
	"log"
	"os"

	"github.com/Pancreasz/Undead-Miles/watcher/internal/database"
	"github.com/Pancreasz/Undead-Miles/watcher/internal/event"
	"github.com/Pancreasz/Undead-Miles/watcher/internal/handler"
	"github.com/gin-gonic/gin" // Import Gin
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// 1. TELEMETRY IMPORTS
	"github.com/Pancreasz/Undead-Miles/watcher/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// ---------------------------------------------------------
	// 2. INITIALIZE JAEGER TRACER
	// ---------------------------------------------------------
	// Service Name: watcher-service
	// Collector URL: jaeger:4318 (HTTP collector inside Docker network)
	shutdown := telemetry.InitTracer("watcher-service", "jaeger:4318")
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	// 3. Database Setup
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:cpre888@localhost:5556/undeadmiles?sslmode=disable")
	connPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer connPool.Close()
	db := database.New(connPool)

	// 4. RabbitMQ Setup
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

	r := gin.Default()

	// 5. ADD TRACING MIDDLEWARE (Must be FIRST)
	// This automatically tracks every incoming HTTP request
	r.Use(otelgin.Middleware("watcher-service"))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
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
