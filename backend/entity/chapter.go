package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chapter struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MangaID   primitive.ObjectID `bson:"manga_id" json:"mangaId"`
	Number    float64            `bson:"number" json:"number"`
	Title     string             `bson:"title" json:"title"`
	Images    []string           `bson:"images" json:"images"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}
