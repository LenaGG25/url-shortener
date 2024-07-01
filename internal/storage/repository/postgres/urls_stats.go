package postgres

import (
	"context"

	"url-shortener/internal/model"
	db3 "url-shortener/internal/storage/db"

	"github.com/georgysavva/scany/pgxscan"
)

type URLStatsRepo struct {
	db *db3.Database
	tm db3.QueryEngineProvider
}

func NewURLStatsRepo(database *db3.Database, tm db3.QueryEngineProvider) *URLStatsRepo {
	return &URLStatsRepo{
		db: database,
		tm: tm,
	}
}

func (r *URLStatsRepo) AddURLStats(ctx context.Context, urlStats model.URLStats) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		r.tm.GetQueryEngine(ctx),
		`INSERT INTO url_stats(short_url,request_number) VALUES ($1,$2) RETURNING id;`,
		urlStats.ShortURL,
		urlStats.RequestNumber,
	).Scan(&id)

	return id, err
}

func (r *URLStatsRepo) GetURLStats(ctx context.Context, shortURL string) (int64, error) {
	var requestNumber int64
	if err := r.db.Get(
		ctx,
		r.tm.GetQueryEngine(ctx),
		&requestNumber,
		"SELECT request_number FROM url_stats WHERE short_url=$1",
		shortURL,
	); err != nil {
		if pgxscan.NotFound(err) {
			return 0, ErrObjectNotFound
		}
		return 0, err
	}

	return requestNumber, nil
}

func (r *URLStatsRepo) UpdateURLStats(ctx context.Context, shortURL string) error {
	comTag, err := r.db.Exec(
		ctx,
		r.tm.GetQueryEngine(ctx),
		"UPDATE url_stats SET request_number = request_number + 1 WHERE short_url = $1",
		shortURL,
	)
	if err != nil {
		return err
	}

	if comTag.RowsAffected() == 0 {
		return ErrUpdateFailed
	}

	return nil
}
