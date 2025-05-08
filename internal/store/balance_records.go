package store

import (
	"context"
	"database/sql"
)

type BalanceRecord struct {
	ID         int64  `json:"id"`
	Amount     int64  `json:"amount"`
	UserID     int64  `json:"user_id"`
	BalanceID  int64  `json:"balance_id"`
	CompanyID  int64  `json:"company_id"`
	Details    string `json:"details"`
	CurrenctID int64  `json:"currency_id"`
	Type       int64  `json:"type"`
	CreatedAt  string `json:"created_at"`
}

type BalanceRecordStorage struct {
	db *sql.DB
}

func (s *BalanceRecordStorage) Create(ctx context.Context, balanceRecord *BalanceRecord) error {
	query := `INSERT INTO balance_records(amount, user_id, balance_id, company_id, details, currency_id, type)
				VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		balanceRecord.Amount,
		balanceRecord.UserID,
		balanceRecord.BalanceID,
		balanceRecord.CompanyID,
		balanceRecord.Details,
		balanceRecord.CurrenctID,
		balanceRecord.Type,
	).Scan(
		&balanceRecord.ID,
		&balanceRecord.CreatedAt,
	)

	return err
}

func (s *BalanceRecordStorage) Update(ctx context.Context, balanceRecord *BalanceRecord) error {
	query := `UPDATE balance_records SET amount = $1, user_id = $2, balance_id = $3,
				company_id =$4, details = $5, currency_id = $6, type = $6 WHERE id = $7`

	_, err := s.db.ExecContext(
		ctx,
		query,
		balanceRecord.Amount,
		balanceRecord.UserID,
		balanceRecord.BalanceID,
		balanceRecord.CompanyID,
		balanceRecord.Details,
		balanceRecord.CurrenctID,
		balanceRecord.Type,
		balanceRecord.ID,
	)

	return err
}

func (s *BalanceRecordStorage) GetByBalanceId(ctx context.Context, balance_id int64) ([]BalanceRecord, error) {
	query := `SELECT id, amount, user_id, balance_id, company_id, details, currency_id, type, created_at balance_records WHERE balance_id = $1`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		balance_id,
	)
	if err != nil {
		return nil, err
	}

	var balanceRecords []BalanceRecord

	for rows.Next() {
		balance := BalanceRecord{}

		err := rows.Scan(
			&balance.ID,
			&balance.Amount,
			&balance.UserID,
			&balance.BalanceID,
			&balance.CompanyID,
			&balance.Details,
			&balance.CurrenctID,
			&balance.Type,
			&balance.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		balanceRecords = append(balanceRecords, balance)
	}

	return balanceRecords, nil
}

func (s *BalanceRecordStorage) GetByUserId(ctx context.Context, user_id int64) ([]BalanceRecord, error) {
	query := `SELECT id, amount, user_id, balance_id, company_id, details, currency_id, type, created_at balance_records WHERE user_id = $1`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		user_id,
	)
	if err != nil {
		return nil, err
	}

	var balanceRecords []BalanceRecord

	for rows.Next() {
		balance := BalanceRecord{}

		err := rows.Scan(
			&balance.ID,
			&balance.Amount,
			&balance.UserID,
			&balance.BalanceID,
			&balance.CompanyID,
			&balance.Details,
			&balance.CurrenctID,
			&balance.Type,
			&balance.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		balanceRecords = append(balanceRecords, balance)
	}

	return balanceRecords, nil
}

func (s *BalanceRecordStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM balance_records WHERE id = $1`

	_, err := s.db.QueryContext(
		ctx,
		query,
		id,
	)

	return err
}
