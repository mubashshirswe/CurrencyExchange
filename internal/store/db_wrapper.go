package store

import (
	"context"
	"database/sql"
)

type DBWrapper struct {
	db *sql.DB
}

func (d *DBWrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *DBWrapper) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DBWrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *DBWrapper) Commit() error {
	return nil // no-op for DB
}

func (d *DBWrapper) Rollback() error {
	return nil // no-op for DB
}
