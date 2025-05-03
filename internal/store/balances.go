package store

import (
	"context"
	"database/sql"
)

type Balance struct {
	ID        int64  `json:"id"`
	Balance   int64  `json:"balance"`
	UserId    int64  `json:"user_id"`
	InOutLay  int64  `json:"in_out_lay"`
	OutInLay  int64  `json:"out_in_lay"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BalanceStorage struct {
	db *sql.DB
}

func (s *BalanceStorage) Create(ctx context.Context, balance *Balance) error {
	query := `INSERT INTO balances(balance, user_id, in_out_lay, out_in_lay)
				VALUES($1, $2, $3, $4) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		balance.Balance,
		balance.UserId,
		balance.InOutLay,
		balance.OutInLay).Scan(
		&balance.ID,
		&balance.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *BalanceStorage) GetAll(ctx context.Context) ([]Balance, error) {

	return nil, nil
}

func (s *BalanceStorage) GetByUserId(ctx context.Context, userId *int64) ([]Balance, error) {

	return nil, nil
}

func (s *BalanceStorage) GetById(ctx context.Context, id *int64) (*Balance, error) {

	return nil, nil
}

func (s *BalanceStorage) Update(ctx context.Context, balance *Balance) error {

	return nil
}
