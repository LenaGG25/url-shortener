package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"url-shortener/internal/pkg/pb"
)

type Service struct {
	grpcClient     pb.URLShortenerClient
	shortURLPrefix string
}

func New(
	grpcClient pb.URLShortenerClient,
	shortURLPrefix string,
) *Service {
	return &Service{
		grpcClient:     grpcClient,
		shortURLPrefix: shortURLPrefix,
	}
}

func (s *Service) CreateShortUrl(ctx context.Context, originalURL string) (string, error) {
	response, err := s.grpcClient.CreateShortUrl(ctx, &pb.CreateShortUrlRequest{OriginalUrl: originalURL})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			return "", err
		}

		switch status.Code() {
		case codes.InvalidArgument:
			return "", ErrInvalidURL
		default:
			return "", err
		}
	}

	return fmt.Sprintf("%s/short_url/%s", s.shortURLPrefix, response.ShortUrl), nil
}

func (s *Service) GetOriginalUrl(ctx context.Context, shortURL string) (string, error) {
	response, err := s.grpcClient.GetOriginalUrl(ctx, &pb.GetOriginalUrlRequest{ShortUrl: shortURL})
	if err != nil {
		responseStatus, ok := status.FromError(err)
		if !ok {
			return "", err
		}

		switch responseStatus.Code() {
		case codes.InvalidArgument:
			return "", ErrURLNotFound
		default:
			return "", err
		}
	}

	return response.OriginalUrl, nil
}

func (s *Service) GetStatistics(ctx context.Context, shortURL string) (int64, error) {
	response, err := s.grpcClient.GetStatistics(ctx, &pb.GetStatisticsRequest{ShortUrl: shortURL})
	if err != nil {
		responseStatus, ok := status.FromError(err)
		if !ok {
			return 0, err
		}

		switch responseStatus.Code() {
		case codes.InvalidArgument:
			return 0, ErrURLNotFound
		default:
			return 0, err
		}
	}

	return response.RequestNumber, nil
}
