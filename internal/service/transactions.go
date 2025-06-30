package service

import (
	"context"
	"fmt"

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
	balance, err := s.store.Balances.GetByIdAndCurrency(ctx, &transaction.ReceivedUserId, transaction.ReceivedCurrency)
	if err != nil {
		return err
	}

	switch transaction.Type {
	case TYPE_SELL:
		if balance.Balance >= transaction.ReceivedAmount {
			balance.Balance -= transaction.ReceivedAmount
			balance.InOutLay += transaction.ReceivedAmount
		} else {
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE THIS %v", err)
		}
	case TYPE_BUY:
		balance.Balance += transaction.ReceivedAmount
		balance.OutInLay += transaction.ReceivedAmount
	default:
		return fmt.Errorf("FOUND UNKNOWN TYPE")
	}

	if err := s.store.Transactions.Create(ctx, transaction); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Transactions.Create %v", err)
	}

	balanceRecord := &store.BalanceRecord{
		Amount:        transaction.ReceivedAmount,
		Currency:      transaction.ReceivedCurrency,
		BalanceID:     balance.ID,
		CompanyID:     balance.CompanyId,
		UserID:        transaction.ReceivedUserId,
		Type:          transaction.Type,
		TransactionId: &transaction.ID,
	}

	if err := s.store.BalanceRecords.Create(ctx, balanceRecord); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.Create %v", err)
	}

	return nil
}

func (s *TransactionService) CompleteTransaction(ctx context.Context, serialNo string) error {

	return nil
}

func (s *TransactionService) Update(ctx context.Context, transaction *store.Transaction) error {

	if err := s.store.Transactions.Update(ctx, transaction); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING TRANSACTION %v", err)
	}

	return nil
}

func (s *TransactionService) Delete(ctx context.Context, id *int64) error {

	if err := s.store.Transactions.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
