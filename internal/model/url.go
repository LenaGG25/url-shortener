package model

import "time"

type URL struct {
	ID          int64     `db:"id"`
	ShortURL    string    `db:"short_url"`
	OriginalURL string    `db:"original_url"`
	CreatedAt   time.Time `db:"created_at"`
}
