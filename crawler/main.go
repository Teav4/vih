package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Crawler struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	batchSize   int
	wg          sync.WaitGroup
}

type Progress struct {
	TotalItems     int    `json:"total_items"`
	ProcessedItems int    `json:"processed_items"`
	CurrentURL     string `json:"current_url"`
	Status         string `json:"status"`
	LastUpdateTime string `json:"last_update_time"`
}

func NewCrawler() (*Crawler, error) {
	// MongoDB connection
	mongoURI := getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Redis connection
	redisURI := getEnvOrDefault("REDIS_URI", "redis://localhost:6379")
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(opt)

	return &Crawler{
		mongoClient: mongoClient,
		redisClient: redisClient,
		batchSize:   5,
	}, nil
}

func (c *Crawler) Start() error {
	ctx := context.Background()

	// Use WaitGroup to manage goroutines
	c.wg.Add(c.batchSize)

	for i := 0; i < c.batchSize; i++ {
		go func(workerID int) {
			defer c.wg.Done()
			c.crawlWorker(ctx, workerID)
		}(i)
	}

	// Wait for all workers to complete
	c.wg.Wait()
	return nil
}

func (c *Crawler) updateProgress(ctx context.Context, progress Progress) error {
	data, err := json.Marshal(progress)
	if err != nil {
		return err
	}
	return c.redisClient.Set(ctx, "crawler:progress", data, 0).Err()
}

func (c *Crawler) crawlWorker(ctx context.Context, workerID int) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			progress := Progress{
				TotalItems:     100, // Update with actual total
				ProcessedItems: 0,   // Update as items are processed
				CurrentURL:     "current_url_here",
				Status:         "running",
				LastUpdateTime: time.Now().Format(time.RFC3339),
			}

			if err := c.updateProgress(ctx, progress); err != nil {
				log.Printf("Error updating progress: %v", err)
			}

			log.Printf("Worker %d progress: %+v", workerID, progress)
		}
	}
}

func main() {
	log.Println("Starting crawler service...")

	crawler, err := NewCrawler()
	if err != nil {
		log.Fatal(err)
	}

	if err := crawler.Start(); err != nil {
		log.Fatal(err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
