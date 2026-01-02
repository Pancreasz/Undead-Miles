package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	// Import your internal packages
	"github.com/Pancreasz/Undead-Miles/marketplace/internal/database"
	"github.com/Pancreasz/Undead-Miles/marketplace/internal/event"
	"github.com/Pancreasz/Undead-Miles/marketplace/internal/handler"
)

func main() {
	// 1. Database
	dbURL := "postgres://user:password@localhost:5432/deadahead?sslmode=disable"
	connPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer connPool.Close()
	db := database.New(connPool)

	// 2. RabbitMQ
	rabbitURL := "amqp://user:password@localhost:5672/"
	rabbitClient, err := event.Connect(rabbitURL)
	if err != nil {
		log.Fatal("Could not connect to RabbitMQ:", err)
	}
	defer rabbitClient.Close()

	// 3. Initialize Handler (Inject dependencies)
	h := handler.New(db, rabbitClient)

	// 4. Router
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Marketplace Service is OK!"))
	})

	// Use the methods from the new handler package
	r.Get("/trips", h.ListTrips)
	r.Post("/trips", h.CreateTrip)

	// 5. Start Server
	port := "8080"
	fmt.Printf("Marketplace Service running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
