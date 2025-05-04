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
		GetById(context.Context, *int64) (*User, error)
		Delete(context.Context, *int64) error
	}

	Employees interface {
		Login(context.Context, *Employee) error
		Create(context.Context, *Employee) error
		Update(context.Context, *Employee) error
		GetAll(context.Context) ([]Employee, error)
		GetById(context.Context, *int64) (*Employee, error)
		Delete(context.Context, *int64) error
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
		Update(context.Context, *Transaction) error
		GetById(context.Context, *int64) (*Transaction, error)
		GetAllByBalanceId(context.Context, *int64) ([]Transaction, error)
		GetAllByDate(context.Context, string, string, *int64) ([]Transaction, error)
	}

	Currencies interface {
		Create(context.Context, *Currency) error
		GetAll(context.Context) ([]Currency, error)
		Update(context.Context, *Currency) error
		Delete(context.Context, *int64) error
	}
	Cities interface {
		Create(context.Context, *City) error
		GetAll(context.Context) ([]City, error)
		Update(context.Context, *City) error
		Delete(context.Context, *int64) error
	}

	Companies interface {
		Create(context.Context, *Company) error
		GetAll(context.Context) ([]Company, error)
		GetById(context.Context, *int64) (*Company, error)
		Update(context.Context, *Company) error
		Delete(context.Context, *int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users:        &UserStorage{db: db},
		Transactions: &TransactionStorage{db: db},
		Currencies:   &CurrencyStorage{db: db},
		Cities:       &CityStorage{db: db},
		Balances:     &BalanceStorage{db: db},
		Companies:    &CompanyStorage{db: db},
		Employees:    &EmployeeStorage{db: db},
	}
}
