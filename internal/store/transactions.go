package store

import (
	"context"
	"database/sql"
	"errors"
)

type Transaction struct {
	ID                 int64   `json:"id"`
	Amount             int64   `json:"amount"`
	ServiceFee         int64   `json:"service_fee"`
	FromCurrencyTypeId int64   `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64   `json:"to_currency_type_id"`
	SenderId           int64   `json:"sender_id"`
	FromCityId         int64   `json:"from_city_id"`
	ToCityId           int64   `json:"to_city_id"`
	ReceiverName       string  `json:"receiver_name"`
	ReceiverPhone      string  `json:"receiver_phone"`
	Details            string  `json:"details"`
	Type               int64   `json:"type"`
	CompanyId          int64   `json:"company_id"`
	BalanceId          int64   `json:"balance_id"`
	CreatedAt          *string `json:"created_at"`
}

type TransactionStorage struct {
	db *sql.DB
}

func (s *TransactionStorage) Create(ctx context.Context, tr *Transaction) error {
	query := `INSERT INTO transactions(
				amount, service_fee, from_currency_type_id,
				to_currency_type_id, sender_id, from_city_id, to_city_id,
				receiver_name, receiver_phone, details, type, company_id,
				created_at, balance_id) RETURNING id`

	err := s.db.QueryRowContext(
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
		tr.CreatedAt,
		tr.BalanceId,
	).Scan(
		&tr.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStorage) Update(ctx context.Context, tr *Transaction) error {
	query := `UPDATE transactions SET amount = $1, service_fee = $2, from_currency_type_id = $3,
				to_currency_type_id = $4, sender_id = $5, from_city_id = $6, to_city_id = $7,
				receiver_name = $8, receiver_phone = $9, details = $10, type = $11, 
				WHERE id = $12`

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
		return errors.New("NOT FOUND")
	}

	return nil
}

func (s *TransactionStorage) GetById(ctx context.Context, id *int64) (*Transaction, error) {
	query := `SELECT id, amount, service_fee, from_currency_type_id,
				to_currency_type_id, sender_id, from_city_id, to_city_id,
				receiver_name, receiver_phone, details, type, company_id,
				created_at, balance_id FROM transactions WHERE id = $1`

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
		&tr.CompanyId,
		&tr.CreatedAt,
		&tr.BalanceId,
	)

	if err != nil {
		return nil, err
	}

	return tr, nil
}

func (s *TransactionStorage) GetAll(ctx context.Context) ([]Transaction, error) {
	query := `SELECT id, amount, service_fee, from_currency_type_id,
				to_currency_type_id, sender_id, from_city_id, to_city_id,
				receiver_name, receiver_phone, details, type, company_id,
				created_at, balance_id FROM transactions ORDER BY created_at DESC`
	var transactions []Transaction
	rows, err := s.db.QueryContext(
		ctx,
		query,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			&tr.CompanyId,
			&tr.CreatedAt,
			&tr.BalanceId,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, *tr)
	}

	return transactions, nil
}

func (s *TransactionStorage) GetAllByDate(ctx context.Context, from string, to string) ([]Transaction, error) {
	query := `SELECT id, amount, service_fee, from_currency_type_id,
				to_currency_type_id, sender_id, from_city_id, to_city_id,
				receiver_name, receiver_phone, details, type, company_id, balance_id,
				created_at FROM transactions WHERE created_at between $1 and $2`
	var transactions []Transaction

	rows, err := s.db.QueryContext(
		ctx,
		query,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			&tr.CompanyId,
			&tr.BalanceId,
			&tr.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, *tr)
	}

	return transactions, nil
}
