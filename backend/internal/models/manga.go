package models

import "time"

type Manga struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	CoverURL    string    `json:"cover_url" bson:"cover_url"`
	Genres      []string  `json:"genres" bson:"genres"`
	Status      string    `json:"status" bson:"status"`
	Author      string    `json:"author" bson:"author"`
	ViewCount   int64     `json:"view_count" bson:"view_count"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type Chapter struct {
	ID            string    `json:"id" bson:"_id"`
	MangaID       string    `json:"manga_id" bson:"manga_id"`
	Title         string    `json:"title" bson:"title"` 
	ChapterNumber float32   `json:"chapter_number" bson:"chapter_number"`
	ViewCount     int64     `json:"view_count" bson:"view_count"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

type ChapterContent struct {
	ID         string   `json:"id" bson:"_id"`
	ChapterID  string   `json:"chapter_id" bson:"chapter_id"`
	ImageURLs  []string `json:"image_urls" bson:"image_urls"`
}
