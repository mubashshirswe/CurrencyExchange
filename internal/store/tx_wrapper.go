package store

import (
	"context"
	"database/sql"
)

type TxWrapper struct {
	tx *sql.Tx
}

func (t *TxWrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *TxWrapper) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *TxWrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *TxWrapper) Commit() error {
	return t.tx.Commit()
}

func (t *TxWrapper) Rollback() error {
	return t.tx.Rollback()
}
