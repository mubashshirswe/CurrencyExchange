package service

import (
	"context"
	"fmt"
	"log"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type DebtsService struct {
	store store.Storage
}

func (s *DebtsService) Create(ctx context.Context, debt *store.Debts) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	debtsStorage := store.NewDebtsStorage(tx)

	balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debt.UserID, debt.ReceivedCurrency)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetByUserIdAndCurrency")
	}
	debt.CompanyID = balance.CompanyId

	switch debt.Type {
	case TYPE_SELL:
		if balance.Balance < debt.ReceivedAmount {
			return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
		}
		balance.Balance -= debt.ReceivedAmount
		balance.InOutLay += debt.ReceivedAmount

		debt.ReceivedAmount = -debt.ReceivedAmount
		debt.DebtedAmount = -debt.DebtedAmount

		debt.State = 1
	case TYPE_BUY:
		balance.Balance += debt.ReceivedAmount
		balance.OutInLay += debt.ReceivedAmount
	}

	debtor := &store.Debtors{
		FullName:  debt.FullName,
		Balance:   debt.DebtedAmount,
		Currency:  debt.DebtedCurrency,
		UserID:    debt.UserID,
		CompanyID: debt.CompanyID,
		Phone:     debt.Phone,
	}

	if err := debtorsStorage.Create(ctx, debtor); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.Create %v", err)
	}

	debt.DebtorId = debtor.ID
	if err := debtsStorage.Create(ctx, debt); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtsStorage.Create %v", err)
	}

	record := &store.BalanceRecord{
		Amount:    debt.ReceivedAmount,
		UserID:    debt.UserID,
		CompanyID: balance.CompanyId,
		BalanceID: balance.ID,
		Type:      int64(debt.Type),
		Details:   debt.Details,
		Currency:  debt.ReceivedCurrency,
		DebtId:    &debt.ID,
	}

	if err := balanceRecordsStorage.Create(ctx, record); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balanceRecordsStorage.Create %v", err)
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.Update %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE tx.Commit: %w", err)
	}
	return nil
}

func (s *DebtsService) Transaction(ctx context.Context, debt *store.Debts) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	debtsStorage := store.NewDebtsStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)

	debtor, err := debtorsStorage.GetById(ctx, debt.DebtorId)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.GetById %v", err)
	}

	if debtor.Currency != debt.DebtedCurrency {
		return fmt.Errorf("DEBTED CURRENCIES ARE NOT MATCH %v", err)
	}

	balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debtor.UserID, debt.ReceivedCurrency)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetByUserIdAndCurrency %v", err)
	}

	debt.CompanyID = balance.CompanyId
	debtor.CompanyID = balance.CompanyId
	debt.UserID = debtor.UserID

	switch debt.Type {
	case TYPE_SELL:
		if balance.Balance >= debt.ReceivedAmount {
			balance.Balance -= debt.ReceivedAmount
			balance.InOutLay += debt.ReceivedAmount

			debt.ReceivedAmount = -debt.ReceivedAmount
			debt.DebtedAmount = -debt.DebtedAmount
			debtor.Balance += debt.DebtedAmount
		} else {
			return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
		}
	case TYPE_BUY:
		balance.Balance += debt.ReceivedAmount
		balance.OutInLay += debt.ReceivedAmount

		debtor.Balance += debt.DebtedAmount
	}

	if err := debtsStorage.Create(ctx, debt); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtsStorage.Create %v", err)
	}

	record := &store.BalanceRecord{
		Amount:    debt.ReceivedAmount,
		UserID:    debt.UserID,
		CompanyID: balance.CompanyId,
		BalanceID: balance.ID,
		Type:      int64(debt.Type),
		Details:   debt.Details,
		Currency:  debt.ReceivedCurrency,
		DebtId:    &debt.ID,
	}

	if err := balanceRecordsStorage.Create(ctx, record); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.Create %v", err)
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.Update %v", err)
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.Update %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE tx.Commit: %w", err)
	}
	return nil
}

