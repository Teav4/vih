package server

import (
	"context"
	pb "vih/backend/proto"
)

type MangaService struct {
	pb.UnimplementedMangaServiceServer
}

func NewMangaService() *MangaService {
	return &MangaService{}
}

func (s *MangaService) GetMangaList(ctx context.Context, req *pb.MangaListRequest) (*pb.MangaListResponse, error) {
	// TODO: Implement manga list retrieval
	return &pb.MangaListResponse{}, nil
}

func (s *MangaService) GetMangaDetail(ctx context.Context, req *pb.MangaDetailRequest) (*pb.MangaDetailResponse, error) {
	// TODO: Implement manga detail retrieval
	return &pb.MangaDetailResponse{}, nil
}

func (s *MangaService) GetChapterList(ctx context.Context, req *pb.ChapterListRequest) (*pb.ChapterListResponse, error) {
	// TODO: Implement chapter list retrieval
	return &pb.ChapterListResponse{}, nil
}

func (s *MangaService) GetChapterContent(ctx context.Context, req *pb.ChapterContentRequest) (*pb.ChapterContentResponse, error) {
	// TODO: Implement chapter content retrieval
	return &pb.ChapterContentResponse{}, nil
}

func (s *MangaService) UpdateViewCount(ctx context.Context, req *pb.UpdateViewCountRequest) (*pb.UpdateViewCountResponse, error) {
	// TODO: Implement view count update
	return &pb.UpdateViewCountResponse{}, nil
}
