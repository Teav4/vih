package server

import (
	"context"
	"encoding/json"
	"net/http"

	pb "github.com/Teav4/vih/backend/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewHTTPServer(grpcPort, httpPort string) *http.Server {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	
	// Register gRPC-Gateway handlers
	err := pb.RegisterMangaServiceHandlerFromEndpoint(
		ctx, 
		mux,
		grpcPort,
		opts,
	)
	if err != nil {
		panic(err)
	}

	return &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}
}
