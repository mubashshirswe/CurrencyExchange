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
		return err
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

	isFlag := false
	var balanceId int64
	for _, o := range balances {
		if o.CurrencyId == transaction.ToCurrencyTypeId && o.Balance > transaction.Amount {
			isFlag = true
			balanceId = o.ID
		}
	}

	if !isFlag {
		return nil, fmt.Errorf("RECEIVER HAS NO BALANCE ACCOUNT OR BALANCE AMOUNT LIKE REQUESTED")
	}

	return &balanceId, nil
}
