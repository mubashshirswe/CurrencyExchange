package store

import (
	"context"
	"database/sql"
)

type Storage struct {
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
		GetByIdAndCurrency(context.Context, *int64, string) (*Balance, error)
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
	return Storage{
		Debtors:        &DebtorsStorage{db: db},
		Users:          &UserStorage{db: db},
		Transactions:   &TransactionStorage{db: db},
		Balances:       &BalanceStorage{db: db},
		Companies:      &CompanyStorage{db: db},
		BalanceRecords: &BalanceRecordStorage{db: db},
	}
}

type TxWrapper struct {
	tx *sql.Tx
}

func (t *TxWrapper) Commit() error {
	return t.tx.Commit()
}

func (t *TxWrapper) Rollback() error {
	return t.tx.Rollback()
}

func (t *TxWrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *TxWrapper) ExecContext(ctx context.Context, query string, args ...any) (*sql.Result, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	return &result, err
}
