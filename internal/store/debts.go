package store

import (
	"context"
	"fmt"
	"time"

	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type Debts struct {
	ID                 int64     `json:"id"`
	FullName           string    `json:"full_name"`
	ReceivedAmount     int64     `json:"received_amount"`
	ReceivedCurrency   string    `json:"received_currency"`
	DebtedAmount       int64     `json:"debted_amount"`
	DebtedCurrency     string    `json:"debted_currency"`
	State              int64     `json:"state"`
	UserID             int64     `json:"user_id"`
	CompanyID          int64     `json:"company_id"`
	DebtorId           int64     `json:"debtor_id"`
	Details            *string   `json:"details"`
	Phone              *string   `json:"phone"`
	IsBalanceEffect    int       `json:"is_balance_effect"`
	Type               int       `json:"type"`
	CreatedAt          time.Time `json:"-"`
	CreatedAtFormatted string    `json:"created_at"`
}

/*
1. Ism Phone details

2. original summa, currency bilan

3. qabul qilinadigan summa, currency bilan

4.  type (qarz olish yoki berish)

5. balancega tasir qilsinmi?
*/

type DebtsStorage struct {
	db DBTX
}

func NewDebtsStorage(db DBTX) *DebtsStorage {
	return &DebtsStorage{db: db}
}

func (s *DebtsStorage) Create(ctx context.Context, credits *Debts) error {
	query := `
				INSERT INTO debts (received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, company_id, debtor_id, state)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at
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
		credits.CompanyID,
		credits.DebtorId,
		credits.State,
	).Scan(
		&credits.ID,
		&credits.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *DebtsStorage) GetByCompanyId(ctx context.Context, companyId int64) ([]Debts, error) {
	query := `
				SELECT id, received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
				FROM debts WHERE company_id = $1 
			`

	var credits []Debts
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
		var credit Debts
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
			&credit.CreatedAt,
			&credit.CompanyID,
			&credit.DebtorId,
			&credit.State,
		)

		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("Asia/Tashkent")
		createdAtInTashkent := credit.CreatedAt.In(loc)
		credit.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

		credits = append(credits, credit)

	}

	return credits, nil
}

func (s *DebtsStorage) GetByDebtorId(ctx context.Context, debtorId int64, pagination types.Pagination) ([]Debts, error) {
	query := `
				SELECT id, received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
				FROM debts WHERE debtor_id = $1 	ORDER BY created_at DESC
	` + fmt.Sprintf(" OFFSET %v LIMIT %v", pagination.Offset, pagination.Limit)

	var credits []Debts
	rows, err := s.db.QueryContext(
		ctx,
		query,
		debtorId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var credit Debts
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
			&credit.CreatedAt,
			&credit.CompanyID,
			&credit.DebtorId,
			&credit.State,
		)

		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("Asia/Tashkent")
		createdAtInTashkent := credit.CreatedAt.In(loc)
		credit.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

		credits = append(credits, credit)

	}

	return credits, nil
}

func (s *DebtsStorage) GetByUserId(ctx context.Context, userId int64, pagination types.Pagination) ([]Debts, error) {
	query := `
				SELECT id, received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
				FROM debts WHERE user_id = $1 	ORDER BY created_at DESC
	` + fmt.Sprintf(" OFFSET %v LIMIT %v", pagination.Offset, pagination.Limit)

	var credits []Debts
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
		var credit Debts
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
			&credit.CreatedAt,
			&credit.CompanyID,
			&credit.DebtorId,
			&credit.State,
		)

		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("Asia/Tashkent")
		createdAtInTashkent := credit.CreatedAt.In(loc)
		credit.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

		credits = append(credits, credit)

	}

	return credits, nil
}

func (s *DebtsStorage) GetById(ctx context.Context, id int64) (*Debts, error) {
	query := `
				SELECT id, received_amount, received_currency, debted_amount, debted_currency, user_id, 
				details, phone, is_balance_effect, type, created_at, company_id, debtor_id, state
				FROM debts WHERE id = $1
			`

	fmt.Printf("GetById ID %v", id)

	credit := &Debts{}
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
		&credit.CreatedAt,
		&credit.CompanyID,
		&credit.DebtorId,
		&credit.State,
	)

	loc, _ := time.LoadLocation("Asia/Tashkent")
	createdAtInTashkent := credit.CreatedAt.In(loc)
	credit.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

	if err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *DebtsStorage) Update(ctx context.Context, credit *Debts) error {
	query := `
				UPDATE debts SET received_amount = $1, received_currency = $2, debted_amount = $3, debted_currency = $4, 
				user_id = $5, details = $6, phone = $7, is_balance_effect = $8, type = $9, company_id = $10, debtor_id = $11, state = $12  WHERE id = $13
			`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		&credit.ReceivedAmount,
		&credit.ReceivedCurrency,
		&credit.DebtedAmount,
		&credit.DebtedCurrency,
		&credit.UserID,
		&credit.Details,
		&credit.Phone,
		&credit.IsBalanceEffect,
		&credit.Type,
		&credit.CompanyID,
		&credit.DebtorId,
		&credit.State,
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

func (s *DebtsStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM debts WHERE id = $1`

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
		return fmt.Errorf("DEBT NOT FOUND")
	}

	return nil
}
