package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/go-redis/redis/v8"
)

type Crawler struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	batchSize   int
	wg          sync.WaitGroup
}

func NewCrawler() (*Crawler, error) {
	// MongoDB connection
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		return nil, err
	}

	// Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

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

func (c *Crawler) crawlWorker(ctx context.Context, workerID int) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Printf("Worker %d is crawling...", workerID)
			// TODO: Add actual crawling logic here
		}
	}
}

func main() {
	crawler, err := NewCrawler()
	if err != nil {
		log.Fatal(err)
	}

	if err := crawler.Start(); err != nil {
		log.Fatal(err)
	}
}
