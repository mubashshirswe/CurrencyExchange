package store

import (
	"context"
	"database/sql"
	"errors"
)

type Balance struct {
	ID         int64  `json:"id"`
	Balance    int64  `json:"balance"`
	UserId     int64  `json:"user_id"`
	InOutLay   int64  `json:"in_out_lay"`
	OutInLay   int64  `json:"out_in_lay"`
	CompanyId  int64  `json:"company_id"`
	CurrencyId int64  `json:"currency_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type BalanceStorage struct {
	db *sql.DB
}

func (s *BalanceStorage) Create(ctx context.Context, balance *Balance) error {
	query := `INSERT INTO balances(balance, user_id, in_out_lay, out_in_lay, company_id, currency_id)
				VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		balance.Balance,
		balance.UserId,
		balance.InOutLay,
		balance.OutInLay,
		balance.CompanyId,
		balance.CurrencyId).Scan(
		&balance.ID,
		&balance.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *BalanceStorage) GetAll(ctx context.Context) ([]Balance, error) {
	query := `SELECT id, balance, user_id, in_out_lay, out_in_lay, company_id, created_at, updated_at, currency_id FROM balances`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []Balance
	for rows.Next() {
		balance := &Balance{}
		err := rows.Scan(
			&balance.ID,
			&balance.Balance,
			&balance.UserId,
			&balance.InOutLay,
			&balance.OutInLay,
			&balance.CompanyId,
			&balance.CreatedAt,
			&balance.UpdatedAt,
			&balance.CurrencyId,
		)
		if err != nil {
			return nil, err
		}

		balances = append(balances, *balance)
	}

	return balances, nil
}

func (s *BalanceStorage) GetByUserId(ctx context.Context, userId *int64) ([]Balance, error) {
	query := `SELECT id, balance, user_id, in_out_lay, out_in_lay, company_id, created_at, updated_at, currency_id FROM balances WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []Balance
	for rows.Next() {
		balance := &Balance{}
		err := rows.Scan(
			&balance.ID,
			&balance.Balance,
			&balance.UserId,
			&balance.InOutLay,
			&balance.OutInLay,
			&balance.CompanyId,
			&balance.CreatedAt,
			&balance.UpdatedAt,
			&balance.CurrencyId,
		)
		if err != nil {
			return nil, err
		}

		balances = append(balances, *balance)
	}

	return balances, nil
}

func (s *BalanceStorage) GetById(ctx context.Context, id *int64) (*Balance, error) {
	query := `SELECT id, balance, user_id, in_out_lay, out_in_lay, company_id, created_at, updated_at, currency_id FROM balances WHERE id = $1`
	balance := &Balance{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&balance.ID,
		&balance.Balance,
		&balance.UserId,
		&balance.InOutLay,
		&balance.OutInLay,
		&balance.CompanyId,
		&balance.CreatedAt,
		&balance.UpdatedAt,
		&balance.CurrencyId,
	)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (s *BalanceStorage) Update(ctx context.Context, balance *Balance) error {
	query := `UPDATE balances SET balance = $1, in_out_lay = $2, out_in_lay = $3
		WHERE id = $4`

	rows, err := s.db.ExecContext(ctx, query, balance.Balance, balance.InOutLay, balance.OutInLay, balance.ID)
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
