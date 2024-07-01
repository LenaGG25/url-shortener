package repository

import (
	"context"

	"url-shortener/internal/model"
)

type URLRepo interface {
	AddURL(ctx context.Context, url model.URL) (int64, error)
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
}

type URLStatsRepo interface {
	AddURLStats(ctx context.Context, urlStats model.URLStats) (int64, error)
	GetURLStats(ctx context.Context, shortURL string) (int64, error)
	UpdateURLStats(ctx context.Context, shortURL string) error
}

type Redis interface {
	Set(ctx context.Context, shortURL string, originalURL string) error
	Get(ctx context.Context, shortURL string) (string, error)
}
