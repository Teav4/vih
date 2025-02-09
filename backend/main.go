package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Teav4/vih/backend/handler"
	"github.com/Teav4/vih/backend/repository"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := mongoClient.Database("manga_db")

	// Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Initialize repositories
	mangaRepo := repository.NewMangaRepository(db)

	// Initialize handlers
	mangaHandler := handler.NewMangaHandler(mangaRepo, redisClient)
	seedHandler := handler.NewSeedHandler(mangaRepo)

	// Router setup
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	// Register routes
	api.HandleFunc("/records", mangaHandler.GetRecords).Methods("GET", "OPTIONS")
	api.HandleFunc("/progress", mangaHandler.GetProgress).Methods("GET", "OPTIONS")

	// Development only routes
	if os.Getenv("ENV") != "production" {
		r.HandleFunc("/dev/seed", seedHandler.SeedData).Methods("POST")
	}

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
