package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	serialNo := GenerateSerialNo(time.Hour.Microseconds())
	balanceRecord.SerialNo = serialNo

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
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING TRANSACTION %v", err)
	}

	err = s.store.BalanceRecords.Create(ctx, balanceRecord)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE CREATING BALANCE RECORD %v", err)
	}

	return nil
}

func (s *BalanceRecordService) RollbackBalanceRecord(ctx context.Context, serialNo string) error {
	balanceRecord, err := s.store.BalanceRecords.GetBySerialNo(ctx, serialNo)
	if err != nil {
		return fmt.Errorf("BALANCE RECORD IS NOT FOUND, %v", err)
	}

	balance, err := s.store.Balances.GetById(ctx, &balanceRecord.BalanceID)
	if err != nil {
		return fmt.Errorf("BALANCE SELECTED IS NOT FOUND, %v", err)
	}

	if balance.CurrencyId != balanceRecord.CurrenctID {
		return errors.New("CURRENCY IS NOT MATCH TO BALANCE CURRENCY")
	}

	if balanceRecord.Type == SELL {
		balance.Balance += balanceRecord.Amount
		balance.InOutLay -= balanceRecord.Amount
	} else {
		if balance.Balance >= balanceRecord.Amount {
			balance.Balance -= balanceRecord.Amount
			balance.OutInLay -= balanceRecord.Amount
		} else {
			return fmt.Errorf("BALANCE AMOUNT IS NOT ENOUGH TO BACK TRANSACTION")
		}
	}

	if err := s.store.BalanceRecords.Delete(ctx, balanceRecord.ID); err != nil {
		return err
	}

	return nil
}

func (s *BalanceRecordService) UpdateRecord(ctx context.Context, balanceRecord *store.BalanceRecord) error {

	if err := s.RollbackBalanceRecord(ctx, balanceRecord.SerialNo); err != nil {
		return err
	}

	if err := s.PerformBalanceRecord(ctx, balanceRecord); err != nil {
		return err
	}

	return nil
}
