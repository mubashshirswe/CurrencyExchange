package store

import (
	"context"
	"database/sql"
)

type Currency struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Sell      *int64 `json:"sell"`
	Buy       *int64 `json:"buy"`
	CreatedAt string `json:"created_at"`
}

type CurrencyStorage struct {
	db *sql.DB
}

func (s *CurrencyStorage) Create(ctx context.Context, currnecy *Currency) error {

	return nil
}

func (s *CurrencyStorage) GetAll(ctx context.Context) ([]Currency, error) {

	return nil, nil
}
