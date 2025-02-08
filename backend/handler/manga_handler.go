package handler

import (
	"net/http"

	"github.com/Teav4/vih/backend/repository"

	"github.com/gorilla/mux"
)

type MangaHandler struct {
	mangaRepo repository.MangaRepository
}

func NewMangaHandler(mangaRepo repository.MangaRepository) *MangaHandler {
	return &MangaHandler{mangaRepo: mangaRepo}
}

func (h *MangaHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/manga", h.GetMangas).Methods("GET")
	r.HandleFunc("/api/manga/{id}", h.GetManga).Methods("GET")
	r.HandleFunc("/api/manga", h.CreateManga).Methods("POST")
	r.HandleFunc("/api/manga/{id}", h.UpdateManga).Methods("PUT")
	r.HandleFunc("/api/manga/{id}", h.DeleteManga).Methods("DELETE")
}

func (h *MangaHandler) GetMangas(w http.ResponseWriter, r *http.Request) {
	// Implementation
}

// Implement other handler methods
