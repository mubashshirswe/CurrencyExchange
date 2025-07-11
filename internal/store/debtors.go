package store

import (
	"context"
	"fmt"
)

type Debtors struct {
	ID        int64   `json:"id"`
	Balance   int64   `json:"balance"`
	Currency  string  `json:"currency"`
	UserID    int64   `json:"user_id"`
	CompanyID int64   `json:"company_id"`
	Phone     *string `json:"phone"`
	CreatedAt string  `json:"created_at"`
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
				INSERT INTO debtors (balance, currency, user_id, phone, company_id)
				VALUES($1, $2, $3, $4, $5) RETURNING id, created_at
			`

	err := s.db.QueryRowContext(
		ctx,
		query,
		credits.Balance,
		credits.Currency,
		credits.UserID,
		credits.Phone,
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
				SELECT id, balance, currency, user_id, phone, company_id, created_at
				FROM debtors WHERE company_id = $1
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
			&credit.Balance,
			&credit.Currency,
			&credit.UserID,
			&credit.Phone,
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
				SELECT id, balance, currency, user_id, phone, company_id, created_at
				FROM debtors WHERE user_id = $1
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
			&credit.Balance,
			&credit.Currency,
			&credit.UserID,
			&credit.Phone,
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
				SELECT id, balance, currency, user_id, phone, company_id, created_at
				FROM debtors WHERE id = $1
			`

	fmt.Printf("GetById ID %v", id)

	credit := &Debtors{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&credit.ID,
		&credit.Balance,
		&credit.Currency,
		&credit.UserID,
		&credit.Phone,
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
				UPDATE debtors SET balance = $1, currency = $2, user_id = $3, phone = $4, 
				company_id = $5 WHERE id = $6
			`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		&credit.Balance,
		&credit.Currency,
		&credit.UserID,
		&credit.Phone,
		&credit.CompanyID,
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
