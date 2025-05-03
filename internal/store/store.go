package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		Login(context.Context, *User) error
		Create(context.Context, *User) error
		Update(context.Context, *User) error
		GetAll(context.Context) ([]User, error)
	}

	Balances interface {
		Create(context.Context, *Balance) error
		GetById(context.Context, *int64) (*Balance, error)
		GetByUserId(context.Context, *int64) ([]Balance, error)
		GetAll(context.Context) ([]Balance, error)
		Update(context.Context, *Balance) error
	}

	Transactions interface {
		Create(context.Context, *Transaction) error
		GetById(context.Context, *int64) (*Transaction, error)
		GetAll(context.Context) ([]Transaction, error)
		GetAllByDate(context.Context, string, string) ([]Transaction, error)
	}

	Currencies interface {
		Create(context.Context, *Currency) error
		GetAll(context.Context) ([]Currency, error)
	}
	Cities interface {
		Create(context.Context, *City) error
		GetAll(context.Context) ([]City, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:        &UserStorage{db: db},
		Transactions: &TransactionStorage{db: db},
		Currencies:   &CurrencyStorage{db: db},
		Cities:       &CityStorage{db: db},
		Balances:     &BalanceStorage{db: db},
	}
}
