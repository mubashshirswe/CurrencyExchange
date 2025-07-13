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
			return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
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
			return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
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

	tran.Status = TRANSACTION_STATUS_COMPLETED
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
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
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
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
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

func (s *TransactionService) GetByCompanyId(ctx context.Context, companyId int64) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.GetByField(ctx, "delivered_company_id", companyId)
	if err != nil {
		return nil, err
	}

	companies, err := s.store.Companies.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	users, err := s.store.Users.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	getCurrencies := make(map[string]int64)
	giveCurrencies := make(map[string]int64)

	for _, tran := range trans {
		if tran.Status == TRANSACTION_STATUS_PENDING {
			if tran.Type == TYPE_SELL {
				getCurrencies[tran.DeliveredCurrency] += tran.DeliveredAmount
			} else {
				giveCurrencies[tran.DeliveredCurrency] += tran.DeliveredAmount
			}
			res := map[string]interface{}{
				"marked_service_fee":    tran.MarkedServiceFee,
				"received_amount":       tran.ReceivedAmount,
				"received_company":      GetCompany(companies, tran.ReceivedCompanyId).Name,
				"received_user":         GetUser(users, &tran.ReceivedUserId),
				"received_currency":     tran.ReceivedCurrency,
				"delivered_currency":    tran.DeliveredCurrency,
				"delivered_amount":      tran.DeliveredAmount,
				"delivered_user":        GetUser(users, tran.DeliveredUserId).Username,
				"delivered_service_fee": tran.DeliveredServiceFee,
				"phone":                 tran.Phone,
				"details":               tran.Details,
				"created_at":            tran.CreatedAt,
				"type":                  tran.Type,
			}

			response = append(response, res)
		}
	}
	response = append(response, map[string]interface{}{
		"get_currencies": getCurrencies,
		"giveCurrencies": giveCurrencies,
	})

	return response, nil
}

func (s *TransactionService) GetByField(ctx context.Context, fieldName string, value any) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.GetByField(ctx, fieldName, value)
	if err != nil {
		return nil, err
	}

	companies, err := s.store.Companies.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	users, err := s.store.Users.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, tran := range trans {
		res := map[string]interface{}{
			"marked_service_fee":    tran.MarkedServiceFee,
			"received_amount":       tran.ReceivedAmount,
			"received_company":      GetCompany(companies, tran.ReceivedCompanyId).Name,
			"received_company_id":   tran.ReceivedCompanyId,
			"received_user":         GetUser(users, &tran.ReceivedUserId).Username,
			"received_user_id":      tran.ReceivedUserId,
			"received_currency":     tran.ReceivedCurrency,
			"delivered_currency":    tran.DeliveredCurrency,
			"delivered_amount":      tran.DeliveredAmount,
			"delivered_company_id":  GetCompany(companies, tran.DeliveredCompanyId).Name,
			"delivered_user":        GetUser(users, tran.DeliveredUserId),
			"delivered_user_id":     tran.DeliveredUserId,
			"delivered_service_fee": tran.DeliveredServiceFee,
			"phone":                 tran.Phone,
			"details":               tran.Details,
			"created_at":            tran.CreatedAt,
			"type":                  tran.Type,
		}

		response = append(response, res)
	}

	return response, nil
}

func (s *TransactionService) Archived(ctx context.Context) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.Archived(ctx)
	if err != nil {
		return nil, err
	}

	companies, err := s.store.Companies.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	users, err := s.store.Users.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, tran := range trans {
		res := map[string]interface{}{
			"marked_service_fee":    tran.MarkedServiceFee,
			"received_amount":       tran.ReceivedAmount,
			"received_company":      GetCompany(companies, tran.ReceivedCompanyId).Name,
			"received_company_id":   tran.ReceivedCompanyId,
			"received_user":         GetUser(users, &tran.ReceivedUserId).Username,
			"received_user_id":      tran.ReceivedUserId,
			"received_currency":     tran.ReceivedCurrency,
			"delivered_currency":    tran.DeliveredCurrency,
			"delivered_amount":      tran.DeliveredAmount,
			"delivered_company_id":  GetCompany(companies, tran.DeliveredCompanyId).Name,
			"delivered_user":        GetUser(users, tran.DeliveredUserId).Username,
			"delivered_user_id":     tran.DeliveredUserId,
			"delivered_service_fee": tran.DeliveredServiceFee,
			"phone":                 tran.Phone,
			"details":               tran.Details,
			"created_at":            tran.CreatedAt,
			"type":                  tran.Type,
		}

		response = append(response, res)
	}
	return response, nil
}

func (s *TransactionService) GetInfos(ctx context.Context, companyId int64) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.GetByField(ctx, "delivered_company_id", companyId)
	if err != nil {
		return nil, fmt.Errorf("ERROR OCCURRED WHILE Transactions.GetByField %v", err)
	}

	balances, err := s.store.Balances.GetByCompanyId(ctx, &companyId)
	if err != nil {
		return nil, fmt.Errorf("ERROR OCCURRED WHILE Balances.GetByCompanyId %v", err)
	}

	var response []map[string]interface{}
	getCurrencies := make(map[string]int64)
	giveCurrencies := make(map[string]int64)
	currencies := make(map[string]int64)
	free_currencies := make(map[string]int64)

	for _, balance := range balances {
		currencies[balance.Currency] += balance.Balance
	}

	for _, tran := range trans {
		if tran.Status == TRANSACTION_STATUS_PENDING {
			if tran.Type == TYPE_SELL {
				getCurrencies[tran.DeliveredCurrency] += tran.DeliveredAmount
			} else {
				giveCurrencies[tran.DeliveredCurrency] += tran.DeliveredAmount
			}
		}
	}

	for i, cur := range currencies {
		getOne := GetOne(getCurrencies, i)
		giveOne := GetOne(giveCurrencies, i)
		cur += getOne
		cur -= giveOne

		free_currencies[i] += cur
	}

	response = append(response, map[string]interface{}{
		"balances":       currencies,
		"get_currencies": getCurrencies,
		"giveCurrencies": giveCurrencies,
		"free_balances":  free_currencies,
	})

	return response, nil
}

func GetOne(ids map[string]int64, id string) int64 {
	for i, company := range ids {
		if i == id {
			return company
		}
	}
	return 0
}

func GetCompany(companies []store.Company, companyId int64) *store.Company {
	for _, company := range companies {
		if company.ID == companyId {
			return &company
		}
	}
	return nil
}
