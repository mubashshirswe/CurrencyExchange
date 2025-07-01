package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type Service struct {
	Debtors interface {
		Create(context.Context, *store.Debtors) error
		Transaction(context.Context, *store.Debtors) error
		GetByCompanyId(context.Context, int64) (map[string]interface{}, error)
		Update(context.Context, *store.BalanceRecord) error
		Delete(context.Context, int64) error
	}
	BalanceRecords interface {
		PerformBalanceRecord(context.Context, types.BalanceRecordPayload) error
		RollbackBalanceRecord(context.Context, int64) error
		UpdateRecord(context.Context, store.BalanceRecord) error
	}

	Transactions interface {
		PerformTransaction(context.Context, *store.Transaction) error
		CompleteTransaction(context.Context, types.TransactionComplete) error
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
