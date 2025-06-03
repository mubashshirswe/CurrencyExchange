package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Debtors struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	Amount          int64  `json:"amount"`
	SerialNo        string `json:"serial_no"`
	BalanceId       int64  `json:"balance_id"`
	CompanyId       int64  `json:"company_id"`
	Details         string `json:"details"`
	DebtorsName     string `json:"debtors_name"`
	DebtorsPhone    string `json:"debtors_phone"`
	CurrencyId      int64  `json:"currency_id"`
	CurrencyType    string `json:"currency_type"`
	Type            int    `json:"type"`
	Status          int    `json:"status"`
	IsBalanceEffect int    `json:"is_balance_effect"`
	CreatedAt       string `json:"created_at"`
}

type DebtorsStorage struct {
	db *sql.DB
}

func (s *DebtorsStorage) Create(ctx context.Context, credits *Debtors) error {
	query := `INSERT INTO debtors(user_id, amount, serial_no, balance_id, company_id, details, debtors_name, debtors_phone, currency_id, type, is_balance_effect, currency_type, status)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, created_at`
	err := s.db.QueryRowContext(
		ctx,
		query,
		credits.UserID,
		credits.Amount,
		credits.SerialNo,
		credits.BalanceId,
		credits.CompanyId,
		credits.Details,
		credits.DebtorsName,
		credits.DebtorsPhone,
		credits.CurrencyId,
		credits.Type,
		credits.IsBalanceEffect,
		credits.CurrencyType,
		credits.Status,
	).Scan(
		&credits.ID,
		&credits.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *DebtorsStorage) GetByUserId(ctx context.Context, userId int64) ([]Debtors, error) {
	query := `SELECT id, user_id, amount, serial_no, balance_id, company_id, 
						details, debtors_name, debtors_phone, currency_id, 
						type, created_at, is_balance_effect, currency_type, status FROM debtors WHERE user_id = $1 and status == 1`

	var credits []Debtors
	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var credit Debtors
		err := rows.Scan(
			&credit.ID,
			&credit.UserID,
			&credit.Amount,
			&credit.SerialNo,
			&credit.BalanceId,
			&credit.CompanyId,
			&credit.Details,
			&credit.DebtorsName,
			&credit.DebtorsPhone,
			&credit.CurrencyId,
			&credit.Type,
			&credit.CreatedAt,
			&credit.IsBalanceEffect,
			&credit.CurrencyType,
			&credit.Status,
		)

		if err != nil {
			return nil, err
		}

		credits = append(credits, credit)

	}

	return credits, nil
}

func (s *DebtorsStorage) GetById(ctx context.Context, id int64) (*Debtors, error) {
	query := `SELECT id, user_id, amount, serial_no, balance_id, company_id, 
						details, debtors_name, debtors_phone, currency_id, 
						type, created_at, is_balance_effect, currency_type, status FROM debtors WHERE id = $1 and status == 1`

	fmt.Printf("GetById ID %v", id)

	credit := &Debtors{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&credit.ID,
		&credit.UserID,
		&credit.Amount,
		&credit.SerialNo,
		&credit.BalanceId,
		&credit.CompanyId,
		&credit.Details,
		&credit.DebtorsName,
		&credit.DebtorsPhone,
		&credit.CurrencyId,
		&credit.Type,
		&credit.CreatedAt,
		&credit.IsBalanceEffect,
		&credit.CurrencyType,
		&credit.Status,
	)

	if err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *DebtorsStorage) Update(ctx context.Context, credits *Debtors) error {
	query := `UPDATE debtors SET amount = $1, balance_id = $2, company_id = $3, details = $4, debtors_name = $5, 
					debtors_phone = $6, currency_id = $7, type = $8, is_balance_effect = $9, currency_type = $10, status = $11 WHERE id = $12 and status == 1`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		credits.Amount,
		credits.BalanceId,
		credits.CompanyId,
		credits.Details,
		credits.DebtorsName,
		credits.DebtorsPhone,
		credits.CurrencyId,
		credits.Type,
		credits.IsBalanceEffect,
		credits.CurrencyType,
		credits.Status,
		credits.ID,
	)

	if err != nil {
		return err
	}

	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return fmt.Errorf("DEBTORS NOT FOUND")
	}

	return nil
}

func (s *DebtorsStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM debtors  WHERE id = $1`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return fmt.Errorf("DEBTORS NOT FOUND")
	}

	return nil
}
