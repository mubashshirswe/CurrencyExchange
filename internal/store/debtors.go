package store

import (
	"context"
	"fmt"
	"time"

	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type Debtors struct {
	ID                 int64     `json:"id"`
	Balance            int64     `json:"balance"`
	Currency           string    `json:"currency"`
	UserID             int64     `json:"user_id"`
	CompanyID          int64     `json:"company_id"`
	Phone              string    `json:"phone"`
	FullName           string    `json:"full_name"`
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

type DebtorsStorage struct {
	db DBTX
}

func NewDebtorsStorage(db DBTX) *DebtorsStorage {
	return &DebtorsStorage{db: db}
}

func (s *DebtorsStorage) Create(ctx context.Context, credits *Debtors) error {
	query := `
				INSERT INTO debtors (balance, currency, user_id, phone, company_id, full_name, created_at)
				VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at
			`
	loc, _ := time.LoadLocation("Asia/Tashkent")
	nowUz := time.Now().In(loc)

	err := s.db.QueryRowContext(
		ctx,
		query,
		credits.Balance,
		credits.Currency,
		credits.UserID,
		credits.Phone,
		credits.CompanyID,
		credits.FullName,
		nowUz,
	).Scan(
		&credits.ID,
		&credits.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *DebtorsStorage) GetByCompanyId(
	ctx context.Context,
	companyId int64,
	search *string,
	dateFilter *string, // "2025-02-18" or nil/empty = all
	pagination types.Pagination,
) ([]Debtors, error) {

	query := `
        SELECT DISTINCT d.id, d.balance, d.currency, d.user_id, d.phone, d.company_id, d.created_at, d.full_name
        FROM debtors d
        JOIN debts db ON d.id = db.debtor_id
        WHERE d.company_id = $1
    `

	args := []interface{}{companyId}
	argIndex := 2

	// Search filter
	if search != nil && *search != "" {
		searchLike := "%" + *search + "%"
		query += fmt.Sprintf(`
            AND (
                CAST(d.balance AS TEXT) ILIKE $%d OR
                d.currency ILIKE $%d OR
                d.phone ILIKE $%d OR
                d.full_name ILIKE $%d
            )
        `, argIndex, argIndex, argIndex, argIndex)
		args = append(args, searchLike)
		argIndex++
	}

	// Date filter on debt creation date
	if dateFilter != nil && *dateFilter != "" {
		query += fmt.Sprintf(" AND db.created_at::date = $%d", argIndex)
		args = append(args, *dateFilter)
		argIndex++
	}

	// Order & pagination
	orderBy := pagination.OrderBy
	if orderBy == "" {
		orderBy = "d.created_at DESC" // or "MAX(db.created_at) DESC" if you want latest debt
	}

	query += fmt.Sprintf(`
        ORDER BY %s
        OFFSET $%d LIMIT $%d
    `, orderBy, argIndex, argIndex+1)

	args = append(args, pagination.Offset, pagination.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var debtors []Debtors

	for rows.Next() {
		var d Debtors
		err := rows.Scan(
			&d.ID,
			&d.Balance,
			&d.Currency,
			&d.UserID,
			&d.Phone,
			&d.CompanyID,
			&d.CreatedAt,
			&d.FullName,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		loc, _ := time.LoadLocation("Asia/Tashkent")
		d.CreatedAtFormatted = d.CreatedAt.In(loc).Format("2006-01-02 15:04:05")

		debtors = append(debtors, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return debtors, nil
}

func (s *DebtorsStorage) GetByBalanceInfo(ctx context.Context, companyId int64) ([]map[string]interface{}, error) {
	query :=
		`
		SELECT
   			 currency,
   			 SUM(CASE WHEN balance > 0 THEN balance ELSE 0 END) AS positive_balance,
   			 SUM(CASE WHEN balance < 0 THEN balance ELSE 0 END) AS negative_balance
		FROM debtors WHERE company_id = $1 GROUP BY currency
	`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		companyId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var currency string
		var positiveBalance int64
		var negativeBalance int64

		err := rows.Scan(&currency, &positiveBalance, &negativeBalance)

		if err != nil {
			return nil, err
		}

		info := map[string]interface{}{
			"currency":         currency,
			"positive_balance": positiveBalance,
			"negative_balance": negativeBalance,
		}

		result = append(result, info)
	}
	return result, nil
}

func (s *DebtorsStorage) GetByUserId(ctx context.Context, userId int64, pagination types.Pagination) ([]Debtors, error) {
	query := `
				SELECT id, balance, currency, user_id, phone, company_id, created_at, full_name
				FROM debtors WHERE user_id = $1 ORDER BY balance DESC
	` + fmt.Sprintf(" OFFSET %v LIMIT %v", pagination.Offset, pagination.Limit)

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
			&credit.CompanyID,
			&credit.CreatedAt,
			&credit.FullName,
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

func (s *DebtorsStorage) GetById(ctx context.Context, id int64) (*Debtors, error) {
	query := `
				SELECT id, balance, currency, user_id, phone, company_id, created_at, full_name
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
		&credit.CompanyID,
		&credit.CreatedAt,
		&credit.FullName,
	)

	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Tashkent")
	createdAtInTashkent := credit.CreatedAt.In(loc)
	credit.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

	return credit, nil
}

func (s *DebtorsStorage) Update(ctx context.Context, credit *Debtors) error {
	query := `
				UPDATE debtors SET balance = $1, currency = $2, user_id = $3, phone = $4, full_name = $5,
				company_id = $6 WHERE id = $7
			`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		&credit.Balance,
		&credit.Currency,
		&credit.UserID,
		&credit.Phone,
		&credit.FullName,
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
	query := `DELETE FROM debtors WHERE id = $1 and balance = 0`

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
