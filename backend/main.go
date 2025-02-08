package main

import (
	"context"
	"log"
	"net/http"

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
	mangaHandler := handler.NewMangaHandler(mangaRepo)

	// Router setup
	r := mux.NewRouter()
	mangaHandler.RegisterRoutes(r)

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
