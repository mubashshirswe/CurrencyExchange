package service

import (
	"context"
	"fmt"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type DebtorsService struct {
	store store.Storage
}

func (s *DebtorsService) Create(ctx context.Context, debtor *store.Debtors) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)

	balance, err := balancesStorage.GetByIdAndCurrency(ctx, &debtor.UserID, debtor.ReceivedCurrency)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetByIdAndCurrency")
	}

	debtor.CompanyID = balance.CompanyId

	switch debtor.Type {
	case TYPE_SELL:
		if balance.Balance >= debtor.ReceivedAmount {
			balance.Balance -= debtor.ReceivedAmount
			balance.InOutLay += debtor.ReceivedAmount
		} else {
			return fmt.Errorf("ERROR OCCURRED: BALANCE HAS NO ENOUGH MONEY TO OPERATE %v >= %v ", balance.Balance, debtor.ReceivedAmount)
		}
	case TYPE_BUY:
		balance.Balance += debtor.ReceivedAmount
		balance.OutInLay += debtor.ReceivedAmount
	}

	if err := debtorsStorage.Create(ctx, debtor); err != nil {
		tx.Rollback()
		return err
	}

	record := &store.BalanceRecord{
		Amount:    debtor.ReceivedAmount,
		UserID:    debtor.UserID,
		CompanyID: balance.CompanyId,
		BalanceID: balance.ID,
		Type:      int64(debtor.Type),
		Details:   debtor.Details,
		Currency:  debtor.ReceivedCurrency,
	}

	if err := balanceRecordsStorage.Create(ctx, record); err != nil {
		tx.Rollback()
		return fmt.Errorf("ROLLBACK ISHLADI")
	}

	if err := balancesStorage.Update(ctx, balance); err == nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *DebtorsService) Transaction(ctx context.Context, debtor *store.Debtors) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	debtorsStorage := store.NewDebtorsStorage(tx)
	balancesStorage := store.NewBalanceStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)

	oldDebtor, err := debtorsStorage.GetById(ctx, debtor.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if oldDebtor.DebtedCurrency != debtor.DebtedCurrency {
		tx.Rollback()
		return fmt.Errorf("DEBTED CURRENCIES ARE NOT MATCH %v", err)
	}

	balance, err := balancesStorage.GetByIdAndCurrency(ctx, &debtor.UserID, debtor.ReceivedCurrency)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetByIdAndCurrency %v", err)
	}

	debtor.CompanyID = balance.CompanyId

	switch debtor.Type {
	case TYPE_SELL:
		if balance.Balance >= debtor.ReceivedAmount {
			balance.Balance -= debtor.ReceivedAmount
			balance.InOutLay += debtor.ReceivedAmount

			oldDebtor.DebtedAmount -= debtor.DebtedAmount
		} else {
			tx.Rollback()
			return fmt.Errorf("ERROR OCCURRED: BALANCE HAS NO ENOUGH MONEY TO OPERATE %v >= %v ", balance.Balance, debtor.ReceivedAmount)
		}
	case TYPE_BUY:
		balance.Balance += debtor.ReceivedAmount
		balance.OutInLay += debtor.ReceivedAmount

		oldDebtor.DebtedAmount += debtor.DebtedAmount
	}

	record := &store.BalanceRecord{
		Amount:    debtor.ReceivedAmount,
		UserID:    debtor.UserID,
		CompanyID: balance.CompanyId,
		BalanceID: balance.ID,
		Type:      int64(debtor.Type),
		Details:   debtor.Details,
		Currency:  debtor.ReceivedCurrency,
	}

	if err := balanceRecordsStorage.Create(ctx, record); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.Create %v", err)
	}

	if err := debtorsStorage.Update(ctx, oldDebtor); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.Update %v", err)
	}

	if err := balancesStorage.Update(ctx, balance); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.Update %v", err)
	}

	tx.Commit()
	return nil
}

func (s *DebtorsService) Update(ctx context.Context, record *store.BalanceRecord) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	debtorsStorage := store.NewDebtorsStorage(tx)
	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	balanceStorage := store.NewBalanceStorage(tx)

	debtor, err := debtorsStorage.GetById(ctx, *record.DebtorId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.GetById(ctx, *record.DebtorId) %v", err)
	}

	Oldrecords, err := balanceRecordsStorage.GetByField(ctx, "id", &record.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	oldRecord := Oldrecords[0]

	balance, err := balanceStorage.GetById(ctx, &oldRecord.BalanceID)
	if err != nil {
		tx.Rollback()
		return err
	}

	debtor.CompanyID = balance.CompanyId

	switch oldRecord.Type {
	case TYPE_SELL:
		balance.Balance += oldRecord.Amount
		balance.InOutLay -= oldRecord.Amount

		debtor.DebtedAmount += oldRecord.Amount

	case TYPE_BUY:
		if balance.Balance >= oldRecord.Amount {
			balance.Balance -= oldRecord.Amount
			balance.OutInLay -= oldRecord.Amount

			debtor.DebtedAmount -= oldRecord.Amount
		} else {
			tx.Rollback()
			return fmt.Errorf("BALANCE THAT IS ROLLBACKING HAS NO ENOUGH MONEY")
		}
	}

	switch record.Type {
	case TYPE_SELL:
		balance.Balance -= record.Amount
		balance.InOutLay += record.Amount

		debtor.DebtedAmount += record.Amount
	case TYPE_BUY:
		balance.Balance += record.Amount
		balance.OutInLay += record.Amount

		debtor.DebtedAmount -= record.Amount
	}

	if err := balanceStorage.Update(ctx, balance); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING BALANCE %v", err)
	}

	if err := balanceRecordsStorage.Update(ctx, record); err != nil {
		tx.Rollback()
		return err
	}

	if err := debtorsStorage.Update(ctx, debtor); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *DebtorsService) Delete(ctx context.Context, balanceRecordId int64) error {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	balanceRecordsStorage := store.NewBalanceRecordStorage(tx)
	balancesStorage := store.NewBalanceStorage(tx)
	debtorsStorage := store.NewDebtorsStorage(tx)

	records, err := balanceRecordsStorage.GetByField(ctx, "id", balanceRecordId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.GetByField %v", err)
	}

	record := records[0]
	balance, err := balancesStorage.GetById(ctx, &record.BalanceID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetById %v", err)
	}

	debtor, err := debtorsStorage.GetById(ctx, *record.DebtorId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.GetById(ctx, *record.DebtorId) %v", err)
	}

	switch record.Type {
	case TYPE_SELL:
		balance.Balance += record.Amount
		balance.InOutLay -= record.Amount

		debtor.DebtedAmount -= record.Amount
	case TYPE_BUY:
		if balance.Balance >= record.Amount {
			balance.Balance -= record.Amount
			balance.OutInLay -= record.Amount

			debtor.DebtedAmount += record.Amount
		} else {
			tx.Rollback()
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE %v >= %v", balance.Balance, record.Amount)
		}
	}

	if err := balanceRecordsStorage.Delete(ctx, balanceRecordId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (s *DebtorsService) GetByCompanyId(ctx context.Context, companyId int64) (map[string]interface{}, error) {
	s.store.Debtors.GetByUserId(ctx, companyId)

	return nil, nil
}
