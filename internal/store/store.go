package store

import (
	"context"
	"database/sql"
)

const (
	STATUS_CREATED   = 1
	STATUS_COMPLETED = 2
	STATUS_ARCHIVED  = 3
)

type DBTX interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}

type Storage struct {
	DB *sql.DB

	Exchanges interface {
		Create(context.Context, *Exchange) error
		Update(context.Context, *Exchange) error
		GetById(context.Context, int64) (*Exchange, error)
		GetByField(context.Context, string, any) ([]Exchange, error)
		Delete(context.Context, int64) error
	}

	Debtors interface {
		Create(context.Context, *Debtors) error
		Update(context.Context, *Debtors) error
		GetById(context.Context, int64) (*Debtors, error)
		GetByUserId(context.Context, int64) ([]Debtors, error)
		GetByCompanyId(context.Context, int64) ([]Debtors, error)
		Delete(context.Context, int64) error
	}

	Users interface {
		Login(context.Context, *User) error
		Create(context.Context, *User) error
		Update(context.Context, *User) error
		GetAll(context.Context) ([]User, error)
		GetById(context.Context, *int64) (*User, error)
		Delete(context.Context, *int64) error
	}

	Balances interface {
		Create(context.Context, *Balance) error
		GetById(context.Context, *int64) (*Balance, error)
		GetByUserIdAndCurrency(context.Context, *int64, string) (*Balance, error)
		GetByUserId(context.Context, *int64) ([]Balance, error)
		GetAll(context.Context) ([]Balance, error)
		Update(context.Context, *Balance) error
		Delete(context.Context, int64) error
	}

	BalanceRecords interface {
		Create(context.Context, *BalanceRecord) error
		GetByField(context.Context, string, any) ([]BalanceRecord, error)
		GetByFieldAndDate(context.Context, string, string, string, any) ([]BalanceRecord, error)
		Update(context.Context, *BalanceRecord) error
		Delete(context.Context, int64) error
	}

	Transactions interface {
		Create(context.Context, *Transaction) error
		Update(context.Context, *Transaction) error
		Delete(context.Context, *int64) error
		GetByField(context.Context, string, any) ([]Transaction, error)
		GetByFieldAndDate(context.Context, string, string, string, any) ([]Transaction, error)
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
	dbwrapper := &DBWrapper{db: db}

	return Storage{
		DB:             db,
		Exchanges:      &ExchangeStorage{db: dbwrapper},
		Debtors:        &DebtorsStorage{db: dbwrapper},
		Users:          &UserStorage{db: dbwrapper},
		Transactions:   &TransactionStorage{db: dbwrapper},
		Balances:       &BalanceStorage{db: dbwrapper},
		Companies:      &CompanyStorage{db: dbwrapper},
		BalanceRecords: &BalanceRecordStorage{db: dbwrapper},
	}
}

func (s *Storage) BeginTx(ctx context.Context) (DBTX, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &TxWrapper{tx: tx}, nil
}
