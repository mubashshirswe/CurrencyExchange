package store

import (
	"context"
	"database/sql"
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
	CreatedAt          *string `json:"created_at"`
}

type TransactionStorage struct {
	db *sql.DB
}

func (s *TransactionStorage) Create(ctx context.Context, tr *Transaction) error {
	query := `INSERT INTO transactions(
				amount, service_fee, from_currency_type_id,
				to_currency_type_id, sender_id, from_city_id,
				receiver_name, receiver_phone, details, type,
				created_at) RETURNING id`

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
		tr.CreatedAt,
	).Scan(
		&tr.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionStorage) GetById(ctx context.Context, id *int64) (*Transaction, error) {
	query := `SELECT * FROM transactions WHERE id =$1`
	tr := &Transaction{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		tr.ID,
	).Scan(
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
		&tr.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return tr, nil
}

func (s *TransactionStorage) GetAll(ctx context.Context) ([]Transaction, error) {
	query := `SELECT * FROM transactions`
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
		rows.Scan(
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
			&tr.CreatedAt,
		)

		transactions = append(transactions, *tr)
	}

	return transactions, nil
}

func (s *TransactionStorage) GetAllByDate(ctx context.Context, from string, to string) ([]Transaction, error) {
	query := `SELECT * FROM transactions WHERE created_at between $1 and $2`
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
		rows.Scan(
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
			&tr.CreatedAt,
		)

		transactions = append(transactions, *tr)
	}

	return transactions, nil
}
