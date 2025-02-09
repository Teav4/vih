package handlers

import (
	"net/http"

	"github.com/go-redis/redis/v8"
)

type Handler struct {
	redisClient *redis.Client
}

func (h *Handler) GetProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get progress from Redis
	progressData, err := h.redisClient.Get(ctx, "crawler:progress").Result()
	if err != nil {
		http.Error(w, "Failed to get progress", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(progressData))
}
