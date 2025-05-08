package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

const (
	SELL = 1
	BUY  = 1
)

type BalanceRecordService struct {
	store store.Storage
}

func (s *BalanceRecordService) PerformBalanceRecord(ctx context.Context, balanceRecord *store.BalanceRecord) error {

	balance, err := s.store.Balances.GetById(ctx, &balanceRecord.BalanceID)
	if err != nil {
		return fmt.Errorf("BALANCE IS SELECTED NOT FOUND, %v", err)
	}

	if balance.CurrencyId != balanceRecord.CurrenctID {
		return errors.New("CURRENCY IS NOT MATCH TO BALANCE CURRENCY")
	}

	if balanceRecord.Type == SELL {
		if balance.Balance < balanceRecord.Amount {
			return errors.New("BALANCE IS NOT ENOUGH TO TRANSACTION")
		} else {
			balance.Balance -= balanceRecord.Amount
			balance.InOutLay += balanceRecord.Amount
		}
	} else {
		balance.Balance += balanceRecord.Amount
		balance.OutInLay += balanceRecord.Amount
	}

	err = s.store.Balances.Update(ctx, balance)
	if err != nil {
		return err
	}

	err = s.store.BalanceRecords.Create(ctx, balanceRecord)
	if err != nil {
		return err
	}

	return nil
}
