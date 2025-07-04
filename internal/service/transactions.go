package service

import (
	"context"
	"database/sql"
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
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	transactionsStorage := store.NewTransactionStorage(tx)

	balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &transaction.ReceivedUserId, transaction.ReceivedCurrency)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return err
		} else {
			return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.GetByUserIdAndCurrency %v", err)
		}
	}

	switch transaction.Type {
	case TYPE_SELL:
		if balance.Balance >= transaction.ReceivedAmount {
			balance.Balance -= transaction.ReceivedAmount
			balance.InOutLay += transaction.ReceivedAmount
		} else {
			tx.Rollback()
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE THIS %v", err)
		}
	case TYPE_BUY:
		balance.Balance += transaction.ReceivedAmount
		balance.OutInLay += transaction.ReceivedAmount
	default:
		tx.Rollback()
		return fmt.Errorf("FOUND UNKNOWN TYPE")
	}

	if err := transactionsStorage.Create(ctx, transaction); err != nil {
		tx.Rollback()
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

	if err := balanceRecordsStorage.Create(ctx, balanceRecord); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.Create %v", err)
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.Create %v", err)
	}

	tx.Commit()
	return nil
}

func (s *TransactionService) CompleteTransaction(ctx context.Context, transaction types.TransactionComplete) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	transactionsStorage := store.NewTransactionStorage(tx)

	tran, err := transactionsStorage.GetById(ctx, transaction.TransactionID)
	if err != nil {
		tx.Rollback()
		return err
	}

	balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &transaction.DeliveredUserId, tran.DeliveredCurrency)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.GetByUserIdAndCurrency( %v", err)
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
			tx.Rollback()
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE THIS ACTION")
		}
	}

	balanceRecord := &store.BalanceRecord{
		Amount:        tran.DeliveredAmount,
		Currency:      tran.DeliveredCurrency,
		BalanceID:     balance.ID,
		UserID:        transaction.DeliveredUserId,
		TransactionId: &tran.ID,
		CompanyID:     tran.DeliveredCompanyId,
		Type:          recordType,
	}

	if err := balanceRecordsStorage.Create(ctx, balanceRecord); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE balanceRecordsStorage.Create %v", err)
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.Create %v", err)
	}

	tran.Status = 3
	tran.DeliveredServiceFee = &transaction.RecievedServiceFee

	if err := transactionsStorage.Update(ctx, tran); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE transactionsStorage.Update %v", err)
	}

	tx.Commit()
	return nil
}

func (s *TransactionService) Update(ctx context.Context, transaction *store.Transaction) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	transactionsStorage := store.NewTransactionStorage(tx)

	records, err := balanceRecordsStorage.GetByField(ctx, "transaction_id", transaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, record := range records {
		balance, err := balancesStorage.GetById(ctx, &record.BalanceID)
		if err != nil {
			tx.Rollback()
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
				tx.Rollback()
				return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY")
			}
		}

		if err := balanceRecordsStorage.Delete(ctx, record.ID); err != nil {
			tx.Rollback()
			return err
		}

		if err := balancesStorage.Update(ctx, balance); err != nil {
			tx.Rollback()
			return err
		}
	}

	if transaction.ReceivedUserId != 0 {
		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &transaction.ReceivedUserId, transaction.ReceivedCurrency)
		if err != nil {
			tx.Rollback()
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

		balance.Balance += transaction.ReceivedAmount
		balance.OutInLay += transaction.ReceivedAmount

		if err := balancesStorage.Update(ctx, balance); err != nil {
			tx.Rollback()
			return err
		}

		if err := balanceRecordsStorage.Update(ctx, record); err != nil {
			tx.Rollback()
			return err
		}
	}

	if transaction.DeliveredUserId != nil {
		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, transaction.DeliveredUserId, transaction.DeliveredCurrency)
		if err != nil {
			tx.Rollback()
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

		if err := balancesStorage.Update(ctx, balance); err != nil {
			tx.Rollback()
			return err
		}

		if err := balanceRecordsStorage.Update(ctx, record); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := transactionsStorage.Update(ctx, transaction); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING TRANSACTION %v", err)
	}

	tx.Commit()
	return nil
}

func (s *TransactionService) Delete(ctx context.Context, id *int64) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	transactionsStorage := store.NewTransactionStorage(tx)

	tran, err := transactionsStorage.GetById(ctx, *id)
	if err != nil {
		tx.Rollback()
		return err
	}

	records, err := balanceRecordsStorage.GetByField(ctx, "transaction_id", tran.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, record := range records {
		balance, err := balancesStorage.GetById(ctx, &record.BalanceID)
		if err != nil {
			tx.Rollback()
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
				tx.Rollback()
				return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE THE ACTION")
			}
		}
		if err := balancesStorage.Create(ctx, balance); err != nil {
			tx.Rollback()
			return err
		}
		if err := balanceRecordsStorage.Delete(ctx, record.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := transactionsStorage.Delete(ctx, &tran.ID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
