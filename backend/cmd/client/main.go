package main

import (
	"context"
	"log"
	"google.golang.org/grpc"
	pb "github.com/Teav4/vih/backend/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if (err != nil) {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMangaServiceClient(conn)

	// Get manga list
	resp, err := client.GetMangaList(context.Background(), &pb.MangaListRequest{
		Page:     1,
		PageSize: 10,
		SortBy:   "latest",
	})
	if err != nil {
		log.Fatalf("GetMangaList failed: %v", err)
	}

	// Print results
	for _, manga := range resp.Mangas {
		log.Printf("Manga: %s", manga.Title)
	}
}
