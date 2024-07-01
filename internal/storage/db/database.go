package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Database struct{}

func newDatabase() *Database {
	return &Database{}
}

// Get возвращает в переменную результат запроса
func (db *Database) Get(ctx context.Context, cluster QueryEngine, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, cluster, dest, query, args...)
}

// Select получить несколько строк данных и записать их в слайс
func (db *Database) Select(ctx context.Context, cluster QueryEngine, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, cluster, dest, query, args...)
}

// Exec выполнить sql запрос, когда не важен вывод, просто хотим изменить данные
func (db *Database) Exec(ctx context.Context, cluster QueryEngine, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return cluster.Exec(ctx, query, args...)
}

// ExecQueryRow выполнить sql запрос и вернуть строчку, затронутую в результате выполнения запроса
func (db *Database) ExecQueryRow(ctx context.Context, cluster QueryEngine, query string, args ...interface{}) pgx.Row {
	return cluster.QueryRow(ctx, query, args...)
}
