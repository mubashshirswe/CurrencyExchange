package service

import (
	"context"
	"fmt"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
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

func (s *TransactionService) CompleteTransaction(ctx context.Context, transaction types.TransactionComplete) error {
	transactions, err := s.store.Transactions.GetByField(ctx, "id", transaction.TransactionID)
	if err != nil {
		return err
	}
	tran := transactions[0]

	balance, err := s.store.Balances.GetByIdAndCurrency(ctx, tran.DeliveredUserId, tran.DeliveredCurrency)
	if err != nil {
		return err
	}

	var recordType int64
	if tran.Type == TYPE_SELL {
		recordType = TYPE_BUY
		balance.Balance += tran.DeliveredAmount
		balance.OutInLay += tran.ReceivedAmount
	} else {
		recordType = TYPE_SELL
		if balance.Balance >= tran.DeliveredAmount {
			balance.Balance -= tran.DeliveredAmount
			balance.InOutLay += tran.ReceivedAmount
		} else {
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE THIS ACTION")
		}
	}

	balanceRecord := &store.BalanceRecord{
		Amount:        tran.DeliveredAmount,
		Currency:      tran.DeliveredCurrency,
		BalanceID:     balance.ID,
		UserID:        *tran.DeliveredUserId,
		TransactionId: &tran.ID,
		Type:          recordType,
	}

	if err := s.store.BalanceRecords.Create(context.Background(), balanceRecord); err != nil {
		return err
	}

	if err := s.store.Balances.Create(context.Background(), balance); err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) Update(ctx context.Context, transaction *store.Transaction) error {
	records, err := s.store.BalanceRecords.GetByField(ctx, "transaction_id", transaction.ID)
	if err != nil {
		return err
	}

	for _, record := range records {
		balance, err := s.store.Balances.GetById(ctx, &record.BalanceID)
		if err != nil {
			return err
		}

		if record.Type == TYPE_SELL {
			balance.Balance += record.Amount
			balance.InOutLay -= record.Amount
		} else {
			if balance.Balance >= record.Amount {
				balance.Balance -= record.Amount
				balance.OutInLay -= record.Amount
			} else {
				return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY")
			}
		}

		if err := s.store.BalanceRecords.Delete(ctx, record.ID); err != nil {
			return err
		}

		if err := s.store.Balances.Update(ctx, balance); err != nil {
			return err
		}
	}

	if transaction.ReceivedUserId != 0 {
		balance, err := s.store.Balances.GetByIdAndCurrency(ctx, &transaction.ReceivedUserId, transaction.ReceivedCurrency)
		if err != nil {
			return err
		}

		record := &store.BalanceRecord{
			Amount:        transaction.ReceivedAmount,
			Currency:      transaction.ReceivedCurrency,
			UserID:        transaction.ReceivedUserId,
			CompanyID:     balance.CompanyId,
			TransactionId: &transaction.ID,
			BalanceID:     balance.ID,
			Details:       &transaction.Details,
			Type:          transaction.Type,
		}

		if err := s.store.Balances.Update(ctx, balance); err != nil {
			return err
		}

		if err := s.store.BalanceRecords.Update(ctx, record); err != nil {
			return err
		}
	}

	if transaction.DeliveredUserId != nil {
		balance, err := s.store.Balances.GetByIdAndCurrency(ctx, transaction.DeliveredUserId, transaction.DeliveredCurrency)
		if err != nil {
			return err
		}

		record := &store.BalanceRecord{
			Amount:        transaction.DeliveredAmount,
			Currency:      transaction.DeliveredCurrency,
			UserID:        *transaction.DeliveredUserId,
			CompanyID:     balance.CompanyId,
			TransactionId: &transaction.ID,
			BalanceID:     balance.ID,
			Details:       &transaction.Details,
			Type:          transaction.Type,
		}

		if err := s.store.Balances.Update(ctx, balance); err != nil {
			return err
		}

		if err := s.store.BalanceRecords.Update(ctx, record); err != nil {
			return err
		}
	}

	if err := s.store.Transactions.Update(ctx, transaction); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING TRANSACTION %v", err)
	}

	return nil
}

func (s *TransactionService) Delete(ctx context.Context, id *int64) error {
	trans, err := s.store.Transactions.GetByField(ctx, "id", id)
	if err != nil {
		return err
	}
	tran := trans[0]

	records, err := s.store.BalanceRecords.GetByField(ctx, "transaction_id", tran.ID)
	if err != nil {
		return err
	}

	for _, record := range records {
		balance, err := s.store.Balances.GetById(ctx, &record.BalanceID)
		if err != nil {
			return err
		}

		if record.Type == TYPE_SELL {
			balance.Balance += record.Amount
			balance.InOutLay -= record.Amount
		} else {
			if balance.Balance >= record.Amount {
				balance.Balance -= record.Amount
				balance.OutInLay -= record.Amount
			} else {
				return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE THE ACTION")
			}
		}
		if err := s.store.Balances.Create(ctx, balance); err != nil {
			return err
		}
		if err := s.store.BalanceRecords.Delete(ctx, record.ID); err != nil {
			return err
		}
	}

	if err := s.store.Transactions.Delete(ctx, &tran.ID); err != nil {
		return err
	}

	return nil
}
