package repository

import (
	"context"

	"github.com/Teav4/vih/backend/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MangaRepository interface {
	FindAll(ctx context.Context, page, limit int) ([]entity.Manga, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Manga, error)
	Create(ctx context.Context, manga *entity.Manga) error
	Update(ctx context.Context, manga *entity.Manga) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type mangaRepository struct {
	db *mongo.Database
}

func NewMangaRepository(db *mongo.Database) MangaRepository {
	return &mangaRepository{db: db}
}

// Implement the interface methods here
