package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type Service struct {
	Debtors interface {
		Update(context.Context, *store.Debtors) error
		Create(context.Context, *store.Debtors) error
		ReceivedDebt(context.Context, int64) error
		GetByUserId(context.Context, int64) ([]store.Debtors, error)
		Delete(context.Context, int64) error
	}
	BalanceRecords interface {
		PerformBalanceRecord(context.Context, *store.BalanceRecord) error
		RollbackBalanceRecord(context.Context, string) error
		UpdateRecord(context.Context, *store.BalanceRecord) error
	}

	Transactions interface {
		PerformTransaction(context.Context, *store.Transaction) error
		CompleteTransaction(context.Context, string) error
		Update(context.Context, *store.Transaction) error
		Delete(context.Context, *int64) error
	}
}

func NewService(store store.Storage) Service {
	return Service{
		BalanceRecords: &BalanceRecordService{store: store},
		Transactions:   &TransactionService{store: store},
		Debtors:        &DebtorsService{store: store},
	}
}

func GenerateSerialNo(id int64) string {
	return fmt.Sprintf("%v%v", id, rand.Intn(10000000))
}
