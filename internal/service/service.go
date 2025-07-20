package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type Service struct {
	Exchanges interface {
		Create(context.Context, *store.Exchange) error
		Update(context.Context, *store.Exchange) error
		Delete(context.Context, int64) error
	}

	Balances interface {
		GetByCompanyId(context.Context, int64) ([]map[string]interface{}, error)
		GetAll(context.Context) ([]map[string]interface{}, error)
	}

	Debts interface {
		Create(context.Context, *store.Debts) error
		Transaction(context.Context, *store.Debts) error
		Update(context.Context, *store.Debts) error
		Delete(context.Context, int64) error
	}
	BalanceRecords interface {
		PerformBalanceRecord(context.Context, types.BalanceRecordPayload) error
		RollbackBalanceRecord(context.Context, int64) error
		UpdateRecord(context.Context, store.BalanceRecord) error
	}

	Transactions interface {
		GetByField(context.Context, string, any) ([]map[string]interface{}, error)
		PerformTransaction(context.Context, *store.Transaction) error
		CompleteTransaction(context.Context, types.TransactionComplete) error
		GetByCompanyId(context.Context, int64) ([]map[string]interface{}, error)
		GetInfos(context.Context, int64) ([]map[string]interface{}, error)
		Archived(context.Context) ([]map[string]interface{}, error)
		Update(context.Context, *store.Transaction) error
		Delete(context.Context, *int64) error
	}
}

func NewService(store store.Storage) Service {
	return Service{
		Balances:       &BalanceService{store: store},
		Exchanges:      &ExchangeService{store: store},
		BalanceRecords: &BalanceRecordService{store: store},
		Transactions:   &TransactionService{store: store},
		Debts:          &DebtsService{store: store},
	}
}

func GenerateSerialNo(id int64) string {
	return fmt.Sprintf("%v%v", id, rand.Intn(10000000))
}
