package postgres

import (
	"context"

	"url-shortener/internal/model"
	storagedb "url-shortener/internal/storage/db"

	"github.com/georgysavva/scany/pgxscan"
)

type URLRepo struct {
	db *storagedb.Database
	tm storagedb.QueryEngineProvider
}

func NewURLRepo(database *storagedb.Database, tm storagedb.QueryEngineProvider) *URLRepo {
	return &URLRepo{
		db: database,
		tm: tm,
	}
}

func (r *URLRepo) AddURL(ctx context.Context, url model.URL) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		r.tm.GetQueryEngine(ctx),
		`INSERT INTO urls(short_url,original_url,created_at) VALUES ($1,$2,$3) RETURNING id;`,
		url.ShortURL,
		url.OriginalURL,
		url.CreatedAt,
	).Scan(&id)

	return id, err
}

func (r *URLRepo) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	if err := r.db.Get(
		ctx,
		r.tm.GetQueryEngine(ctx),
		&originalURL,
		"SELECT original_url FROM urls WHERE short_url=$1 LIMIT 1",
		shortURL,
	); err != nil {
		if pgxscan.NotFound(err) {
			return "", ErrObjectNotFound
		}
		return "", err
	}

	return originalURL, nil
}
