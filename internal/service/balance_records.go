package service

import (
	"context"
	"fmt"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type BalanceRecordService struct {
	store store.Storage
}

func (s *BalanceRecordService) PerformBalanceRecord(ctx context.Context, balanceRecord types.BalanceRecordPayload) error {
	receivedCurrencyBalance, err := s.store.Balances.GetByIdAndCurrency(ctx, &balanceRecord.UserId, balanceRecord.ReceivedCurrency)
	if err != nil {
		return fmt.Errorf("BALANCE WITH CURRENCY %v NOT FOUND", balanceRecord.ReceivedCurrency)
	}

	selledCurrencyBalance, err := s.store.Balances.GetByIdAndCurrency(ctx, &balanceRecord.UserId, balanceRecord.SelledCurrency)
	if err != nil {
		return fmt.Errorf("BALANCE WITH CURRENCY %v NOT FOUND", balanceRecord.SelledCurrency)
	}

	if err := s.SelledMoneyPerform(ctx, balanceRecord, selledCurrencyBalance); err != nil {
		return err
	}

	if err := s.ReceivedMoneyPerform(ctx, balanceRecord, receivedCurrencyBalance); err != nil {
		return err
	}

	return nil
}

func (s *BalanceRecordService) RollbackBalanceRecord(ctx context.Context, id int64) error {
	records, err := s.store.BalanceRecords.GetByField(ctx, "id", id)
	if err != nil {
		return err
	}
	record := records[0]

	balance, err := s.store.Balances.GetById(ctx, &record.BalanceID)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE GETTING BALANCE WITH ID %v", err)
	}

	switch record.Type {
	case TYPE_SELL:
		balance.Balance += record.Amount
		balance.InOutLay -= record.Amount
	case TYPE_BUY:
		if balance.Balance >= record.Amount {
			balance.Balance -= record.Amount
			balance.OutInLay -= record.Amount
		} else {
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO ROLLBACK TRANSACTION %v >= %v", balance.Balance, record.Amount)
		}
	default:
		return fmt.Errorf("FOUND UNKOWN RECORD TYPE")
	}

	if err := s.store.Balances.Update(ctx, balance); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING BALANCE %v", err)
	}

	if err := s.store.BalanceRecords.Delete(ctx, record.ID); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE DELETING BALANCE RECORD %v", err)
	}

	return nil
}

func (s *BalanceRecordService) UpdateRecord(ctx context.Context, balanceRecord store.BalanceRecord) error {
	record, err := s.store.BalanceRecords.GetByField(ctx, "id", &balanceRecord.ID)
	if err != nil {
		return err
	}
	oldRecord := record[0]

	balance, err := s.store.Balances.GetById(ctx, &oldRecord.BalanceID)
	if err != nil {
		return err
	}

	switch oldRecord.Type {
	case TYPE_SELL:
		balance.Balance += oldRecord.Amount
		balance.InOutLay -= oldRecord.Amount
	case TYPE_BUY:
		if balance.Balance >= oldRecord.Amount {
			balance.Balance -= oldRecord.Amount
			balance.OutInLay -= oldRecord.Amount

			balance.Balance += balanceRecord.Amount
			balance.OutInLay += balanceRecord.Amount
		} else {
			return fmt.Errorf("oldRecord BALANCE  HAS NO ENOUGH MONEY")
		}
	}

	switch balanceRecord.Type {
	case TYPE_SELL:
		if balance.Balance >= balanceRecord.Amount {
			balance.Balance -= balanceRecord.Amount
			balance.InOutLay += balanceRecord.Amount
		} else {
			return fmt.Errorf("balanceRecord BALANCE  HAS NO ENOUGH MONEY")
		}
	case TYPE_BUY:
		balance.Balance += balanceRecord.Amount
		balance.OutInLay += balanceRecord.Amount
	}

	if err := s.store.Balances.Update(ctx, balance); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING BALANCE %v", err)
	}

	if err := s.store.BalanceRecords.Update(ctx, &balanceRecord); err != nil {
		return err
	}

	return nil
}

func (s *BalanceRecordService) ReceivedMoneyPerform(ctx context.Context, balanceRecord types.BalanceRecordPayload, receivedCurrencyBalance *store.Balance) error {

	if receivedCurrencyBalance != nil {
		receivedCurrencyBalance.Balance += balanceRecord.ReceivedMoney
		receivedCurrencyBalance.OutInLay += balanceRecord.ReceivedMoney
	}

	err := s.store.Balances.Update(ctx, receivedCurrencyBalance)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING receivedCurrencyBalance %v", err)
	}

	receivedMoneyRecord := &store.BalanceRecord{
		Amount:    balanceRecord.ReceivedMoney,
		Currency:  balanceRecord.ReceivedCurrency,
		CompanyID: balanceRecord.CompanyID,
		BalanceID: receivedCurrencyBalance.ID,
		Details:   &balanceRecord.Details,
		UserID:    balanceRecord.UserId,
		Type:      TYPE_BUY,
	}

	err = s.store.BalanceRecords.Create(ctx, receivedMoneyRecord)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE CREATING BALANCE RECORD %v", err)
	}

	return nil
}

func (s *BalanceRecordService) SelledMoneyPerform(ctx context.Context, balanceRecord types.BalanceRecordPayload, selledCurrencyBalance *store.Balance) error {
	if selledCurrencyBalance.Balance >= balanceRecord.ReceivedMoney {
		selledCurrencyBalance.Balance -= balanceRecord.SelledMoney
		selledCurrencyBalance.InOutLay += balanceRecord.SelledMoney
	} else {
		return fmt.Errorf("SELLED CURRENCY HAVE NO ENOUGH MONEY")
	}

	err := s.store.Balances.Update(ctx, selledCurrencyBalance)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING selledCurrencyBalance %v", err)
	}

	selledMoneyRecord := &store.BalanceRecord{
		Amount:    balanceRecord.SelledMoney,
		Currency:  balanceRecord.SelledCurrency,
		CompanyID: balanceRecord.CompanyID,
		BalanceID: selledCurrencyBalance.ID,
		Details:   &balanceRecord.Details,
		UserID:    balanceRecord.UserId,
		Type:      TYPE_SELL,
	}

	err = s.store.BalanceRecords.Create(ctx, selledMoneyRecord)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE CREATING BALANCE RECORD %v", err)
	}
	return nil
}
