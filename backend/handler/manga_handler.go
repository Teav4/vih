package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Teav4/vih/backend/entity"
	"github.com/Teav4/vih/backend/repository"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MangaHandler struct {
	mangaRepo   repository.MangaRepository
	redisClient *redis.Client
}

func NewMangaHandler(mangaRepo repository.MangaRepository, redisClient *redis.Client) *MangaHandler {
	return &MangaHandler{
		mangaRepo:   mangaRepo,
		redisClient: redisClient,
	}
}

func (h *MangaHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/manga", h.GetMangas).Methods("GET")
	r.HandleFunc("/api/manga/{id}", h.GetManga).Methods("GET")
	r.HandleFunc("/api/manga", h.CreateManga).Methods("POST")
	r.HandleFunc("/api/manga/{id}", h.UpdateManga).Methods("PUT")
	r.HandleFunc("/api/manga/{id}", h.DeleteManga).Methods("DELETE")
	r.HandleFunc("/api/records", h.GetRecords).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/progress", h.GetProgress).Methods("GET", "OPTIONS")
}

func (h *MangaHandler) GetMangas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get page and limit from query params, default to page 1 and limit 10
	page := 1
	limit := 10

	mangas, err := h.mangaRepo.FindAll(ctx, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(mangas)
}

func (h *MangaHandler) GetManga(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Try to get from Redis first
	cacheKey := "manga:" + id.Hex()
	if cachedData, err := h.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	}

	// If not in cache, get from MongoDB
	manga, err := h.mangaRepo.FindByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Cache the result
	if jsonData, err := json.Marshal(manga); err == nil {
		h.redisClient.Set(ctx, cacheKey, jsonData, 1*time.Hour)
	}

	json.NewEncoder(w).Encode(manga)
}

func (h *MangaHandler) CreateManga(w http.ResponseWriter, r *http.Request) {
	var manga entity.Manga
	if err := json.NewDecoder(r.Body).Decode(&manga); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.mangaRepo.Create(r.Context(), &manga); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(manga)
}

func (h *MangaHandler) UpdateManga(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var manga entity.Manga
	if err := json.NewDecoder(r.Body).Decode(&manga); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manga.ID = id
	if err := h.mangaRepo.Update(r.Context(), &manga); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(manga)
}

func (h *MangaHandler) DeleteManga(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.mangaRepo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MangaHandler) GetRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	page := 1
	pageSize := 10

	// Get query parameters
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	// Get records from repository
	records, err := h.mangaRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"records":   records,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MangaHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	// Get progress from Redis
	progressData, err := h.redisClient.Get(ctx, "crawler:progress").Result()
	if err == redis.Nil {
		// Return default progress when no data exists
		defaultProgress := map[string]interface{}{
			"total_items":      0,
			"processed_items":  0,
			"current_url":      "",
			"status":           "idle",
			"last_update_time": time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(defaultProgress)
		return
	} else if err != nil {
		http.Error(w, "Failed to connect to Redis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(progressData))
}
