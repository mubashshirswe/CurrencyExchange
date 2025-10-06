package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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
	if err := transactionsStorage.Create(ctx, transaction); err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR OCCURRED WHILE Transactions.Create %v", err)
	}

	for _, tr := range transaction.ReceivedIncomes {
		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &transaction.ReceivedUserId, tr.ReceivedCurrency)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				jsonTr, _ := json.Marshal(transaction)
				log.Println("TRANSACTION:  ")
				log.Println(jsonTr)
				return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.GetByUserIdAndCurrency  %v", err)
			} else {
				return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.GetByUserIdAndCurrency %v", err)
			}
		}

		switch transaction.Type {
		case TYPE_SELL:
			if balance.Balance >= tr.ReceivedAmount {
				balance.Balance -= tr.ReceivedAmount
				balance.InOutLay += tr.ReceivedAmount
			} else {
				tx.Rollback()
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
		case TYPE_BUY:
			balance.Balance += tr.ReceivedAmount
			balance.OutInLay += tr.ReceivedAmount
		default:
			tx.Rollback()
			return fmt.Errorf("FOUND UNKNOWN TYPE")
		}

		balanceRecord := &store.BalanceRecord{
			Amount:        tr.ReceivedAmount,
			Currency:      tr.ReceivedCurrency,
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
		return fmt.Errorf("ERROR OCCURRED WHILE transactionsStorage.GetById %v", err)
	}

	for _, tr := range tran.DeliveredOutcomes {
		balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &transaction.DeliveredUserId, tr.DeliveredCurrency)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ERROR OCCURRED WHILE balancesStorage.GetByUserIdAndCurrency( %v", err)
		}

		var recordType int64
		if tran.Type == TYPE_SELL {
			recordType = TYPE_BUY
			balance.Balance += tr.DeliveredAmount
			balance.OutInLay += tr.DeliveredAmount
		} else {
			recordType = TYPE_SELL
			if balance.Balance >= tr.DeliveredAmount {
				balance.Balance -= tr.DeliveredAmount
				balance.InOutLay += tr.DeliveredAmount
			} else {
				tx.Rollback()
				return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
			}
		}

		balanceRecord := &store.BalanceRecord{
			Amount:        tr.DeliveredAmount,
			Currency:      tr.DeliveredCurrency,
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
	}

	fmt.Println("transaction", tran)
	tran.Status = TRANSACTION_STATUS_COMPLETED
	tran.DeliveredUserId = &transaction.DeliveredUserId

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

	records, err := balanceRecordsStorage.GetByField(ctx, "transaction_id", transaction.ID, types.Pagination{Limit: 100, Offset: 0})
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
		for _, tr := range transaction.ReceivedIncomes {
			balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, &transaction.ReceivedUserId, tr.ReceivedCurrency)
			if err != nil {
				tx.Rollback()
				return err
			}

			record := &store.BalanceRecord{
				Amount:        tr.ReceivedAmount,
				Currency:      tr.ReceivedCurrency,
				UserID:        transaction.ReceivedUserId,
				CompanyID:     balance.CompanyId,
				TransactionId: &transaction.ID,
				BalanceID:     balance.ID,
				Details:       &transaction.Details,
				Type:          transaction.Type,
			}

			if transaction.Type == TYPE_SELL {
				if balance.Balance < tr.ReceivedAmount {
					tx.Rollback()
					return fmt.Errorf(types.BALANCE_NO_ENOUGH_MONEY)
				}
				balance.Balance -= tr.ReceivedAmount
				balance.InOutLay += tr.ReceivedAmount
			} else {
				balance.Balance += tr.ReceivedAmount
				balance.OutInLay += tr.ReceivedAmount
			}

			if err := balancesStorage.Update(ctx, balance); err != nil {
				tx.Rollback()
				return err
			}

			if err := balanceRecordsStorage.Create(ctx, record); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if transaction.DeliveredUserId != nil {
		for _, tr := range transaction.DeliveredOutcomes {
			balance, err := balancesStorage.GetByUserIdAndCurrency(ctx, transaction.DeliveredUserId, tr.DeliveredCurrency)
			if err != nil {
				tx.Rollback()
				return err
			}

			record := &store.BalanceRecord{
				Amount:        tr.DeliveredAmount,
				Currency:      tr.DeliveredCurrency,
				UserID:        *transaction.DeliveredUserId,
				CompanyID:     balance.CompanyId,
				TransactionId: &transaction.ID,
				BalanceID:     balance.ID,
				Details:       &transaction.Details,
				Type:          transaction.Type,
			}

			balance.Balance += tr.DeliveredAmount
			balance.OutInLay += tr.DeliveredAmount

			if err := balancesStorage.Update(ctx, balance); err != nil {
				tx.Rollback()
				return err
			}

			if err := balanceRecordsStorage.Create(ctx, record); err != nil {
				tx.Rollback()
				return err
			}
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

	records, err := balanceRecordsStorage.GetByField(ctx, "transaction_id", tran.ID, types.Pagination{Limit: 100, Offset: 0})
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
		if err := balancesStorage.Update(ctx, balance); err != nil {
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

func (s *TransactionService) GetByCompanyId(ctx context.Context, companyId int64, pagination types.Pagination) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.GetByField(ctx, "delivered_company_id", companyId, pagination)
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
				for _, tr := range tran.DeliveredOutcomes {
					getCurrencies[tr.DeliveredCurrency] += tr.DeliveredAmount
				}
			} else {
				for _, tr := range tran.DeliveredOutcomes {
					giveCurrencies[tr.DeliveredCurrency] += tr.DeliveredAmount
				}
			}

			deliveryUser := ""
			if tran.DeliveredUserId != nil {
				deliveryUser = GetUser(users, *tran.DeliveredUserId).Username
			}
			res := map[string]interface{}{
				"service_fee":        tran.ServiceFee,
				"received_incomes":   tran.ReceivedIncomes,
				"delivered_outcomes": tran.DeliveredOutcomes,
				"received_company":   GetCompany(companies, tran.ReceivedCompanyId).Name,
				"received_user":      GetUser(users, tran.ReceivedUserId).Username,
				"delivered_user":     deliveryUser,
				"phone":              tran.Phone,
				"details":            tran.Details,
				"created_at":         tran.CreatedAt,
				"type":               tran.Type,
				"status":             tran.Status,
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

func (s *TransactionService) GetByField(ctx context.Context, fieldName string, value any, pagination types.Pagination) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.GetByField(ctx, fieldName, value, pagination)
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
		receiverUser := GetUser(users, tran.ReceivedUserId).Username
		deliveryUser := " "
		if tran.DeliveredUserId != nil {
			user := GetUser(users, *tran.DeliveredUserId)
			if user != nil {
				deliveryUser = user.Username
			}
		}

		res := map[string]interface{}{
			"id":                   tran.ID,
			"received_company_id":  tran.ReceivedCompanyId,
			"received_company":     GetCompany(companies, tran.ReceivedCompanyId).Name,
			"received_user_id":     tran.ReceivedUserId,
			"received_user":        receiverUser,
			"received_incomes":     tran.ReceivedIncomes,
			"delivered_outcomes":   tran.DeliveredOutcomes,
			"delivered_company_id": tran.DeliveredCompanyId,
			"delivered_company":    GetCompany(companies, tran.DeliveredCompanyId).Name,
			"delivered_user":       deliveryUser,
			"delivered_user_id":    tran.DeliveredUserId,
			"service_fee":          tran.ServiceFee,
			"phone":                tran.Phone,
			"details":              tran.Details,
			"created_at":           tran.CreatedAt,
			"type":                 tran.Type,
			"status":               tran.Status,
		}

		response = append(response, res)
	}

	return response, nil
}

func (s *TransactionService) Archived(ctx context.Context, pagination types.Pagination) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.Archived(ctx, pagination)
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
		DeliveredUser := ""
		if tran.DeliveredUserId != nil {
			DeliveredUser = GetUser(users, *tran.DeliveredUserId).Username
		}
		res := map[string]interface{}{
			"id":                   tran.ID,
			"received_company_id":  tran.ReceivedCompanyId,
			"received_company":     GetCompany(companies, tran.ReceivedCompanyId).Name,
			"received_user_id":     tran.ReceivedUserId,
			"received_user":        GetUser(users, tran.ReceivedUserId).Username,
			"received_incomes":     tran.ReceivedIncomes,
			"delivered_outcomes":   tran.DeliveredOutcomes,
			"delivered_company_id": tran.DeliveredCompanyId,
			"delivered_company":    GetCompany(companies, tran.DeliveredCompanyId).Name,
			"delivered_user":       DeliveredUser,
			"delivered_user_id":    tran.DeliveredUserId,
			"service_fee":          tran.ServiceFee,
			"phone":                tran.Phone,
			"details":              tran.Details,
			"created_at":           tran.CreatedAt,
			"type":                 tran.Type,
			"status":               tran.Status,
		}

		response = append(response, res)
	}
	return response, nil
}

func (s *TransactionService) GetInfos(ctx context.Context, companyId int64, pagination types.Pagination) ([]map[string]interface{}, error) {
	trans, err := s.store.Transactions.GetByField(ctx, "delivered_company_id", companyId, pagination)
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
				for _, tr := range tran.DeliveredOutcomes {
					getCurrencies[tr.DeliveredCurrency] += tr.DeliveredAmount
				}
			} else {
				for _, tr := range tran.DeliveredOutcomes {
					giveCurrencies[tr.DeliveredCurrency] += tr.DeliveredAmount
				}
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
