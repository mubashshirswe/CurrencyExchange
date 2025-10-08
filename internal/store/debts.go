package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type Debts struct {
	ID                 int64                   `json:"id"`
	FullName           string                  `json:"full_name"`
	ReceivedIncomes    []types.ReceivedIncomes `json:"received_incomes"`
	DebtedAmount       int64                   `json:"debted_amount"`
	DebtedCurrency     string                  `json:"debted_currency"`
	State              int64                   `json:"state"`
	UserID             int64                   `json:"user_id"`
	CompanyID          int64                   `json:"company_id"`
	DebtorID           int64                   `json:"debtor_id"` // Fixed: consistent naming
	Details            *string                 `json:"details"`
	Phone              *string                 `json:"phone"`
	IsBalanceEffect    int                     `json:"is_balance_effect"`
	Type               int                     `json:"type"`
	CreatedAt          time.Time               `json:"-"`
	CreatedAtFormatted string                  `json:"created_at"`
}

type DebtsStorage struct {
	db DBTX
}

func NewDebtsStorage(db DBTX) *DebtsStorage {
	return &DebtsStorage{db: db}
}

func (s *DebtsStorage) Create(ctx context.Context, debts *Debts) error {
	query := `
		INSERT INTO debts (received_incomes, debted_amount, debted_currency, user_id, 
		details, phone, is_balance_effect, type, company_id, debtor_id, state)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at
	`

	// Convert ReceivedIncomes to JSON for storage
	incomesJSON, err := json.Marshal(debts.ReceivedIncomes)
	if err != nil {
		return fmt.Errorf("failed to marshal received_incomes: %w", err)
	}

	err = s.db.QueryRowContext(
		ctx,
		query,
		incomesJSON,
		debts.DebtedAmount,
		debts.DebtedCurrency,
		debts.UserID,
		debts.Details,
		debts.Phone,
		debts.IsBalanceEffect,
		debts.Type,
		debts.CompanyID,
		debts.DebtorID,
		debts.State,
	).Scan(
		&debts.ID,
		&debts.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create debt: %w", err)
	}

	return nil
}

func (s *DebtsStorage) GetByCompanyID(ctx context.Context, companyID int64) ([]Debts, error) {
	query := `
		SELECT id, received_incomes, debted_amount, debted_currency, user_id, 
		details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
		FROM debts WHERE company_id = $1 ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query debts: %w", err)
	}
	defer rows.Close()

	return s.scanDebts(rows)
}

func (s *DebtsStorage) GetByDebtorID(ctx context.Context, debtorID int64, pagination types.Pagination) ([]Debts, error) {
	query := `
		SELECT id, received_incomes, debted_amount, debted_currency, user_id, 
		details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
		FROM debts WHERE debtor_id = $1 ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, debtorID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query debts: %w", err)
	}
	defer rows.Close()

	return s.scanDebts(rows)
}

func (s *DebtsStorage) GetByUserID(ctx context.Context, userID int64, pagination types.Pagination) ([]Debts, error) {
	query := `
		SELECT id, received_incomes, debted_amount, debted_currency, user_id, 
		details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
		FROM debts WHERE user_id = $1 ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query debts: %w", err)
	}
	defer rows.Close()

	return s.scanDebts(rows)
}

func (s *DebtsStorage) GetByID(ctx context.Context, id int64) (*Debts, error) {
	query := `
		SELECT id, received_incomes, debted_amount, debted_currency, user_id, 
		details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
		FROM debts WHERE id = $1
	`

	debt := &Debts{}
	var incomesJSON []byte

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&debt.ID,
		&incomesJSON,
		&debt.DebtedAmount,
		&debt.DebtedCurrency,
		&debt.UserID,
		&debt.Details,
		&debt.Phone,
		&debt.IsBalanceEffect,
		&debt.Type,
		&debt.CreatedAt,
		&debt.CompanyID,
		&debt.DebtorID,
		&debt.State,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("debt not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get debt: %w", err)
	}

	// Unmarshal received incomes
	if len(incomesJSON) > 0 {
		if err := json.Unmarshal(incomesJSON, &debt.ReceivedIncomes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal received_incomes: %w", err)
		}
	}

	debt.CreatedAtFormatted = s.formatTime(debt.CreatedAt)

	return debt, nil
}

func (s *DebtsStorage) Update(ctx context.Context, debt *Debts) error {
	query := `
		UPDATE debts SET received_incomes = $1, debted_amount = $2, debted_currency = $3, 
		user_id = $4, details = $5, phone = $6, is_balance_effect = $7, type = $8, 
		company_id = $9, debtor_id = $10, state = $11  WHERE id = $12
	`

	// Convert ReceivedIncomes to JSON
	incomesJSON, err := json.Marshal(debt.ReceivedIncomes)
	if err != nil {
		return fmt.Errorf("failed to marshal received_incomes: %w", err)
	}

	result, err := s.db.ExecContext(
		ctx,
		query,
		incomesJSON,
		debt.DebtedAmount,
		debt.DebtedCurrency,
		debt.UserID,
		debt.Details,
		debt.Phone,
		debt.IsBalanceEffect,
		debt.Type,
		debt.CompanyID,
		debt.DebtorID,
		debt.State,
		debt.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update debt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("debt not found")
	}

	return nil
}

func (s *DebtsStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM debts WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete debt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("debt not found")
	}

	return nil
}

// Helper function to scan multiple debts from rows
func (s *DebtsStorage) scanDebts(rows *sql.Rows) ([]Debts, error) {
	var debts []Debts

	for rows.Next() {
		var debt Debts
		var incomesJSON []byte

		err := rows.Scan(
			&debt.ID,
			&incomesJSON,
			&debt.DebtedAmount,
			&debt.DebtedCurrency,
			&debt.UserID,
			&debt.Details,
			&debt.Phone,
			&debt.IsBalanceEffect,
			&debt.Type,
			&debt.CreatedAt,
			&debt.CompanyID,
			&debt.DebtorID,
			&debt.State,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan debt: %w", err)
		}

		// Unmarshal received incomes
		if len(incomesJSON) > 0 {
			if err := json.Unmarshal(incomesJSON, &debt.ReceivedIncomes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal received_incomes: %w", err)
			}
		}

		debt.CreatedAtFormatted = s.formatTime(debt.CreatedAt)
		debts = append(debts, debt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return debts, nil
}

// Helper function to format time consistently
func (s *DebtsStorage) formatTime(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Tashkent")
	if err != nil {
		// Fallback to UTC if location load fails
		return t.Format("2006-01-02 15:04:05")
	}
	return t.In(loc).Format("2006-01-02 15:04:05")
}
