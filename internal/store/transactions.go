package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Transaction struct {
	ID                 int64     `json:"id"`
	Amount             int64     `json:"amount"`
	ServiceFee         int64     `json:"service_fee"`
	FromCurrencyTypeId int64     `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64     `json:"to_currency_type_id"`
	SenderId           int64     `json:"sender_id"`
	SerialNo           string    `json:"serial_no"`
	ReceiverId         int64     `json:"receiver_id"`
	FromCityId         int64     `json:"from_city_id"`
	ToCityId           int64     `json:"to_city_id"`
	ReceiverName       string    `json:"receiver_name"`
	ReceiverPhone      string    `json:"receiver_phone"`
	Details            string    `json:"details"`
	Type               int64     `json:"type"`
	FromCurrencyType   string    `json:"from_currency_type"`
	ToCurrencyType     string    `json:"to_currency_type"`
	Status             int64     `json:"status"`
	CompanyId          int64     `json:"company_id"`
	BalanceId          int64     `json:"balance_id"`
	CreatedAt          time.Time `json:"created_at"`
}

type TransactionStorage struct {
	db *sql.DB
}

func (s *TransactionStorage) Create(ctx context.Context, tr *Transaction) error {
	query := `INSERT INTO transactions(
                amount, service_fee, from_currency_type_id,
                to_currency_type_id, sender_id, from_city_id, to_city_id,
                receiver_name, receiver_phone, details, type, company_id,
                balance_id, status, receiver_id, serial_no) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err := s.db.QueryContext(
		ctx,
		query,
		tr.Amount,
		tr.ServiceFee,
		tr.FromCurrencyTypeId,
		tr.ToCurrencyTypeId,
		tr.SenderId,
		tr.FromCityId,
		tr.ToCityId,
		tr.ReceiverName,
		tr.ReceiverPhone,
		tr.Details,
		tr.Type,
		tr.CompanyId,
		tr.BalanceId,
		tr.Status,
		tr.ReceiverId,
		tr.SerialNo,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStorage) Update(ctx context.Context, tr *Transaction) error {
	query := `UPDATE transactions SET amount = $1, service_fee = $2, from_currency_type_id = $3,
				to_currency_type_id = $4, sender_id = $5, from_city_id = $6, to_city_id = $7,
				receiver_name = $8, receiver_phone = $9, details = $10, type = $11, receiver_id = $12, serial_no = $13, status = $14
				WHERE id = $15`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		tr.Amount,
		tr.ServiceFee,
		tr.FromCurrencyTypeId,
		tr.ToCurrencyTypeId,
		tr.SenderId,
		tr.FromCityId,
		tr.ToCityId,
		tr.ReceiverName,
		tr.ReceiverPhone,
		tr.Details,
		tr.Type,
		tr.ReceiverId,
		tr.SerialNo,
		tr.Status,
		tr.ID,
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

func (s *TransactionStorage) GetByField(ctx context.Context, fieldName string, fieldValue string) (*Transaction, error) {
	query := `SELECT 
    t.id, 
    t.amount, 
    t.service_fee, 
    t.from_currency_type_id,
    t.to_currency_type_id, 
    t.sender_id, 
    t.from_city_id, 
    t.to_city_id,
    t.receiver_name, 
    t.receiver_phone, 
    t.details, 
    t.type, 
	t.status, 
    t.company_id,
    t.created_at, 
    t.balance_id, 
    t.receiver_id, 
	t.serial_no, 
    cf.name AS from_currency_name,
    ct.name AS to_currency_name
FROM 
    transactions t
LEFT JOIN 
    currencies cf ON t.from_currency_type_id = cf.id
LEFT JOIN 
    currencies ct ON t.to_currency_type_id = ct.id
WHERE ` + fmt.Sprintf("%v", fieldName) + ` = $1
ORDER BY 
    t.serial_no DESC`

	tr := &Transaction{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		fieldValue,
	).Scan(
		&tr.ID,
		&tr.Amount,
		&tr.ServiceFee,
		&tr.FromCurrencyTypeId,
		&tr.ToCurrencyTypeId,
		&tr.SenderId,
		&tr.FromCityId,
		&tr.ToCityId,
		&tr.ReceiverName,
		&tr.ReceiverPhone,
		&tr.Details,
		&tr.Type,
		&tr.Status,
		&tr.CompanyId,
		&tr.CreatedAt,
		&tr.BalanceId,
		&tr.ReceiverId,
		&tr.SerialNo,
		&tr.FromCurrencyType,
		&tr.ToCurrencyType,
	)

	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Tashkent")
	createdAtInTashkent := tr.CreatedAt.In(loc)
	tr.CreatedAt = createdAtInTashkent

	return tr, nil
}

func (s *TransactionStorage) GetAllByStatus(ctx context.Context, status int64) ([]Transaction, error) {
	query := `SELECT 
    t.id, 
    t.amount, 
    t.service_fee, 
    t.from_currency_type_id,
    t.to_currency_type_id, 
    t.sender_id, 
    t.from_city_id, 
    t.to_city_id,
    t.receiver_name, 
    t.receiver_phone, 
    t.details, 
    t.type, 
	t.status, 
    t.company_id,
    t.created_at, 
    t.balance_id, 
    t.receiver_id, 
	t.serial_no, 
    cf.name AS from_currency_name,
    ct.name AS to_currency_name
FROM 
    transactions t
LEFT JOIN 
    currencies cf ON t.from_currency_type_id = cf.id
LEFT JOIN 
    currencies ct ON t.to_currency_type_id = ct.id
WHERE t.status = $1
ORDER BY 
    t.created_at DESC`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		status,
	)

	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetById(ctx context.Context, id *int64) (*Transaction, error) {
	query := `SELECT 
    t.id, 
    t.amount, 
    t.service_fee, 
    t.from_currency_type_id,
    t.to_currency_type_id, 
    t.sender_id, 
    t.from_city_id, 
    t.to_city_id,
    t.receiver_name, 
    t.receiver_phone, 
    t.details, 
    t.type, 
	t.status, 
    t.company_id,
    t.created_at, 
    t.balance_id, 
    t.receiver_id, 
	t.serial_no, 
    cf.name AS from_currency_name,
    ct.name AS to_currency_name
FROM 
    transactions t
LEFT JOIN 
    currencies cf ON t.from_currency_type_id = cf.id
LEFT JOIN 
    currencies ct ON t.to_currency_type_id = ct.id
WHERE 
    t.id = $1
ORDER BY 
    t.created_at DESC`

	tr := &Transaction{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&tr.ID,
		&tr.Amount,
		&tr.ServiceFee,
		&tr.FromCurrencyTypeId,
		&tr.ToCurrencyTypeId,
		&tr.SenderId,
		&tr.FromCityId,
		&tr.ToCityId,
		&tr.ReceiverName,
		&tr.ReceiverPhone,
		&tr.Details,
		&tr.Type,
		&tr.Status,
		&tr.CompanyId,
		&tr.CreatedAt,
		&tr.BalanceId,
		&tr.ReceiverId,
		&tr.FromCurrencyType,
		&tr.ToCurrencyType,
	)

	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Tashkent")
	createdAtInTashkent := tr.CreatedAt.In(loc)
	tr.CreatedAt = createdAtInTashkent

	return tr, nil
}

func (s *TransactionStorage) GetAllByBalanceId(ctx context.Context, balance_id *int64) ([]Transaction, error) {
	query := `SELECT 
    t.id, 
    t.amount, 
    t.service_fee, 
    t.from_currency_type_id,
    t.to_currency_type_id, 
    t.sender_id, 
    t.from_city_id, 
    t.to_city_id,
    t.receiver_name, 
    t.receiver_phone, 
    t.details, 
    t.type, 
	t.status, 
    t.company_id,
    t.created_at, 
    t.balance_id, 
    t.receiver_id, 
	t.serial_no, 
    cf.name AS from_currency_name,
    ct.name AS to_currency_name
FROM 
    transactions t
LEFT JOIN 
    currencies cf ON t.from_currency_type_id = cf.id
LEFT JOIN 
    currencies ct ON t.to_currency_type_id = ct.id
WHERE 
    t.balance_id = $1
ORDER BY 
    t.created_at DESC`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		balance_id,
	)
	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetAllByUserId(ctx context.Context, userId *int64) ([]Transaction, error) {
	query := `SELECT 
    t.id, 
    t.amount, 
    t.service_fee, 
    t.from_currency_type_id,
    t.to_currency_type_id, 
    t.sender_id, 
    t.from_city_id, 
    t.to_city_id,
    t.receiver_name, 
    t.receiver_phone, 
    t.details, 
    t.type, 
	t.status, 
    t.company_id,
    t.created_at, 
    t.balance_id, 
    t.receiver_id, 
	t.serial_no, 
    cf.name AS from_currency_name,
    ct.name AS to_currency_name
FROM 
    transactions t
LEFT JOIN 
    currencies cf ON t.from_currency_type_id = cf.id
LEFT JOIN 
    currencies ct ON t.to_currency_type_id = ct.id
WHERE 
    t.sender_id = $1
ORDER BY 
    t.created_at DESC`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
	)
	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetAllByReceiverId(ctx context.Context, receiverId *int64) ([]Transaction, error) {
	query := `SELECT 
    t.id, 
    t.amount, 
    t.service_fee, 
    t.from_currency_type_id,
    t.to_currency_type_id, 
    t.sender_id, 
    t.from_city_id, 
    t.to_city_id,
    t.receiver_name, 
    t.receiver_phone, 
    t.details, 
    t.type, 
	t.status, 
    t.company_id,
    t.created_at, 
    t.balance_id, 
    t.receiver_id,
	t.serial_no, 
    cf.name AS from_currency_name,
    ct.name AS to_currency_name
FROM 
    transactions t
LEFT JOIN 
    currencies cf ON t.from_currency_type_id = cf.id
LEFT JOIN 
    currencies ct ON t.to_currency_type_id = ct.id
WHERE 
    t.receiver_id = $1
ORDER BY 
    t.created_at DESC`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		receiverId,
	)
	return s.ConvertRowsToObject(rows, err)
}

func (s *TransactionStorage) GetAllByDate(ctx context.Context, from string, to string, balance_id *int64) ([]Transaction, error) {
	query := `SELECT 
    t.id, 
    t.amount, 
    t.service_fee, 
    t.from_currency_type_id,
    t.to_currency_type_id, 
    t.sender_id, 
    t.from_city_id, 
    t.to_city_id,
    t.receiver_name, 
    t.receiver_phone, 
    t.details, 
    t.type, 
	t.status, 
    t.company_id,
    t.created_at, 
    t.balance_id, 
    t.receiver_id, 
	t.serial_no, 
    cf.name AS from_currency_name,
    ct.name AS to_currency_name
FROM 
    transactions t
LEFT JOIN 
    currencies cf ON t.from_currency_type_id = cf.id
LEFT JOIN 
    currencies ct ON t.to_currency_type_id = ct.id
WHERE 
t.created_at BETWEEN $1 and $2
    t.balance_id = $3
ORDER BY 
    t.created_at DESC`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		from,
		to,
		balance_id,
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
			&tr.Amount,
			&tr.ServiceFee,
			&tr.FromCurrencyTypeId,
			&tr.ToCurrencyTypeId,
			&tr.SenderId,
			&tr.FromCityId,
			&tr.ToCityId,
			&tr.ReceiverName,
			&tr.ReceiverPhone,
			&tr.Details,
			&tr.Type,
			&tr.Status,
			&tr.CompanyId,
			&tr.CreatedAt,
			&tr.BalanceId,
			&tr.ReceiverId,
			&tr.SerialNo,
			&tr.FromCurrencyType,
			&tr.ToCurrencyType,
		)

		if err != nil {
			return nil, err
		}

		loc, _ := time.LoadLocation("Asia/Tashkent")
		createdAtInTashkent := tr.CreatedAt.In(loc)
		tr.CreatedAt = createdAtInTashkent

		transactions = append(transactions, *tr)
	}

	return transactions, nil
}
