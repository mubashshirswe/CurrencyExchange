package service

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type DebtsService struct {
	store store.Storage
}

func (s *DebtsService) Create(ctx context.Context, debt *store.Debts) error {
	if len(debt.ReceivedIncomes) == 0 {
		return fmt.Errorf("received incomes cannot be empty")
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	debtsStorage := store.NewDebtsStorage(tx)

	user, err := s.store.Users.GetById(ctx, &debt.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	debtor := &store.Debtors{
		FullName:  debt.FullName,
		Balance:   0,
		Currency:  debt.DebtedCurrency,
		UserID:    debt.UserID,
		CompanyID: user.CompanyId,
		Phone:     debt.Phone,
	}

	if err := debtorsStorage.Create(ctx, debtor); err != nil {
		return fmt.Errorf("failed to create debtor: %w", err)
	}

	// Assume input DebtedAmount is positive; apply sign for storage
	originalPositiveDebted := debt.DebtedAmount // Keep original positive for debtor calc
	var signedDebtedAmount int64
	switch debt.Type {
	case types.TYPE_SELL:
		signedDebtedAmount = -debt.DebtedAmount
		debt.State = 1
	case types.TYPE_BUY:
		signedDebtedAmount = debt.DebtedAmount
	default:
		return fmt.Errorf("invalid debt type: %d", debt.Type)
	}

	debt.DebtedAmount = signedDebtedAmount
	debt.DebtorID = debtor.ID
	debt.CompanyID = user.CompanyId

	// Create debt record
	if err := debtsStorage.Create(ctx, debt); err != nil {
		return fmt.Errorf("failed to create debt: %w", err)
	}

	// Process incomes (assume all in same currency as debt for simplicity; validate if needed)
	for _, tr := range debt.ReceivedIncomes {
		if tr.ReceivedCurrency != debt.DebtedCurrency {
			return fmt.Errorf("income currency %s does not match debt currency %s", tr.ReceivedCurrency, debt.DebtedCurrency)
		}

		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debt.UserID, tr.ReceivedCurrency)
		if err != nil {
			return fmt.Errorf("failed to get balance for currency %s: %w", tr.ReceivedCurrency, err)
		}

		originalPositiveReceived := tr.ReceivedAmount // Keep positive
		var signedReceivedAmount int64
		switch debt.Type {
		case types.TYPE_SELL:
			if balance.Balance < originalPositiveReceived {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
			balance.Balance -= originalPositiveReceived
			balance.InOutLay += originalPositiveReceived
			signedReceivedAmount = -originalPositiveReceived
		case types.TYPE_BUY:
			balance.Balance += originalPositiveReceived
			balance.OutInLay += originalPositiveReceived
			signedReceivedAmount = originalPositiveReceived
		}

		record := &store.BalanceRecord{
			Amount:    signedReceivedAmount,
			UserID:    debt.UserID,
			CompanyID: balance.CompanyId,
			BalanceID: balance.ID,
			Type:      int64(debt.Type),
			Details:   debt.Details,
			Currency:  tr.ReceivedCurrency,
			DebtId:    &debt.ID,
		}

		if err := balanceRecordsStorage.Create(ctx, record); err != nil {
			return fmt.Errorf("failed to create balance record: %w", err)
		}

		if err := balancesStorage.Update(ctx, balance); err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}
	}

	// Update debtor balance: positive if debtor owes user
	if debt.Type == types.TYPE_SELL {
		debtor.Balance += originalPositiveDebted
	} else {
		debtor.Balance -= originalPositiveDebted
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("failed to update debtor: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (s *DebtsService) Transaction(ctx context.Context, debt *store.Debts) error {
	if len(debt.ReceivedIncomes) == 0 {
		return fmt.Errorf("received incomes cannot be empty")
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	debtsStorage := store.NewDebtsStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)

	debtor, err := debtorsStorage.GetById(ctx, debt.DebtorID)
	if err != nil {
		return fmt.Errorf("failed to get debtor: %w", err)
	}

	if debtor.Currency != debt.DebtedCurrency {
		return fmt.Errorf("debted currencies do not match: %s != %s", debtor.Currency, debt.DebtedCurrency)
	}

	originalPositiveDebted := debt.DebtedAmount // Positive input
	var signedDebtedAmount int64
	switch debt.Type {
	case types.TYPE_SELL:
		signedDebtedAmount = -debt.DebtedAmount
	case types.TYPE_BUY:
		signedDebtedAmount = debt.DebtedAmount
	default:
		return fmt.Errorf("invalid debt type: %d", debt.Type)
	}

	debt.DebtedAmount = signedDebtedAmount
	debt.CompanyID = debtor.CompanyID

	// Create transaction debt record
	if err := debtsStorage.Create(ctx, debt); err != nil {
		return fmt.Errorf("failed to create debt: %w", err)
	}

	// Process incomes
	for _, tr := range debt.ReceivedIncomes {
		if tr.ReceivedCurrency != debt.DebtedCurrency {
			return fmt.Errorf("income currency %s does not match debt currency %s", tr.ReceivedCurrency, debt.DebtedCurrency)
		}

		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debt.UserID, tr.ReceivedCurrency)
		if err != nil {
			return fmt.Errorf("failed to get balance for currency %s: %w", tr.ReceivedCurrency, err)
		}

		originalPositiveReceived := tr.ReceivedAmount
		var signedReceivedAmount int64
		switch debt.Type {
		case types.TYPE_SELL:
			if balance.Balance < originalPositiveReceived {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
			balance.Balance -= originalPositiveReceived
			balance.InOutLay += originalPositiveReceived
			signedReceivedAmount = -originalPositiveReceived
		case types.TYPE_BUY:
			balance.Balance += originalPositiveReceived
			balance.OutInLay += originalPositiveReceived
			signedReceivedAmount = originalPositiveReceived
		}

		record := &store.BalanceRecord{
			Amount:    signedReceivedAmount,
			UserID:    debt.UserID,
			CompanyID: balance.CompanyId,
			BalanceID: balance.ID,
			Type:      int64(debt.Type),
			Details:   debt.Details,
			Currency:  tr.ReceivedCurrency,
			DebtId:    &debt.ID,
		}

		if err := balanceRecordsStorage.Create(ctx, record); err != nil {
			return fmt.Errorf("failed to create balance record: %w", err)
		}

		if err := balancesStorage.Update(ctx, balance); err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}
	}

	// Update debtor balance
	if debt.Type == types.TYPE_SELL {
		debtor.Balance += originalPositiveDebted
	} else {
		debtor.Balance -= originalPositiveDebted
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("failed to update debtor: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (s *DebtsService) Update(ctx context.Context, debt *store.Debts) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	log.Println("DEBT UPDATE==============    ", debt)

	debtsStorage := store.NewDebtsStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	balanceStorage := store.NewBalanceStorage(tx)

	oldDebt, err := debtsStorage.GetByID(ctx, debt.ID)
	if err != nil {
		return fmt.Errorf("failed to get old debt: %w", err)
	}

	debtor, err := debtorsStorage.GetById(ctx, oldDebt.DebtorID)
	if err != nil {
		return fmt.Errorf("failed to get debtor: %w", err)
	}

	// Use old currency for reversal
	balance, err := balanceStorage.GetByUserIdAndCurrency(ctx, &debt.UserID, oldDebt.DebtedCurrency)
	if err != nil {
		return fmt.Errorf("failed to get balance for old currency %s: %w", oldDebt.DebtedCurrency, err)
	}

	debt.CompanyID = balance.CompanyId

	// Reverse old effects: add opposite of stored signed amounts
	for _, oldTr := range oldDebt.ReceivedIncomes {
		// Reverse balance effect: balance += -old_signed_received
		balance.Balance -= oldTr.ReceivedAmount // Since Amount is signed, this adds the opposite
		if oldDebt.Type == types.TYPE_SELL {
			balance.InOutLay -= (-oldTr.ReceivedAmount) // Undo: for negative old, - positive
		} else {
			balance.OutInLay -= oldTr.ReceivedAmount // Undo positive
		}
	}

	// Reverse debtor effect
	originalOldPositive := math.Abs(float64(oldDebt.DebtedAmount))
	if oldDebt.Type == types.TYPE_SELL {
		debtor.Balance -= int64(originalOldPositive) // Undo + positive
	} else {
		debtor.Balance += int64(originalOldPositive) // Undo - positive
	}

	// Delete old records
	if err := balanceRecordsStorage.DeleteByDebtId(ctx, oldDebt.ID); err != nil {
		return fmt.Errorf("failed to delete old balance records: %w", err)
	}

	// Apply new effects (similar to Create)
	originalPositiveDebted := debt.DebtedAmount // Positive input
	var signedDebtedAmount int64
	switch debt.Type {
	case types.TYPE_SELL:
		signedDebtedAmount = -debt.DebtedAmount
	case types.TYPE_BUY:
		signedDebtedAmount = debt.DebtedAmount
	default:
		return fmt.Errorf("invalid new debt type: %d", debt.Type)
	}
	debt.DebtedAmount = signedDebtedAmount

	if len(debt.ReceivedIncomes) == 0 {
		return fmt.Errorf("new received incomes cannot be empty")
	}

	for _, tr := range debt.ReceivedIncomes {
		if tr.ReceivedCurrency != debt.DebtedCurrency {
			return fmt.Errorf("new income currency %s does not match new debt currency %s", tr.ReceivedCurrency, debt.DebtedCurrency)
		}

		originalPositiveReceived := tr.ReceivedAmount
		var signedReceivedAmount int64
		switch debt.Type {
		case types.TYPE_SELL:
			if balance.Balance < originalPositiveReceived {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
			balance.Balance -= originalPositiveReceived
			balance.InOutLay += originalPositiveReceived
			signedReceivedAmount = -originalPositiveReceived
		case types.TYPE_BUY:
			balance.Balance += originalPositiveReceived
			balance.OutInLay += originalPositiveReceived
			signedReceivedAmount = originalPositiveReceived
		}

		record := &store.BalanceRecord{
			Amount:    signedReceivedAmount,
			UserID:    debt.UserID,
			CompanyID: balance.CompanyId,
			BalanceID: balance.ID,
			Type:      int64(debt.Type),
			Details:   debt.Details,
			Currency:  tr.ReceivedCurrency,
			DebtId:    &debt.ID,
		}

		if err := balanceRecordsStorage.Create(ctx, record); err != nil {
			return fmt.Errorf("failed to create new balance record: %w", err)
		}
	}

	// Apply new debtor effect
	if debt.Type == types.TYPE_SELL {
		debtor.Balance += originalPositiveDebted
	} else {
		debtor.Balance -= originalPositiveDebted
	}

	if err := debtsStorage.Update(ctx, debt); err != nil {
		return fmt.Errorf("failed to update debt: %w", err)
	}

	if err := balanceStorage.Update(ctx, balance); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("failed to update debtor: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}

func (s *DebtsService) Delete(ctx context.Context, debtId int64) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	balancesStorage := store.NewBalanceStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	debtsStorage := store.NewDebtsStorage(tx)

	debt, err := debtsStorage.GetByID(ctx, debtId)
	if err != nil {
		return fmt.Errorf("failed to get debt: %w", err)
	}

	if err := balanceRecordsStorage.DeleteByDebtId(ctx, debtId); err != nil {
		return fmt.Errorf("failed to delete balance records: %w", err)
	}

	debtor, err := debtorsStorage.GetById(ctx, debt.DebtorID)
	if err != nil {
		return fmt.Errorf("failed to get debtor: %w", err)
	}

	balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debt.UserID, debt.DebtedCurrency)
	if err != nil {
		return fmt.Errorf("failed to get balance for currency %s: %w", debt.DebtedCurrency, err)
	}

	// Reverse effects: add opposite of stored signed amounts
	for _, tr := range debt.ReceivedIncomes {
		balance.Balance -= tr.ReceivedAmount // Opposite of stored signed
		if debt.Type == types.TYPE_SELL {
			balance.InOutLay -= (-tr.ReceivedAmount) // Undo for SELL
		} else {
			balance.OutInLay -= tr.ReceivedAmount // Undo for BUY
		}
	}

	// Reverse debtor effect
	originalPositive := math.Abs(float64(debt.DebtedAmount))
	if debt.Type == types.TYPE_SELL {
		debtor.Balance -= int64(originalPositive) // Undo + positive
	} else {
		debtor.Balance += int64(originalPositive) // Undo - positive
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("failed to update debtor: %w", err)
	}

	if err := debtsStorage.Delete(ctx, debtId); err != nil {
		return fmt.Errorf("failed to delete debt: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	return nil
}
