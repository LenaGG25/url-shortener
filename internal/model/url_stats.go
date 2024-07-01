package model

type URLStats struct {
	ID            int64  `db:"id"`
	ShortURL      string `db:"short_url"`
	RequestNumber int64  `db:"request_number"`
}
