package main

import (
	"context"
	"encoding/json"
	"fmt"
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

// Add Manga struct
type Manga struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Categories  []string `json:"categories"`
	Status      string   `json:"status"`
	URL         string   `json:"url"`
}

func NewCrawler() (*Crawler, error) {
	// MongoDB connection
	mongoURI := getEnvOrDefault("MONGODB_URI", "mongodb://root:example@mongodb:27017")
	log.Printf("Connecting to MongoDB: %s", mongoURI)
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Redis connection - use simpler connection approach
	redisAddr := getEnvOrDefault("REDIS_ADDR", "redis:6379")
	log.Printf("Connecting to Redis: %s", redisAddr)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	return &Crawler{
		mongoClient: mongoClient,
		redisClient: redisClient,
		batchSize:   1, // Set to 1 for testing
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
		return fmt.Errorf("failed to marshal progress: %v", err)
	}

	// Add debug logging
	log.Printf("Updating progress: %s", string(data))

	err = c.redisClient.Set(ctx, "crawler:progress", data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set progress in Redis: %v", err)
	}

	// Verify the data was written
	_, err = c.redisClient.Get(ctx, "crawler:progress").Result()
	if err != nil {
		return fmt.Errorf("failed to verify progress in Redis: %v", err)
	}

	return nil
}

// Modify the generateSampleData function
func (c *Crawler) generateSampleData(ctx context.Context) error {
	totalRecords := 1000

	// Initial progress update
	startProgress := Progress{
		TotalItems:     totalRecords,
		ProcessedItems: 0,
		CurrentURL:     "",
		Status:         "starting",
		LastUpdateTime: time.Now().Format(time.RFC3339),
	}
	if err := c.updateProgress(ctx, startProgress); err != nil {
		return fmt.Errorf("failed to update initial progress: %v", err)
	}

	// Fixed categories array
	categories := [][]string{
		{"Action"},
		{"Adventure"},
		{"Romance"},
		{"Comedy"},
		{"Fantasy"},
		{"Magic"},
		{"Sports"},
		{"School"},
		{"Horror"},
		{"Mystery"},
	}

	// Generate 1000 sample manga records
	sampleManga := make([]Manga, totalRecords)
	for i := 0; i < totalRecords; i++ {
		// Pick 2 random categories
		cat1 := categories[i%len(categories)]
		cat2 := categories[(i+1)%len(categories)]
		combinedCats := append([]string{}, cat1...)
		combinedCats = append(combinedCats, cat2...)

		sampleManga[i] = Manga{
			Title:       fmt.Sprintf("Manga %d", i+1),
			Description: fmt.Sprintf("This is description for manga %d", i+1),
			Author:      fmt.Sprintf("Author %d", (i%50)+1),
			Categories:  combinedCats,
			Status:      []string{"Ongoing", "Completed"}[i%2],
			URL:         fmt.Sprintf("http://example.com/manga/%d", i+1),
		}
	}

	// Save to MongoDB
	collection := c.mongoClient.Database("manga_db").Collection("mangas")

	// Process in batches of 50
	batchSize := 50
	for i := 0; i < len(sampleManga); i += batchSize {
		end := i + batchSize
		if end > len(sampleManga) {
			end = len(sampleManga)
		}

		batch := sampleManga[i:end]

		// Convert batch to interface slice for MongoDB
		var documents []interface{}
		for _, manga := range batch {
			documents = append(documents, manga)
		}

		// Insert batch
		_, err := collection.InsertMany(ctx, documents)
		if err != nil {
			log.Printf("❌ Error saving batch %d-%d: %v", i, end, err)
		} else {
			log.Printf("✅ Saved batch %d-%d successfully", i, end)
		}

		// Simulate processing time
		time.Sleep(2 * time.Second)

		// Update progress
		progress := Progress{
			TotalItems:     totalRecords,
			ProcessedItems: end,
			CurrentURL:     batch[len(batch)-1].URL,
			Status:         "processing",
			LastUpdateTime: time.Now().Format(time.RFC3339),
		}

		if err := c.updateProgress(ctx, progress); err != nil {
			log.Printf("❌ Failed to update progress: %v", err)
		} else {
			log.Printf("✅ Updated progress: %d/%d", end, totalRecords)
		}
	}

	// Final progress update
	finalProgress := Progress{
		TotalItems:     totalRecords,
		ProcessedItems: totalRecords,
		CurrentURL:     "",
		Status:         "completed",
		LastUpdateTime: time.Now().Format(time.RFC3339),
	}

	return c.updateProgress(ctx, finalProgress)
}

// Modify the crawlWorker function
func (c *Crawler) crawlWorker(ctx context.Context, workerID int) {
	log.Printf("Worker %d starting...", workerID)

	// Test Redis connection first
	err := c.redisClient.Set(ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		log.Printf("Worker %d Redis connection error: %v", workerID, err)
		return
	}

	// Initial progress
	initialProgress := Progress{
		TotalItems:     2,
		ProcessedItems: 0,
		CurrentURL:     "",
		Status:         fmt.Sprintf("Worker %d starting", workerID),
		LastUpdateTime: time.Now().Format(time.RFC3339),
	}

	if err := c.updateProgress(ctx, initialProgress); err != nil {
		log.Printf("Worker %d failed to update initial progress: %v", workerID, err)
		return
	}

	// For testing, only worker 0 will generate sample data
	if workerID == 0 {
		if err := c.generateSampleData(ctx); err != nil {
			log.Printf("Worker %d error: %v", workerID, err)
		}
	} else {
		log.Printf("Worker %d waiting...", workerID)
		time.Sleep(10 * time.Second)
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
