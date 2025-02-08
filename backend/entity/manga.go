package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manga struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title           string             `bson:"title" json:"title"`
	AlternateTitles []string           `bson:"alternate_titles" json:"alternateTitles"`
	Description     string             `bson:"description" json:"description"`
	Author          string             `bson:"author" json:"author"`
	Artist          string             `bson:"artist" json:"artist"`
	Genres          []string           `bson:"genres" json:"genres"`
	Status          string             `bson:"status" json:"status"` // Ongoing, Completed
	CoverImage      string             `bson:"cover_image" json:"coverImage"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}
