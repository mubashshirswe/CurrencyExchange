package store

import (
	"context"
	"database/sql"
	"errors"
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
	query := `INSERT INTO currencies(name, sell, buy) VALUES($1, $2, $3) RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query, currnecy.Name, currnecy.Sell, currnecy.Buy).Scan(&currnecy.ID, &currnecy.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *CurrencyStorage) GetAll(ctx context.Context) ([]Currency, error) {
	query := `SELECT * FROM currencies`
	var currencies []Currency
	rows, err := s.db.QueryContext(
		ctx,
		query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		currnecy := &Currency{}
		err := rows.Scan(
			&currnecy.ID,
			&currnecy.Name,
			&currnecy.Sell,
			&currnecy.Buy,
			&currnecy.CreatedAt)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, *currnecy)
	}

	return currencies, nil
}

func (s *CurrencyStorage) Update(ctx context.Context, currecy *Currency) error {
	query := `UPDATE currencies SET name = $1, sell = $2, buy = $3 RETURNING id`
	rows, err := s.db.ExecContext(
		ctx,
		query,
		currecy.Name,
		currecy.Sell,
		currecy.Buy,
	)

	if err != nil {
		return err
	}

	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return errors.New("NOT FOUND")
	}

	return nil
}

func (s *CurrencyStorage) Delete(ctx context.Context) ([]Currency, error) {

	return nil, nil
}
