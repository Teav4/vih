package handler

import (
	"net/http"

	"github.com/Teav4/vih/backend/repository"
	"github.com/Teav4/vih/backend/testdata"
)

type SeedHandler struct {
	mangaRepo repository.MangaRepository
}

func NewSeedHandler(mangaRepo repository.MangaRepository) *SeedHandler {
	return &SeedHandler{
		mangaRepo: mangaRepo,
	}
}

func (h *SeedHandler) SeedData(w http.ResponseWriter, r *http.Request) {
	testMangas := testdata.GetTestMangas()

	for _, manga := range testMangas {
		err := h.mangaRepo.Create(r.Context(), &manga)
		if err != nil {
			http.Error(w, "Failed to seed manga data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test data seeded successfully"))
}
