package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewDB constructor for Database
func NewDB(ctx context.Context, dsn string) (*Database, *TransactionManager, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}

	transactionManager := NewTransactionManager(pool)

	return newDatabase(), transactionManager, nil
}
