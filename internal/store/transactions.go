package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type Transaction struct {
	ID                 int64                     `json:"id"`
	ReceivedCompanyId  int64                     `json:"received_company_id"`
	ReceivedUserId     int64                     `json:"received_user_id"`
	ReceivedIncomes    []types.ReceivedIncomes   `json:"received_incomes"`
	DeliveredOutcomes  []types.DeliveredOutcomes `json:"delivered_outcomes"`
	DeliveredCompanyId int64                     `json:"delivered_company_id"`
	DeliveredUserId    *int64                    `json:"delivered_user_id"`
	ServiceFee         string                    `json:"service_fee"`
	Phone              string                    `json:"phone"`
	Details            string                    `json:"details"`
	Status             int64                     `json:"status"`
	Type               int64                     `json:"type"`
	CreatedAt          time.Time                 `json:"-"`
	CreatedAtFormatted string                    `json:"created_at"`
}

type TransactionStorage struct {
	db DBTX
}

func NewTransactionStorage(db DBTX) *TransactionStorage {
	return &TransactionStorage{db: db}
}

func (s *TransactionStorage) Archive(ctx context.Context, companyId int64) error {
	query := `UPDATE transactions SET status = $1 WHERE status = $2 and company_id = $3`
	rows, err := s.db.ExecContext(ctx, query, STATUS_ARCHIVED, STATUS_COMPLETED, companyId)
	if err != nil {
		return err
	}

	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return fmt.Errorf("TRANSACTION NOT FOUND")
	}

	return nil
}

func (s *TransactionStorage) Create(ctx context.Context, tr *Transaction) error {
	receivedIncomesJSON, err := json.Marshal(tr.ReceivedIncomes)
	if err != nil {
		return err
	}

	deliveredOutcomesJSON, err := json.Marshal(tr.DeliveredOutcomes)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Tashkent")
	nowUz := time.Now().In(loc)

	query := `
			INSERT INTO transactions(
				service_fee, received_incomes, delivered_outcomes,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at`

	err = s.db.QueryRowContext(
		ctx,
		query,
		tr.ServiceFee,
		receivedIncomesJSON,
		deliveredOutcomesJSON,
		tr.ReceivedCompanyId,
		tr.DeliveredCompanyId,
		tr.ReceivedUserId,
		tr.DeliveredUserId,
		tr.Phone,
		tr.Details,
		STATUS_CREATED,
		tr.Type,
		nowUz,
	).Scan(
		&tr.ID,
		&tr.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
func (s *TransactionStorage) Update(ctx context.Context, tr *Transaction) error {
	receivedIncomesJSON, err := json.Marshal(tr.ReceivedIncomes)
	if err != nil {
		return err
	}

	deliveredOutcomesJSON, err := json.Marshal(tr.DeliveredOutcomes)
	if err != nil {
		return err
	}

	query := `	
		UPDATE transactions SET
			service_fee = $1,
			received_incomes = $2,
			delivered_outcomes = $3,
			received_company_id = $4,
			delivered_company_id = $5,
			received_user_id = $6,
			delivered_user_id = $7,
			phone = $8,
			details = $9,
			status = $10,
			type = $11
		WHERE id = $12 AND status = $13
	`

	result, err := s.db.ExecContext(
		ctx,
		query,
		tr.ServiceFee,
		receivedIncomesJSON,
		deliveredOutcomesJSON,
		tr.ReceivedCompanyId,
		tr.DeliveredCompanyId,
		tr.ReceivedUserId,
		tr.DeliveredUserId,
		tr.Phone,
		tr.Details,
		tr.Status,
		tr.Type,
		tr.ID,
		STATUS_CREATED,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("transaction to update not found")
	}

	return nil
}

func (s *TransactionStorage) GetById(ctx context.Context, id int64) (*Transaction, error) {
	query := `
				SELECT id, service_fee, received_incomes, delivered_outcomes,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE id = $1 AND status = $2 ORDER BY created_at DESC
			`

	tr := &Transaction{}
	var receivedIncomesJSON []byte
	var deliveredOutcomesJSON []byte

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
		STATUS_CREATED,
	).Scan(
		&tr.ID,
		&tr.ServiceFee,
		&receivedIncomesJSON,
		&deliveredOutcomesJSON,
		&tr.ReceivedCompanyId,
		&tr.DeliveredCompanyId,
		&tr.ReceivedUserId,
		&tr.DeliveredUserId,
		&tr.Phone,
		&tr.Details,
		&tr.Status,
		&tr.Type,
		&tr.CreatedAt)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(receivedIncomesJSON, &tr.ReceivedIncomes); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(deliveredOutcomesJSON, &tr.DeliveredOutcomes); err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Tashkent")
	createdAtInTashkent := tr.CreatedAt.In(loc)
	tr.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

	return tr, nil
}

func (s *TransactionStorage) Archived(ctx context.Context, pagination types.Pagination) ([]Transaction, error) {
	query := `
				SELECT id, service_fee, received_incomes, delivered_outcomes,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE status = $1   ORDER BY created_at DESC ` + fmt.Sprintf("OFFSET %v LIMIT %v", pagination.Offset, pagination.Limit)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		STATUS_ARCHIVED,
	)

	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetByField(ctx context.Context, fieldName string, fieldValue any, pagination types.Pagination) ([]Transaction, error) {
	query := `
				SELECT id, service_fee, received_incomes, delivered_outcomes,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE ` + fmt.Sprintf("%v", fieldName) + ` = $1 AND status != $2   ORDER BY created_at DESC ` + fmt.Sprintf("OFFSET %v LIMIT %v", pagination.Offset, pagination.Limit)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		fieldValue,
		STATUS_ARCHIVED,
	)

	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetInfos(ctx context.Context, companyId int64) ([]Transaction, error) {
	query := `
				SELECT id, service_fee, received_incomes, delivered_outcomes,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE delivered_company_id = $1 AND status = $2
			`
	rows, err := s.db.QueryContext(
		ctx,
		query,
		companyId,
		STATUS_CREATED,
	)

	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetByFieldAndDate(ctx context.Context, fieldName, from, to string, fieldValue any, pagination types.Pagination) ([]Transaction, error) {
	query := `
				SELECT id, service_fee, received_incomes, delivered_outcomes,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE ` + fmt.Sprintf("%v", fieldName) + ` = $1 AND created_at BETWEEN $2 AND $3 AND status != $4  ` + fmt.Sprintf("ORDER BY created_at DESC OFFSET %v LIMIT %v", pagination.Offset, pagination.Limit)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		fieldValue,
		from,
		to,
		STATUS_ARCHIVED,
	)

	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) Delete(ctx context.Context, id *int64) error {
	query := `DELETE FROM transactions WHERE id = $1`

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
		return errors.New("TRANSACTION NOT FOUND")
	}

	return nil
}

