package seeds

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"backend/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Seeder struct {
	db *mongo.Database
}

func NewSeeder(db *mongo.Database) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) LoadAll(ctx context.Context) error {
	if err := s.LoadMangas(ctx); err != nil {
		return err
	}
	if err := s.LoadChapters(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Seeder) LoadMangas(ctx context.Context) error {
	var mangas []models.Manga
	if err := loadJSONFile("seeds/manga.json", &mangas); err != nil {
		return err
	}
	
	collection := s.db.Collection("mangas")
	for _, manga := range mangas {
		_, err := collection.InsertOne(ctx, manga)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Seeder) LoadChapters(ctx context.Context) error {
	var chapters []models.Chapter
	if err := loadJSONFile("seeds/chapters.json", &chapters); err != nil {
		return err
	}

	collection := s.db.Collection("chapters")
	for _, chapter := range chapters {
		_, err := collection.InsertOne(ctx, chapter)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadJSONFile(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}
