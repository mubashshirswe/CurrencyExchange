package service

import (
	"context"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type Service struct {
	BalanceRecords interface {
		PerformBalanceRecord(context.Context, *store.BalanceRecord) error
	}

	Transactions interface {
		CreateTransaction(context.Context, *store.Transaction) error
	}
}

func NewService(store store.Storage) Service {
	return Service{
		BalanceRecords: &BalanceRecordService{store: store},
		Transactions:   &TransactionService{store: store},
	}
}
