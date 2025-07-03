package store

import (
	"context"
	"fmt"
)

type Debtors struct {
	ID               int64   `json:"id"`
	ReceivedAmount   int64   `json:"received_amount"`
	ReceivedCurrency string  `json:"received_currency"`
	DebtedAmount     int64   `json:"debted_amount"`
	DebtedCurrency   string  `json:"debted_currency"`
	UserID           int64   `json:"user_id"`
	CompanyID        int64   `json:"company_id"`
	Details          *string `json:"details"`
	Phone            *string `json:"phone"`
	IsBalanceEffect  int     `json:"is_balance_effect"`
	Type             int     `json:"type"`
	Status           int     `json:"status"`
	CreatedAt        string  `json:"created_at"`
}

/*
1. Ism Phone details

2. original summa, currency bilan

3. qabul qilinadigan summa, currency bilan

4.  type (qarz olish yoki berish)

5. balancega tasir qilsinmi?
*/

type DebtorsStorage struct {
	db DBTX
}

func NewDebtorsStorage(db DBTX) *DebtorsStorage {
	return &DebtorsStorage{db: db}
}

func (s *DebtorsStorage) Create(ctx context.Context, credits *Debtors) error {
	query := `
				INSERT INTO debtors (received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, status, company_id)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at
			`

	err := s.db.QueryRowContext(
		ctx,
		query,
		credits.ReceivedAmount,
		credits.ReceivedCurrency,
		credits.DebtedAmount,
		credits.DebtedCurrency,
		credits.UserID,
		credits.Details,
		credits.Phone,
		credits.IsBalanceEffect,
		credits.Type,
		credits.Status,
		credits.CompanyID,
	).Scan(
		&credits.ID,
		&credits.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *DebtorsStorage) GetByCompanyId(ctx context.Context, companyId int64) ([]Debtors, error) {
	query := `
				SELECT id, received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, status, created_at, company_id
				FROM debtors WHERE company_id = $1 and status != -1
			`

	var credits []Debtors
	rows, err := s.db.QueryContext(
		ctx,
		query,
		companyId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var credit Debtors
		err := rows.Scan(
			&credit.ID,
			&credit.ReceivedAmount,
			&credit.ReceivedCurrency,
			&credit.DebtedAmount,
			&credit.DebtedCurrency,
			&credit.UserID,
			&credit.Details,
			&credit.Phone,
			&credit.IsBalanceEffect,
			&credit.Type,
			&credit.Status,
			&credit.CreatedAt,
			&credit.CompanyID,
		)

		if err != nil {
			return nil, err
		}

		credits = append(credits, credit)

	}

	return credits, nil
}

func (s *DebtorsStorage) GetByUserId(ctx context.Context, userId int64) ([]Debtors, error) {
	query := `
				SELECT id, received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, status, created_at, company_id
				FROM debtors WHERE user_id = $1 and status != -1
			`

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
			&credit.ReceivedAmount,
			&credit.ReceivedCurrency,
			&credit.DebtedAmount,
			&credit.DebtedCurrency,
			&credit.UserID,
			&credit.Details,
			&credit.Phone,
			&credit.IsBalanceEffect,
			&credit.Type,
			&credit.Status,
			&credit.CreatedAt,
			&credit.CompanyID,
		)

		if err != nil {
			return nil, err
		}

		credits = append(credits, credit)

	}

	return credits, nil
}

func (s *DebtorsStorage) GetById(ctx context.Context, id int64) (*Debtors, error) {
	query := `
				SELECT id, received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, status, created_at, company_id
				FROM debtors WHERE id = $1 and status != -1
			`

	fmt.Printf("GetById ID %v", id)

	credit := &Debtors{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&credit.ID,
		&credit.ReceivedAmount,
		&credit.ReceivedCurrency,
		&credit.DebtedAmount,
		&credit.DebtedCurrency,
		&credit.UserID,
		&credit.Details,
		&credit.Phone,
		&credit.IsBalanceEffect,
		&credit.Type,
		&credit.Status,
		&credit.CreatedAt,
		&credit.CompanyID,
	)

	if err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *DebtorsStorage) Update(ctx context.Context, credit *Debtors) error {
	query := `
				UPDATE debtors SET received_amount = $1, received_currency = $2, debted_amount = $3, debted_currency = $4, 
				user_id = $5, details = $6, phone = $7, is_balance_effect = $8, type = $9, status = $10, company_id = $11 WHERE id = $12 and status != -1
			`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		credit.ReceivedAmount,
		credit.ReceivedCurrency,
		credit.DebtedAmount,
		credit.DebtedCurrency,
		credit.UserID,
		credit.Details,
		credit.Phone,
		credit.IsBalanceEffect,
		credit.Type,
		credit.Status,
		credit.CompanyID,
		credit.ID,
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
	query := `UPDATE debtors SET status = -1 WHERE id = $1`

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
