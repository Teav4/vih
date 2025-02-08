package tests

import (
	"context"
	"testing"
	
	"backend/internal/seeds"
	pb "backend/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetMangaList(t *testing.T) {
	ctx := context.Background()
	client := setupTestClient(t)
	
	resp, err := client.GetMangaList(ctx, &pb.MangaListRequest{
		Page: 1,
		PageSize: 10,
	})
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Mangas, 2)
	assert.Equal(t, "One Piece", resp.Mangas[0].Title)
}

func TestGetMangaDetail(t *testing.T) {
	ctx := context.Background()
	client := setupTestClient(t)
	
	resp, err := client.GetMangaDetail(ctx, &pb.MangaDetailRequest{
		MangaId: "manga-001",
	})
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "One Piece", resp.Manga.Title)
}
