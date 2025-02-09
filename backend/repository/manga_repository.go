package repository

import (
	"context"

	"github.com/Teav4/vih/backend/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *mangaRepository) FindAll(ctx context.Context, page, limit int) ([]entity.Manga, error) {
	collection := r.db.Collection("mangas")
	skip := int64((page - 1) * limit)

	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mangas []entity.Manga
	if err := cursor.All(ctx, &mangas); err != nil {
		return nil, err
	}

	return mangas, nil
}

func (r *mangaRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Manga, error) {
	collection := r.db.Collection("mangas")

	var manga entity.Manga
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&manga)
	if err != nil {
		return nil, err
	}

	return &manga, nil
}

func (r *mangaRepository) Create(ctx context.Context, manga *entity.Manga) error {
	collection := r.db.Collection("mangas")

	result, err := collection.InsertOne(ctx, manga)
	if err != nil {
		return err
	}

	manga.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *mangaRepository) Update(ctx context.Context, manga *entity.Manga) error {
	collection := r.db.Collection("mangas")

	filter := bson.M{"_id": manga.ID}
	update := bson.M{"$set": manga}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mangaRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection("mangas")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