func (s *DebtsService) Update(ctx context.Context, debt *store.Debts) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	log.Println("DEBT UPDATE==============    ", debt)

	debtsStorage := store.NewDebtsStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	balanceStorage := store.NewBalanceStorage(tx)

	old, err := debtsStorage.GetById(ctx, debt.ID)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtsStorage.GetById %v", err)
	}

	debtor, err := debtorsStorage.GetById(ctx, old.DebtorId)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.GetById %v", err)
	}

	balance, err := balanceStorage.GetByUserIdAndCurrency(ctx, &debtor.UserID, debtor.Currency)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balanceStorage.GetByUserIdAndCurrency %v", err)
	}

	debt.CompanyID = balance.CompanyId

	switch old.Type {
	case TYPE_SELL:
		balance.Balance += old.ReceivedAmount
		balance.InOutLay -= old.ReceivedAmount

		debtor.Balance -= old.DebtedAmount
	case TYPE_BUY:
		if balance.Balance >= old.ReceivedAmount {
			balance.Balance -= old.ReceivedAmount
			balance.OutInLay -= old.ReceivedAmount

			debtor.Balance -= old.DebtedAmount
		} else {
			return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
		}
	}

	switch debt.Type {
	case TYPE_SELL:
		if balance.Balance >= debt.ReceivedAmount {
			balance.Balance -= debt.ReceivedAmount
			balance.InOutLay += debt.ReceivedAmount

			if debt.ReceivedAmount > 0 {
				debt.ReceivedAmount = -debt.ReceivedAmount
			}
			if debt.DebtedAmount > 0 {
				debt.DebtedAmount = -debt.DebtedAmount
			}

			debtor.Balance += debt.DebtedAmount
		} else {
			return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
		}

	case TYPE_BUY:
		balance.Balance += debt.ReceivedAmount
		balance.OutInLay += debt.ReceivedAmount

		debtor.Balance += debt.DebtedAmount
	}

	if err := debtsStorage.Update(ctx, debt); err != nil {
		return fmt.Errorf("ERROR OCCURRED debtsStorage.Create( %v", err)
	}

	record := &store.BalanceRecord{
		Amount:    debt.ReceivedAmount,
		UserID:    debt.UserID,
		CompanyID: balance.CompanyId,
		BalanceID: balance.ID,
		Type:      int64(debt.Type),
		Details:   debt.Details,
		Currency:  debt.ReceivedCurrency,
		DebtId:    &debt.ID,
	}

	if err := balanceRecordsStorage.DeleteByDebtId(ctx, old.ID); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balanceRecordsStorage.DeleteByDebtId %v", err)
	}

	if err := balanceRecordsStorage.Create(ctx, record); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balanceRecordsStorage.Create %v", err)
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.Update %v", err)
	}

	if err := balanceStorage.Update(ctx, balance); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balanceStorage.Update %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE tx.Commit: %w", err)
	}
	return nil
}

func (s *DebtsService) Delete(ctx context.Context, debtId int64) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	balancesStorage := store.NewBalanceStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)
	debtsStorage := store.NewDebtsStorage(tx)

	if err := balanceRecordsStorage.DeleteByDebtId(ctx, debtId); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE  balanceRecordsStorage.DeleteByDebtId %v", err)
	}

	debt, err := debtsStorage.GetById(ctx, debtId)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.GetById(ctx, *record.DebtorId) %v", err)
	}

	debtor, err := debtorsStorage.GetById(ctx, debt.DebtorId)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.GetById(ctx, *record.DebtorId) %v", err)
	}

	balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debtor.UserID, debtor.Currency)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetById %v", err)
	}

	switch debt.Type {
	case TYPE_SELL:
		balance.Balance += debt.ReceivedAmount
		balance.InOutLay -= debt.ReceivedAmount

		debtor.Balance -= debt.DebtedAmount
	case TYPE_BUY:
		if balance.Balance >= debt.ReceivedAmount {
			balance.Balance -= debt.ReceivedAmount
			balance.OutInLay -= debt.ReceivedAmount

			debtor.Balance -= debt.DebtedAmount
		} else {
			return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
		}
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return err
	}

	if err := debtsStorage.Delete(ctx, debtId); err != nil {
		return err
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE tx.Commit: %w", err)
	}
	return nil
}
