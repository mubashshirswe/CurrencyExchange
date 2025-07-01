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
	// tx,err := s.store.B

	balance, err := s.store.Balances.GetByIdAndCurrency(ctx, &debtor.UserID, debtor.ReceivedCurrency)
	if err != nil {
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

	if err := s.store.Debtors.Create(ctx, debtor); err != nil {
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

	if err := s.store.BalanceRecords.Create(ctx, record); err != nil {
		return err
	}

	if err := s.store.Balances.Update(ctx, balance); err != nil {
		return err
	}

	return nil
}

func (s *DebtorsService) Transaction(ctx context.Context, debtor *store.Debtors) error {
	oldDebtor, err := s.store.Debtors.GetById(ctx, debtor.ID)
	if err != nil {
		return err
	}

	if oldDebtor.DebtedCurrency != debtor.DebtedCurrency {
		return fmt.Errorf("DEBTED CURRENCIES ARE NOT MATCH %v", err)
	}

	balance, err := s.store.Balances.GetByIdAndCurrency(ctx, &debtor.UserID, debtor.ReceivedCurrency)
	if err != nil {
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

	if err := s.store.BalanceRecords.Create(ctx, record); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.Create %v", err)
	}

	if err := s.store.Debtors.Update(ctx, oldDebtor); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.Update %v", err)
	}

	if err := s.store.Balances.Update(ctx, balance); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.Update %v", err)
	}

	return nil
}

func (s *DebtorsService) Update(ctx context.Context, record *store.BalanceRecord) error {
	debtor, err := s.store.Debtors.GetById(ctx, *record.DebtorId)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Debtors.GetById(ctx, *record.DebtorId) %v", err)
	}
	Oldrecords, err := s.store.BalanceRecords.GetByField(ctx, "id", &record.ID)
	if err != nil {
		return err
	}
	oldRecord := Oldrecords[0]

	balance, err := s.store.Balances.GetById(ctx, &oldRecord.BalanceID)
	if err != nil {
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

	if err := s.store.Balances.Update(ctx, balance); err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE UPDATING BALANCE %v", err)
	}

	if err := s.store.BalanceRecords.Update(ctx, record); err != nil {
		return err
	}

	return s.store.Debtors.Update(ctx, debtor)
}

func (s *DebtorsService) Delete(ctx context.Context, balanceRecordId int64) error {
	records, err := s.store.BalanceRecords.GetByField(ctx, "id", balanceRecordId)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE BalanceRecords.GetByField %v", err)
	}
	record := records[0]
	balance, err := s.store.Balances.GetById(ctx, &record.BalanceID)
	if err != nil {
		return fmt.Errorf("ERROR OCCURRED WHILE Balances.GetById %v", err)
	}

	debtor, err := s.store.Debtors.GetById(ctx, *record.DebtorId)
	if err != nil {
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
			return fmt.Errorf("BALANCE HAS NO ENOUGH MONEY TO OPERATE %v >= %v", balance.Balance, record.Amount)
		}
	}

	return s.store.BalanceRecords.Delete(ctx, balanceRecordId)
}

func (s *DebtorsService) GetByCompanyId(ctx context.Context, companyId int64) (map[string]interface{}, error) {
	s.store.Debtors.GetByUserId(ctx, companyId)

	return nil, nil
}
