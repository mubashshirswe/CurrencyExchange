package service

import (
	"context"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type DebtorsService struct {
	store store.Storage
}

func (s *DebtorsService) GetByCompanyId(ctx context.Context, companyId int64, pagination types.Pagination) ([]map[string]interface{}, error) {

	debtors, err := s.store.Debtors.GetByCompanyId(ctx, companyId, pagination)
	if err != nil {
		return nil, err
	}

	users, err := s.store.Users.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]map[string]interface{}, 0, len(debtors))

	for _, debtor := range debtors {
		debts, err := s.store.Debts.GetByDebtorId(context.Background(), debtor.ID, pagination)
		if err != nil {
			return nil, err
		}

		var fullname *string
		if len(debts) > 0 {
			fullname = debts[0].Details
		}

		res = append(res, map[string]interface{}{
			"id":         debtor.ID,
			"balance":    debtor.Balance,
			"currency":   debtor.Currency,
			"username":   GetUser(users, debtor.UserID).Username,
			"user_id":    debtor.UserID,
			"company_id": debtor.CompanyID,
			"phone":      debtor.Phone,
			"full_name":  fullname,
			"created_at": debtor.CreatedAt,
		})
	}

	return res, nil
}
