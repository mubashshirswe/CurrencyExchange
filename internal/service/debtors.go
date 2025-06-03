package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type DebtorsService struct {
	store store.Storage
}

func (s *DebtorsService) Create(ctx context.Context, debtor *store.Debtors) error {
	service := NewService(s.store)
	debtor.SerialNo = GenerateSerialNo(time.Hour.Microseconds())

	if debtor.IsBalanceEffect == 1 && debtor.Type == TYPE_SELL {
		if err := service.BalanceRecords.PerformBalanceRecord(ctx, ConvertDebitsDataToBalanceRecords(debtor)); err != nil {
			return err
		}
	}

	if err := s.store.Debtors.Create(ctx, debtor); err != nil {
		return err
	}

	return nil
}

func (s *DebtorsService) Update(ctx context.Context, debtor *store.Debtors) error {
	service := NewService(s.store)
	if debtor.IsBalanceEffect == 1 {
		if err := service.BalanceRecords.RollbackBalanceRecord(ctx, debtor.SerialNo); err != nil {
			return err
		}

		if err := service.BalanceRecords.PerformBalanceRecord(ctx, ConvertDebitsDataToBalanceRecords(debtor)); err != nil {
			return err
		}
	}

	return s.store.Debtors.Update(ctx, debtor)
}

func (s *DebtorsService) GetByUserId(ctx context.Context, id int64) ([]store.Debtors, error) {
	return s.store.Debtors.GetByUserId(ctx, id)
}

func (s *DebtorsService) Delete(ctx context.Context, id int64) error {
	service := NewService(s.store)
	debtor, err := s.store.Debtors.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE GET BY ID %e", err)
	}

	if debtor.IsBalanceEffect == 1 {
		if err := service.BalanceRecords.RollbackBalanceRecord(ctx, debtor.SerialNo); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE RollbackBalanceRecord %e", err)
		}
	}
	debtor.Status = -1
	return s.store.Debtors.Update(ctx, debtor)
}

func (s *DebtorsService) ReceivedDebt(ctx context.Context, id int64) error {
	service := NewService(s.store)
	debtor, err := s.store.Debtors.GetById(ctx, id)
	if err == nil {
		if debtor.IsBalanceEffect == 1 {
			if err := service.BalanceRecords.PerformBalanceRecord(ctx, ConvertDebitsDataToBalanceRecords(debtor)); err != nil {
				return err
			}
		}
	}

	debtor.Status = -1
	return s.store.Debtors.Update(ctx, debtor)
}

func ConvertDebitsDataToBalanceRecords(debtor *store.Debtors) *store.BalanceRecord {
	return &store.BalanceRecord{
		Amount:       debtor.Amount,
		UserID:       debtor.UserID,
		SerialNo:     debtor.SerialNo,
		BalanceID:    debtor.BalanceId,
		CompanyID:    debtor.CompanyId,
		Details:      debtor.Details,
		CurrenctID:   debtor.CurrencyId,
		Type:         int64(debtor.Type),
		CurrencyType: debtor.CurrencyType,
	}
}
