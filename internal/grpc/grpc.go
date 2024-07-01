package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"url-shortener/internal/pkg/pb"
	shortenerservice "url-shortener/internal/service"
)

type service interface {
	CreateShortUrl(ctx context.Context, originalURL string) (string, error)
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	GetURLStats(ctx context.Context, shortURL string) (int64, error)
}

type server struct {
	pb.UnimplementedURLShortenerServer

	service service
}

// NewURLServer init new URL service
func NewURLServer(
	s service,
) pb.URLShortenerServer {
	return &server{
		service: s,
	}
}

func (s *server) CreateShortUrl(
	ctx context.Context,
	request *pb.CreateShortUrlRequest,
) (*pb.CreateShortUrlResponse, error) {
	shortURL, err := s.service.CreateShortUrl(ctx, request.OriginalUrl)
	if err != nil {
		if errors.Is(err, shortenerservice.ErrInvalidURL) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid url: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create shorten url: %v", err)
	}

	return &pb.CreateShortUrlResponse{ShortUrl: shortURL}, nil
}

func (s *server) GetOriginalUrl(
	ctx context.Context,
	request *pb.GetOriginalUrlRequest,
) (*pb.GetOriginalUrlResponse, error) {
	originalURL, err := s.service.GetOriginalURL(ctx, request.ShortUrl)
	if err != nil {
		if errors.Is(err, shortenerservice.ErrURLNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, "url not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get original url: %v", err)
	}

	return &pb.GetOriginalUrlResponse{OriginalUrl: originalURL}, nil
}

func (s *server) GetStatistics(
	ctx context.Context,
	request *pb.GetStatisticsRequest,
) (*pb.GetStatisticsResponse, error) {
	requestNumber, err := s.service.GetURLStats(ctx, request.ShortUrl)
	if err != nil {
		if errors.Is(err, shortenerservice.ErrURLNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, "url not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get url stats: %v", err)
	}

	return &pb.GetStatisticsResponse{RequestNumber: requestNumber}, nil
}
