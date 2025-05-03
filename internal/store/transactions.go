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
	ReceivedTime       *string `json:"received_time"`
	DeliveredTime      *string `json:"delivered_time"`
}

type TransactionStorage struct {
	db *sql.DB
}

func (s *TransactionStorage) Create(ctx context.Context, tr *Transaction) error {

	return nil
}

func (s *TransactionStorage) GetById(ctx context.Context, id *int64) (*Transaction, error) {

	return nil, nil
}

func (s *TransactionStorage) GetAll(ctx context.Context) ([]Transaction, error) {

	return nil, nil
}

func (s *TransactionStorage) GetAllByDate(ctx context.Context, from string, to string) ([]Transaction, error) {

	return nil, nil
}
