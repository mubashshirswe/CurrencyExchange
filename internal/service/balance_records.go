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
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)

	receivedCurrencyBalance, err := balancesStorage.GetByIdAndCurrency(ctx, &balanceRecord.UserId, balanceRecord.ReceivedCurrency)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("BALANCE WITH CURRENCY %v NOT FOUND", balanceRecord.ReceivedCurrency)
	}

	selledCurrencyBalance, err := balancesStorage.GetByIdAndCurrency(ctx, &balanceRecord.UserId, balanceRecord.SelledCurrency)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("BALANCE WITH CURRENCY %v NOT FOUND", balanceRecord.SelledCurrency)
	}

	/// SELLED MONEY PEFPERFORMROM
	if selledCurrencyBalance.Balance >= balanceRecord.ReceivedMoney {
		selledCurrencyBalance.Balance -= balanceRecord.SelledMoney
		selledCurrencyBalance.InOutLay += balanceRecord.SelledMoney
	} else {
		tx.Rollback()
		return fmt.Errorf("SELLED CURRENCY HAVE NO ENOUGH MONEY")
	}

	if err := balancesStorage.Update(ctx, selledCurrencyBalance); err != nil {
		tx.Rollback()
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

	if err := balanceRecordsStorage.Create(ctx, selledMoneyRecord); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE CREATING BALANCE RECORD %v", err)
	}

	// RECEIVED MONEY PERFORM
	if receivedCurrencyBalance != nil {
		receivedCurrencyBalance.Balance += balanceRecord.ReceivedMoney
		receivedCurrencyBalance.OutInLay += balanceRecord.ReceivedMoney
	}

	if err := balancesStorage.Update(ctx, receivedCurrencyBalance); err != nil {
		tx.Rollback()
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

	if err := balanceRecordsStorage.Create(ctx, receivedMoneyRecord); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE CREATING BALANCE RECORD %v", err)
	}

	tx.Commit()
	return nil
}

func (s *BalanceRecordService) RollbackBalanceRecord(ctx context.Context, id int64) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)

	records, err := balanceRecordsStorage.GetByField(ctx, "id", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	record := records[0]

	balance, err := balancesStorage.GetById(ctx, &record.BalanceID)
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO ROLLBACK TRANSACTION %v >= %v", balance.Balance, record.Amount)
		}
	default:
		tx.Rollback()
		return fmt.Errorf("FOUND UNKOWN RECORD TYPE")
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING BALANCE %v", err)
	}

	if err := balanceRecordsStorage.Delete(ctx, record.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE DELETING BALANCE RECORD %v", err)
	}

	tx.Commit()
	return nil
}

func (s *BalanceRecordService) UpdateRecord(ctx context.Context, balanceRecord store.BalanceRecord) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)

	record, err := balanceRecordsStorage.GetByField(ctx, "id", &balanceRecord.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	oldRecord := record[0]

	balance, err := balancesStorage.GetById(ctx, &oldRecord.BalanceID)
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return fmt.Errorf("oldRecord BALANCE  HAS NO ENOUGH MONEY")
		}
	}

	switch balanceRecord.Type {
	case TYPE_SELL:
		if balance.Balance >= balanceRecord.Amount {
			balance.Balance -= balanceRecord.Amount
			balance.InOutLay += balanceRecord.Amount
		} else {
			tx.Rollback()
			return fmt.Errorf("balanceRecord BALANCE  HAS NO ENOUGH MONEY")
		}
	case TYPE_BUY:
		balance.Balance += balanceRecord.Amount
		balance.OutInLay += balanceRecord.Amount
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING BALANCE %v", err)
	}

	if err := balanceRecordsStorage.Update(ctx, &balanceRecord); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
