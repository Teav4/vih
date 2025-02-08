package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Teav4/vih/backend/internal/server"
	pb "github.com/Teav4/vih/backend/proto"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":50051"
	httpPort = ":8080"
)

func main() {
	// Create TCP listener for gRPC server
	lis, err := net.Listen("tcp", grpcPort)
	if (err != nil) {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	mangaService := server.NewMangaService()
	pb.RegisterMangaServiceServer(grpcServer, mangaService)

	// Create HTTP server
	httpServer := server.NewHTTPServer(grpcPort, httpPort)

	// Start servers
	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		log.Printf("Starting HTTP server on port %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}
}
