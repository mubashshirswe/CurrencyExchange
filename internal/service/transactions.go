package service

import (
	"context"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type TransactionService struct {
	store store.Storage
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *store.Transaction) error {

	return nil
}
