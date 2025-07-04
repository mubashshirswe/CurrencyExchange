package store

import (
	"context"
	"database/sql"
	"time"
)

type BalanceRecord struct {
	ID            int64     `json:"id"`
	Amount        int64     `json:"amount"`
	UserID        int64     `json:"user_id"`
	BalanceID     int64     `json:"balance_id"`
	CompanyID     int64     `json:"company_id"`
	TransactionId *int64    `json:"transaction_id"`
	DebtorId      *int64    `json:"debtor_id"`
	ExchangeId    *int64    `json:"exchange_id"`
	Details       *string   `json:"details"`
	Currency      string    `json:"currency"`
	Type          int64     `json:"type"`
	CreatedAt     time.Time `json:"created_at"`
}

type BalanceRecordStorage struct {
	db DBTX
}

func NewBalanceRecordStorage(db DBTX) *BalanceRecordStorage {
	return &BalanceRecordStorage{db: db}
}

func (s *BalanceRecordStorage) Create(ctx context.Context, balanceRecord *BalanceRecord) error {
	query := `
				INSERT INTO balance_records(amount, user_id, balance_id, company_id, transaction_id, debtor_id, exchange_id, details, currency, type)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		balanceRecord.Amount,
		balanceRecord.UserID,
		balanceRecord.BalanceID,
		balanceRecord.CompanyID,
		balanceRecord.TransactionId,
		balanceRecord.DebtorId,
		balanceRecord.ExchangeId,
		balanceRecord.Details,
		balanceRecord.Currency,
		balanceRecord.Type,
	).Scan(
		&balanceRecord.ID,
		&balanceRecord.CreatedAt,
	)

	return err
}

func (s *BalanceRecordStorage) Update(ctx context.Context, balanceRecord *BalanceRecord) error {
	query := `
				UPDATE balance_records SET amount = $1, user_id = $2, balance_id = $3, company_id = $4, transaction_id = $5, debtor_id = $6,  
				details = $7, currency = $8, type = $9, exchange_id = $10 WHERE id = $11
			`

	_, err := s.db.ExecContext(
		ctx,
		query,
		balanceRecord.Amount,
		balanceRecord.UserID,
		balanceRecord.BalanceID,
		balanceRecord.CompanyID,
		balanceRecord.TransactionId,
		balanceRecord.DebtorId,
		balanceRecord.Details,
		balanceRecord.Currency,
		balanceRecord.Type,
		balanceRecord.ExchangeId,
		balanceRecord.ID,
	)

	return err
}

func (s *BalanceRecordStorage) GetByFieldAndDate(ctx context.Context, fieldName, from, to string, fieldValue any) ([]BalanceRecord, error) {
	query := `
				SELECT id, amount, user_id, balance_id, company_id, transaction_id, debtor_id, exchange_id, details, currency, type, created_at
				FROM balance_records WHERE ` + fieldName + ` = $1 AND created_at BETWEEN $2 AND $3 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		fieldValue,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.FetchDataFromQuery(rows)
}

func (s *BalanceRecordStorage) GetByField(ctx context.Context, fieldName string, fieldValue any) ([]BalanceRecord, error) {
	query := `
				SELECT id, amount, user_id, balance_id, company_id, transaction_id, debtor_id, exchange_id, details, currency, type, created_at
				FROM balance_records WHERE ` + fieldName + ` = $1 ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		fieldValue,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.FetchDataFromQuery(rows)
}

func (s *BalanceRecordStorage) FetchDataFromQuery(rows *sql.Rows) ([]BalanceRecord, error) {
	var balanceRecords []BalanceRecord
	for rows.Next() {
		balance := BalanceRecord{}

		err := rows.Scan(
			&balance.ID,
			&balance.Amount,
			&balance.UserID,
			&balance.BalanceID,
			&balance.CompanyID,
			&balance.TransactionId,
			&balance.DebtorId,
			&balance.ExchangeId,
			&balance.Details,
			&balance.Currency,
			&balance.Type,
			&balance.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("Asia/Tashkent")
		createdAtInTashkent := balance.CreatedAt.In(loc)
		balance.CreatedAt = createdAtInTashkent

		balanceRecords = append(balanceRecords, balance)
	}

	return balanceRecords, nil
}

func (s *BalanceRecordStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM balance_records WHERE id = $1`

	_, err := s.db.ExecContext(
		ctx,
		query,
		id,
	)

	return err
}

func (s *BalanceRecordStorage) DeleteByExchangeId(ctx context.Context, id int64) error {
	query := `DELETE FROM balance_records WHERE exchange_id = $1`

	_, err := s.db.ExecContext(
		ctx,
		query,
		id,
	)

	return err
}
