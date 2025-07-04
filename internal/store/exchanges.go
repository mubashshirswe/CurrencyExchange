package store

import (
	"context"
	"fmt"
)

type Exchange struct {
	ID               int64   `json:"id"`
	ReceivedMoney    int64   `json:"received_money"`
	ReceivedCurrency string  `json:"received_currency"`
	SelledMoney      int64   `json:"selled_money"`
	SelledCurrency   string  `json:"selled_currency"`
	UserId           int64   `json:"user_id"`
	Details          *string `json:"details"`
	CompanyID        int64   `json:"company_id"`
	CreatedAt        *string `json:"created_at"`
}

type ExchangeStorage struct {
	db DBTX
}

func NewExchangeStorage(db DBTX) *ExchangeStorage {
	return &ExchangeStorage{db: db}
}

func (s *ExchangeStorage) Create(ctx context.Context, exchange *Exchange) error {
	query := `INSERT INTO exchanges(received_money, received_currency, selled_money, selled_currency, user_id, company_id, details)
				VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		exchange.ReceivedMoney,
		exchange.ReceivedCurrency,
		exchange.SelledMoney,
		exchange.SelledCurrency,
		exchange.UserId,
		exchange.CompanyID,
		exchange.Details,
	).Scan(
		&exchange.ID,
		&exchange.CreatedAt,
	)

	return err
}

func (s *ExchangeStorage) Update(ctx context.Context, exchange *Exchange) error {
	query := `
				UPDATE exchanges SET received_money = $1, received_currency = $2, 
				selled_money = $3, selled_currency = $4, user_id = $5, company_id = $6, 
				details = $7 WHERE id  = $8`

	rows, err := s.db.ExecContext(
		ctx,
		query,
		exchange.ReceivedMoney,
		exchange.ReceivedCurrency,
		exchange.SelledMoney,
		exchange.SelledCurrency,
		exchange.UserId,
		exchange.CompanyID,
		exchange.Details,
		exchange.ID,
	)
	if err != nil {
		return err
	}
	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return fmt.Errorf("NOT FOUND")
	}

	return err
}

func (s *ExchangeStorage) GetByField(ctx context.Context, fieldName string, fieldValue any) ([]Exchange, error) {
	query := `
				SELECT id, received_money, received_currency, selled_money,
				selled_currency, user_id, company_id, details, created_at 
				FROM exchanges WHERE ` + fmt.Sprintf(" %v = %v ", fieldName, fieldValue)

	rows, err := s.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exchanges []Exchange
	for rows.Next() {
		exchage := &Exchange{}
		err := rows.Scan(
			&exchage.ID,
			&exchage.ReceivedMoney,
			&exchage.ReceivedCurrency,
			&exchage.SelledMoney,
			&exchage.SelledCurrency,
			&exchage.UserId,
			&exchage.CompanyID,
			&exchage.Details,
			&exchage.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		exchanges = append(exchanges, *exchage)
	}

	return exchanges, nil
}

func (s *ExchangeStorage) GetById(ctx context.Context, id int64) (*Exchange, error) {
	query := `
				SELECT id, received_money, received_currency, selled_money, 
				selled_currency, user_id, company_id, details, created_at 
				FROM exchanges WHERE id = $1`

	exchage := &Exchange{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&exchage.ID,
		&exchage.ReceivedMoney,
		&exchage.ReceivedCurrency,
		&exchage.SelledMoney,
		&exchage.SelledCurrency,
		&exchage.UserId,
		&exchage.CompanyID,
		&exchage.Details,
		&exchage.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return exchage, nil
}

func (s *ExchangeStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM exchanges WHERE id  = $1`

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
		return fmt.Errorf("NOT FOUND")
	}

	return nil
}