func (s *TransactionStorage) ConvertRowsToObject(rows *sql.Rows, err error) ([]Transaction, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []Transaction
	var receivedIncomesJSON []byte
	var deliveredOutcomesJSON []byte

	for rows.Next() {
		tr := &Transaction{}
		err := rows.Scan(
			&tr.ID,
			&tr.ServiceFee,
			&receivedIncomesJSON,
			&deliveredOutcomesJSON,
			&tr.ReceivedCompanyId,
			&tr.DeliveredCompanyId,
			&tr.ReceivedUserId,
			&tr.DeliveredUserId,
			&tr.Phone,
			&tr.Details,
			&tr.Status,
			&tr.Type,
			&tr.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(receivedIncomesJSON, &tr.ReceivedIncomes); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(deliveredOutcomesJSON, &tr.DeliveredOutcomes); err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("Asia/Tashkent")
		createdAtInTashkent := tr.CreatedAt.In(loc)
		tr.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

		transactions = append(transactions, *tr)
	}

	return transactions, nil
}

type CompanyAmount struct {
	CompanyName    string
	Currency       string
	OlinganAmount  float64
	BerilganAmount float64
	Remain         float64
}

func (s *TransactionStorage) GetCompanyFinalAmounts(ctx context.Context, companyIDs []int64, date string) ([]CompanyAmount, error) {
	query := `
with all_outcomes as (
    -- delivered_outcomes
    select 
        t.delivered_company_id as company_id,
        t.type,
        elem->>'delivered_currency' as currency,
        (elem->>'delivered_amount')::numeric as delivered_amount,
        0::numeric as received_amount,
        t.created_at
    from transactions t
    cross join jsonb_array_elements(t.delivered_outcomes) as elem
    where t.delivered_company_id = ANY($1) or t.received_company_id = ANY($1)

    union all

    -- received_incomes
    select 
        t.received_company_id as company_id,
        t.type,
        elem->>'received_currency' as currency,
        0::numeric as delivered_amount,
        (elem->>'received_amount')::numeric as received_amount,
        t.created_at
    from transactions t
    cross join jsonb_array_elements(t.received_incomes) as elem
    where t.delivered_company_id = ANY($1) or t.received_company_id = ANY($1)
)
select 
    c.name as company_name,
    a.currency,
    
    -- Olingan summalar (faqat date filter)
    coalesce(sum(
        case 
            when a.created_at::date = $2 then
                case 
                    when a.type = 1 then a.delivered_amount
                    when a.type = 2 then a.received_amount
                    else 0
                end
            else 0
        end
    ),0) as olingan_amount,

    -- Berilgan summalar (faqat date filter)
    coalesce(sum(
        case 
            when a.created_at::date = $2 then
                case 
                    when a.type = 1 then a.received_amount
                    when a.type = 2 then a.delivered_amount
                    else 0
                end
            else 0
        end
    ),0) as berilgan_amount,

    -- Qolgan summasi (barcha transactionlar)
    coalesce(sum(
        case 
            when a.type = 1 then a.delivered_amount
            when a.type = 2 then a.received_amount
            else 0
        end
    ),0)
    -
    coalesce(sum(
        case 
            when a.type = 1 then a.received_amount
            when a.type = 2 then a.delivered_amount
            else 0
        end
    ),0) as remain

from all_outcomes a
join companies c on c.id = a.company_id
group by a.company_id, c.name, a.currency
order by a.company_id, a.currency;
`

	rows, err := s.db.QueryContext(ctx, query, pq.Array(companyIDs), date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CompanyAmount
	for rows.Next() {
		var ca CompanyAmount
		if err := rows.Scan(&ca.CompanyName, &ca.Currency, &ca.OlinganAmount, &ca.BerilganAmount, &ca.Remain); err != nil {
			return nil, err
		}
		results = append(results, ca)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
