package service

import (
	"context"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type DebtorsService struct {
	store store.Storage
}

func (s *DebtorsService) Create(ctx context.Context, debtor *store.Debtors) error {
	return s.store.Debtors.Create(ctx, debtor)
}

func (s *DebtorsService) Update(ctx context.Context, debtor *store.Debtors) error {
	return s.store.Debtors.Update(ctx, debtor)
}

func (s *DebtorsService) GetByUserId(ctx context.Context, id int64) ([]store.Debtors, error) {
	return s.store.Debtors.GetByUserId(ctx, id)
}

func (s *DebtorsService) Delete(ctx context.Context, id int64) error {
	return s.store.Debtors.Delete(ctx, id)
}
