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

	user, _ := s.store.Users.GetById(ctx, &debt.UserID)

	debtor := &store.Debtors{
		FullName:  debt.FullName,
		Balance:   debt.DebtedAmount,
		Currency:  debt.DebtedCurrency,
		UserID:    debt.UserID,
		CompanyID: user.CompanyId,
		Phone:     debt.Phone,
	}

	if err := debtorsStorage.Create(ctx, debtor); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.Create %v", err)
	}

	for _, tr := range debt.ReceivedIncomes {
		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debt.UserID, tr.ReceivedCurrency)
		if err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetByUserIdAndCurrency")
		}

		switch debt.Type {
		case TYPE_SELL:
			if balance.Balance < tr.ReceivedAmount {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
			balance.Balance -= tr.ReceivedAmount
			balance.InOutLay += tr.ReceivedAmount

			tr.ReceivedAmount = -tr.ReceivedAmount
			debt.DebtedAmount = -debt.DebtedAmount

			debt.State = 1
		case TYPE_BUY:
			balance.Balance += tr.ReceivedAmount
			balance.OutInLay += tr.ReceivedAmount
		}

		debt.DebtorID = debtor.ID
		if err := debtsStorage.Create(ctx, debt); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE debtsStorage.Create %v", err)
		}

		record := &store.BalanceRecord{
			Amount:    tr.ReceivedAmount,
			UserID:    debt.UserID,
			CompanyID: balance.CompanyId,
			BalanceID: balance.ID,
			Type:      int64(debt.Type),
			Details:   debt.Details,
			Currency:  tr.ReceivedCurrency,
			DebtId:    &debt.ID,
		}

		if err := balanceRecordsStorage.Create(ctx, record); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE balanceRecordsStorage.Create %v", err)
		}

		if err := balancesStorage.Update(ctx, balance); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.Update %v", err)
		}
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

	debtor, err := debtorsStorage.GetById(ctx, debt.DebtorID)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.GetById %v", err)
	}

	if debtor.Currency != debt.DebtedCurrency {
		return fmt.Errorf("DEBTED CURRENCIES ARE NOT MATCH %v", err)
	}

	for _, tr := range debt.ReceivedIncomes {

		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debt.UserID, tr.ReceivedCurrency)
		if err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetByUserIdAndCurrency %v", err)
		}

		debt.CompanyID = balance.CompanyId
		debtor.CompanyID = balance.CompanyId
		// debt.UserID = debtor.UserID

		switch debt.Type {
		case TYPE_SELL:
			if balance.Balance >= tr.ReceivedAmount {
				balance.Balance -= tr.ReceivedAmount
				balance.InOutLay += tr.ReceivedAmount

				tr.ReceivedAmount = -tr.ReceivedAmount
				debt.DebtedAmount = -debt.DebtedAmount
				debtor.Balance += debt.DebtedAmount
			} else {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
		case TYPE_BUY:
			balance.Balance += tr.ReceivedAmount
			balance.OutInLay += tr.ReceivedAmount

			debtor.Balance += debt.DebtedAmount
		}

		if err := debtsStorage.Create(ctx, debt); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE debtsStorage.Create %v", err)
		}

		record := &store.BalanceRecord{
			Amount:    tr.ReceivedAmount,
			UserID:    debt.UserID,
			CompanyID: balance.CompanyId,
			BalanceID: balance.ID,
			Type:      int64(debt.Type),
			Details:   debt.Details,
			Currency:  tr.ReceivedCurrency,
			DebtId:    &debt.ID,
		}

		if err := balanceRecordsStorage.Create(ctx, record); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.Create %v", err)
		}

		if err := balancesStorage.Update(ctx, balance); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.Update %v", err)
		}
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.Update %v", err)
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

	old, err := debtsStorage.GetByID(ctx, debt.ID)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtsStorage.GetById %v", err)
	}

	debtor, err := debtorsStorage.GetById(ctx, old.DebtorID)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.GetById %v", err)
	}

	balance, err := balanceStorage.GetByUserIdAndCurrency(ctx, &debtor.UserID, debtor.Currency)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE balanceStorage.GetByUserIdAndCurrency %v", err)
	}

	debt.CompanyID = balance.CompanyId

	for _, tr := range debt.ReceivedIncomes {

		switch old.Type {
		case TYPE_SELL:
			balance.Balance += tr.ReceivedAmount
			balance.InOutLay -= tr.ReceivedAmount

			debtor.Balance -= old.DebtedAmount
		case TYPE_BUY:
			if balance.Balance >= tr.ReceivedAmount {
				balance.Balance -= tr.ReceivedAmount
				balance.OutInLay -= tr.ReceivedAmount

				debtor.Balance -= old.DebtedAmount
			} else {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
		}

		switch debt.Type {
		case TYPE_SELL:
			if balance.Balance >= tr.ReceivedAmount {
				balance.Balance -= tr.ReceivedAmount
				balance.InOutLay += tr.ReceivedAmount

				if tr.ReceivedAmount > 0 {
					tr.ReceivedAmount = -tr.ReceivedAmount
				}
				if debt.DebtedAmount > 0 {
					debt.DebtedAmount = -debt.DebtedAmount
				}

				debtor.Balance += debt.DebtedAmount
			} else {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}

		case TYPE_BUY:
			balance.Balance += tr.ReceivedAmount
			balance.OutInLay += tr.ReceivedAmount

			debtor.Balance += debt.DebtedAmount
		}

		if err := debtsStorage.Update(ctx, debt); err != nil {
			return fmt.Errorf("ERROR OCCURRED debtsStorage.Create( %v", err)
		}

		record := &store.BalanceRecord{
			Amount:    tr.ReceivedAmount,
			UserID:    debt.UserID,
			CompanyID: balance.CompanyId,
			BalanceID: balance.ID,
			Type:      int64(debt.Type),
			Details:   debt.Details,
			Currency:  tr.ReceivedCurrency,
			DebtId:    &debt.ID,
		}

		if err := balanceRecordsStorage.DeleteByDebtId(ctx, old.ID); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE balanceRecordsStorage.DeleteByDebtId %v", err)
		}

		if err := balanceRecordsStorage.Create(ctx, record); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE balanceRecordsStorage.Create %v", err)
		}

		if err := balanceStorage.Update(ctx, balance); err != nil {
			return fmt.Errorf("ERROR OCCURRED WHILE balanceStorage.Update %v", err)
		}
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE debtorsStorage.Update %v", err)
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

	debt, err := debtsStorage.GetByID(ctx, debtId)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.GetById(ctx, *record.DebtorId) %v", err)
	}

	debtor, err := debtorsStorage.GetById(ctx, debt.DebtorID)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.GetById(ctx, *record.DebtorId) %v", err)
	}

	balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &debtor.UserID, debtor.Currency)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetById %v", err)
	}

	for _, tr := range debt.ReceivedIncomes {

		switch debt.Type {
		case TYPE_SELL:
			balance.Balance += tr.ReceivedAmount
			balance.InOutLay -= tr.ReceivedAmount

			debtor.Balance -= debt.DebtedAmount
		case TYPE_BUY:
			if balance.Balance >= tr.ReceivedAmount {
				balance.Balance -= tr.ReceivedAmount
				balance.OutInLay -= tr.ReceivedAmount

				debtor.Balance -= debt.DebtedAmount
			} else {
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
		}

		if err := balancesStorage.Update(ctx, balance); err != nil {
			return err
		}
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		return err
	}

	if err := debtsStorage.Delete(ctx, debtId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE tx.Commit: %w", err)
	}
	return nil
}
