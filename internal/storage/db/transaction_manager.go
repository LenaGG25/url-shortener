package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

const key = "transaction"

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

func (t *TransactionManager) RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	if err = f(context.WithValue(ctx, key, tx)); err != nil { // unknown error
		errRollback := tx.Rollback(ctx)           // one more error
		return multierr.Combine(err, errRollback) // unknown error + one more error
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (t *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok {
		return tx
	}

	return t.pool
}
