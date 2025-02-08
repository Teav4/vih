package tests

import (
	"context"
	"testing"
	
	"backend/internal/seeds"
	pb "backend/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func setupTestClient(t *testing.T) pb.MangaServiceClient {
	// Connect to test database
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}

	// Clean database
	db := client.Database("manga_test")
	db.Drop(ctx)

	// Load seed data
	seeder := seeds.NewSeeder(db)
	if err := seeder.LoadAll(ctx); err != nil {
		t.Fatal(err)
	}

	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	return pb.NewMangaServiceClient(conn)
}
