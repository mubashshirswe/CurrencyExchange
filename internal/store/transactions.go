package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Transaction struct {
	ID                  int64     `json:"id"`
	MarkedServiceFee    *int64    `json:"marked_service_fee"`
	ReceivedCompanyId   int64     `json:"received_company_id"`
	ReceivedUserId      int64     `json:"received_user_id"`
	ReceivedAmount      int64     `json:"received_amount"`
	ReceivedCurrency    string    `json:"received_currency"`
	DeliveredAmount     int64     `json:"delivered_amount"`
	DeliveredCurrency   string    `json:"delivered_currency"`
	DeliveredCompanyId  int64     `json:"delivered_company_id"`
	DeliveredUserId     *int64    `json:"delivered_user_id"`
	DeliveredServiceFee *int64    `json:"delivered_service_fee"`
	Phone               string    `json:"phone"`
	Details             string    `json:"details"`
	Status              int64     `json:"status"`
	Type                int64     `json:"type"`
	CreatedAt           time.Time `json:"-"`
	CreatedAtFormatted  string    `json:"created_at"`
}

type TransactionStorage struct {
	db DBTX
}

func NewTransactionStorage(db DBTX) *TransactionStorage {
	return &TransactionStorage{db: db}
}

func (s *TransactionStorage) Archive(ctx context.Context) error {
	query := `UPDATE transactions SET created_at = $1, status = $2 WHERE status = $3`
	rows, err := s.db.ExecContext(ctx, query, time.Now(), STATUS_ARCHIVED, STATUS_COMPLETED)
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
	query := `
			INSERT INTO transactions(
				marked_service_fee, delivered_service_fee, received_amount, received_currency, delivered_amount, delivered_currency,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		tr.MarkedServiceFee,
		tr.DeliveredServiceFee,
		tr.ReceivedAmount,
		tr.ReceivedCurrency,
		tr.DeliveredAmount,
		tr.DeliveredCurrency,
		tr.ReceivedCompanyId,
		tr.DeliveredCompanyId,
		tr.ReceivedUserId,
		tr.DeliveredUserId,
		tr.Phone,
		tr.Details,
		STATUS_CREATED,
		tr.Type,
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
	query := `	
				UPDATE transactions SET marked_service_fee = $1, delivered_service_fee = $2, received_amount = $3, received_currency = $4, 
				delivered_amount = $5, delivered_currency = $6, received_company_id = $7, delivered_company_id = $8, received_user_id = $9, 
				delivered_user_id = $10, phone = $11, details = $12, status = $13, type = $14 WHERE id = $15 AND status = $16`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		tr.MarkedServiceFee,
		tr.DeliveredServiceFee,
		tr.ReceivedAmount,
		tr.ReceivedCurrency,
		tr.DeliveredAmount,
		tr.DeliveredCurrency,
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
	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return errors.New("TRANSACTION THAT WILL BE UPDATED NOT FOUND")
	}

	return nil
}

func (s *TransactionStorage) GetById(ctx context.Context, id int64) (*Transaction, error) {
	query := `
				SELECT id, marked_service_fee, delivered_service_fee, received_amount, received_currency, delivered_amount, delivered_currency,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE id = $1 AND status = $2 ORDER BY created_at DESC
			`

	tr := &Transaction{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
		STATUS_CREATED,
	).Scan(
		&tr.ID,
		&tr.MarkedServiceFee,
		&tr.DeliveredServiceFee,
		&tr.ReceivedAmount,
		&tr.ReceivedCurrency,
		&tr.DeliveredAmount,
		&tr.DeliveredCurrency,
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

	loc, _ := time.LoadLocation("Asia/Tashkent")
	createdAtInTashkent := tr.CreatedAt.In(loc)
	tr.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

	return tr, nil
}

func (s *TransactionStorage) Archived(ctx context.Context) ([]Transaction, error) {
	query := `
				SELECT id, marked_service_fee, delivered_service_fee, received_amount, received_currency, delivered_amount, delivered_currency,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE status = $1 ORDER BY created_at DESC
			`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		STATUS_ARCHIVED,
	)

	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetByField(ctx context.Context, fieldName string, fieldValue any) ([]Transaction, error) {
	query := `
				SELECT id, marked_service_fee, delivered_service_fee, received_amount, received_currency, delivered_amount, delivered_currency,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE ` + fmt.Sprintf("%v", fieldName) + ` = $1 AND status != $2 ORDER BY created_at DESC
			`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		fieldValue,
		STATUS_ARCHIVED,
	)

	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetByFieldAndDate(ctx context.Context, fieldName, from, to string, fieldValue any) ([]Transaction, error) {
	query := `
				SELECT id, marked_service_fee, delivered_service_fee, received_amount, received_currency, delivered_amount, delivered_currency,
	 			received_company_id, delivered_company_id, received_user_id, delivered_user_id, phone, details, status, type, created_at
				FROM transactions WHERE ` + fmt.Sprintf("%v", fieldName) + ` = $1 AND created_at BETWEEN $2 AND $3 AND status != $4 ORDER BY created_at DESC
			`

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

	for rows.Next() {
		tr := &Transaction{}
		err := rows.Scan(
			&tr.ID,
			&tr.MarkedServiceFee,
			&tr.DeliveredServiceFee,
			&tr.ReceivedAmount,
			&tr.ReceivedCurrency,
			&tr.DeliveredAmount,
			&tr.DeliveredCurrency,
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

		loc, _ := time.LoadLocation("Asia/Tashkent")
		createdAtInTashkent := tr.CreatedAt.In(loc)
		tr.CreatedAtFormatted = createdAtInTashkent.Format("2006-01-02 15:04:05")

		transactions = append(transactions, *tr)
	}

	return transactions, nil
}
