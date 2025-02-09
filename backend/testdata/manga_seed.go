package testdata

import (
	"time"

	"github.com/Teav4/vih/backend/entity"
)

func GetTestMangas() []entity.Manga {
	now := time.Now()
	return []entity.Manga{
		{
			Title:           "One Piece",
			AlternateTitles: []string{"ワンピース", "Đảo Hải Tặc"},
			Description:     "Theo chân cuộc hành trình của Monkey D. Luffy và băng hải tặc Mũ Rơm...",
			Author:          "Oda Eiichiro",
			Artist:          "Oda Eiichiro",
			Genres:          []string{"Action", "Adventure", "Comedy", "Fantasy"},
			Status:          "Ongoing",
			CoverImage:      "https://example.com/onepiece.jpg",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			Title:           "Naruto",
			AlternateTitles: []string{"ナルト"},
			Description:     "Naruto Uzumaki, một ninja trẻ đầy nhiệt huyết...",
			Author:          "Kishimoto Masashi",
			Artist:          "Kishimoto Masashi",
			Genres:          []string{"Action", "Adventure", "Martial Arts"},
			Status:          "Completed",
			CoverImage:      "https://example.com/naruto.jpg",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			Title:           "Dragon Ball",
			AlternateTitles: []string{"ドラゴンボール", "Bi Rồng"},
			Description:     "Cuộc phiêu lưu của Son Goku từ lúc còn nhỏ...",
			Author:          "Toriyama Akira",
			Artist:          "Toriyama Akira",
			Genres:          []string{"Action", "Adventure", "Martial Arts"},
			Status:          "Completed",
			CoverImage:      "https://example.com/dragonball.jpg",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}
}
