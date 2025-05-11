package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

const (
	TRANSACTION_STATUS_PENDING   = 1
	TRANSACTION_STATUS_COMPLETED = 2
	TYPE_SELL                    = 1
	TYPE_BUY                     = 2
)

type TransactionService struct {
	store store.Storage
}

func (s *TransactionService) PerformTransaction(ctx context.Context, transaction *store.Transaction) error {
	service := NewService(s.store)
	serialNo := GenerateSerialNo(time.Hour.Microseconds())

	if _, err := s.CheckReceiverBalance(ctx, transaction); err != nil {
		return err
	}

	balanceRecord := &store.BalanceRecord{
		SerialNo:   serialNo,
		Amount:     transaction.Amount,
		UserID:     transaction.SenderId,
		BalanceID:  transaction.BalanceId,
		CompanyID:  transaction.CompanyId,
		Details:    transaction.Details,
		CurrenctID: transaction.FromCurrencyTypeId,
		Type:       TYPE_SELL,
	}

	if err := service.BalanceRecords.PerformBalanceRecord(ctx, balanceRecord); err != nil {
		return err
	}

	transaction.Status = TRANSACTION_STATUS_PENDING
	transaction.SerialNo = serialNo

	if err := s.store.Transactions.Create(ctx, transaction); err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) CompleteTransaction(ctx context.Context, transaction *store.Transaction) error {
	service := NewService(s.store)
	serialNo := GenerateSerialNo(time.Hour.Microseconds())

	balanceId, err := s.CheckReceiverBalance(ctx, transaction)
	if err != nil {
		return err
	}

	balanceRecord := &store.BalanceRecord{
		SerialNo:   serialNo,
		Amount:     transaction.Amount,
		UserID:     transaction.ReceiverId,
		BalanceID:  *balanceId,
		CompanyID:  transaction.CompanyId,
		Details:    transaction.Details,
		CurrenctID: transaction.FromCurrencyTypeId,
		Type:       TYPE_BUY,
	}

	if err := service.BalanceRecords.PerformBalanceRecord(ctx, balanceRecord); err != nil {
		return err
	}

	transaction.Status = TRANSACTION_STATUS_COMPLETED

	if err := s.store.Transactions.Update(ctx, transaction); err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) Update(ctx context.Context, transaction *store.Transaction) error {
	service := NewService(s.store)

	balanceRecord := &store.BalanceRecord{
		SerialNo:   transaction.SerialNo,
		Amount:     transaction.Amount,
		UserID:     transaction.SenderId,
		BalanceID:  transaction.BalanceId,
		CompanyID:  transaction.CompanyId,
		Details:    transaction.Details,
		CurrenctID: transaction.FromCurrencyTypeId,
		Type:       transaction.Type,
	}

	if err := service.BalanceRecords.UpdateRecord(ctx, balanceRecord); err != nil {
		return err
	}

	if err := s.store.Transactions.Update(ctx, transaction); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING TRANSACTION %v", err)
	}

	return nil
}

func (s *TransactionService) Delete(ctx context.Context, id *int64) error {
	service := NewService(s.store)

	transaction, err := s.store.Transactions.GetById(ctx, id)
	if err != nil {
		return err
	}

	if err := service.BalanceRecords.RollbackBalanceRecord(ctx, transaction.SerialNo); err != nil {
		return err
	}

	if err := s.store.Transactions.Delete(ctx, &transaction.ID); err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) CheckReceiverBalance(ctx context.Context, transaction *store.Transaction) (*int64, error) {
	balances, err := s.store.Balances.GetByUserId(ctx, &transaction.ReceiverId)
	if err != nil {
		return nil, fmt.Errorf("ERROR OCCURRED WHILE CHECKING RECEIVER BALANCE INFO %v", err)
	}

	if balances == nil {
		return nil, fmt.Errorf("BALANCE NOT FOUND WITH ID %v", transaction.ReceiverId)
	}

	var isFlag bool
	var balanceId *int64
	for _, o := range balances {
		if o.CurrencyId == transaction.ToCurrencyTypeId {
			isFlag = true
			if o.Balance > transaction.Amount {
				balanceId = &o.ID
			}
		}
	}

	if !isFlag {
		return nil, fmt.Errorf("RECEIVER BALANCE CURRENCY DO NOT MATCH")
	}

	if balanceId == nil {
		return nil, fmt.Errorf("RECEIVER BALANCE HAVE NO ENOUGH MONEY OR BALANCE NOT FOUND")
	}

	return balanceId, nil
}
