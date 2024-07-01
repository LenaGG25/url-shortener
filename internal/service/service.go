package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"url-shortener/internal/model"
	"url-shortener/internal/storage/db"
	"url-shortener/internal/storage/repository"
	"url-shortener/internal/storage/repository/postgres"
	"url-shortener/internal/storage/repository/redis"
)

type Service struct {
	repoURL      repository.URLRepo
	repoURLStats repository.URLStatsRepo
	tm           *db.TransactionManager
	redis        repository.Redis
	urlGenerator *URLGenerator
}

func New(
	repoURL repository.URLRepo,
	repoURLStats repository.URLStatsRepo,
	tm *db.TransactionManager,
	redis repository.Redis,
	urlGenerator *URLGenerator,
) *Service {
	return &Service{
		repoURL:      repoURL,
		repoURLStats: repoURLStats,
		tm:           tm,
		redis:        redis,
		urlGenerator: urlGenerator,
	}
}

func (s *Service) CreateShortUrl(ctx context.Context, originalURL string) (string, error) {
	if !s.validateURL(originalURL) {
		return "", fmt.Errorf("validate url %v: %w", originalURL, ErrInvalidURL)
	}

	shortURL := s.urlGenerator.GenerateShortURL(originalURL)

	if err := s.tm.RunSerializable(
		ctx,
		func(ctxTX context.Context) error {
			_, repoErr := s.repoURL.GetOriginalURL(ctxTX, shortURL)
			if errors.Is(repoErr, postgres.ErrObjectNotFound) {
				_, repoErr = s.repoURL.AddURL(ctxTX, model.URL{
					ShortURL:    shortURL,
					OriginalURL: originalURL,
					CreatedAt:   time.Now(),
				})
				if repoErr != nil {
					return repoErr
				}

				_, repoErr = s.repoURLStats.AddURLStats(ctxTX, model.URLStats{
					ShortURL:      shortURL,
					RequestNumber: 0,
				})
			}

			return repoErr
		},
	); err != nil {
		return "", err
	}

	if err := s.redis.Set(ctx, shortURL, originalURL); err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *Service) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string

	originalURL, err := s.redis.Get(ctx, shortURL)
	if err != nil {
		if errors.Is(err, redis.ErrFindURL) {
			var repoErr error
			originalURL, repoErr = s.repoURL.GetOriginalURL(ctx, shortURL)
			if repoErr != nil {
				if errors.Is(repoErr, postgres.ErrObjectNotFound) {
					return "", ErrURLNotFound
				}
				return "", repoErr
			}

			if repoErr = s.redis.Set(ctx, shortURL, originalURL); repoErr != nil {
				return "", repoErr
			}
		}

		return "", err
	}

	if err := s.repoURLStats.UpdateURLStats(ctx, shortURL); err != nil {
		if errors.Is(err, postgres.ErrUpdateFailed) {
			return "", ErrURLNotFound
		}
		return "", err
	}

	return originalURL, nil
}

func (s *Service) GetURLStats(ctx context.Context, shortURL string) (int64, error) {
	requestNumbers, err := s.repoURLStats.GetURLStats(ctx, shortURL)
	if err != nil {
		if errors.Is(err, postgres.ErrObjectNotFound) {
			return 0, ErrURLNotFound
		}
		return 0, err
	}

	return requestNumbers, nil
}

func (s *Service) validateURL(originalURL string) bool {
	_, err := url.ParseRequestURI(originalURL)
	return err == nil
}
