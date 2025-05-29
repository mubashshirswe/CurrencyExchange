package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type Debtors struct {
	UserID       int64  `json:"user_id"`
	Amount       int64  `json:"amount"`
	BalanceId    int64  `json:"balance_id"`
	CompanyId    int64  `json:"company_id"`
	Details      string `json:"details"`
	DebtorsName  string `json:"debtors_name"`
	DebtorsPhone string `json:"debtors_phone"`
	CurrencyId   int64  `json:"currency_id"`
	Type         int    `json:"type"`
}

func (app *application) CreateDebtorsHandler(w http.ResponseWriter, r *http.Request) {
	var payload Debtors
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	debtor := &store.Debtors{
		UserID:       payload.UserID,
		Amount:       payload.Amount,
		BalanceId:    payload.BalanceId,
		CompanyId:    payload.CompanyId,
		Details:      payload.Details,
		DebtorsName:  payload.DebtorsName,
		DebtorsPhone: payload.DebtorsPhone,
		CurrencyId:   payload.CurrencyId,
	}

	if err := app.service.Debtors.Create(r.Context(), debtor); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, debtor); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateDebtorsHandler(w http.ResponseWriter, r *http.Request) {
	var payload Debtors
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	debtor := &store.Debtors{
		ID:           getIDFromContext(r),
		UserID:       payload.UserID,
		Amount:       payload.Amount,
		BalanceId:    payload.BalanceId,
		CompanyId:    payload.CompanyId,
		Details:      payload.Details,
		DebtorsName:  payload.DebtorsName,
		DebtorsPhone: payload.DebtorsPhone,
		CurrencyId:   payload.CurrencyId,
	}

	if err := app.service.Debtors.Update(r.Context(), debtor); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, debtor); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetDebtorsByUserIdHandler(w http.ResponseWriter, r *http.Request) {

	debtors, err := app.service.Debtors.GetByUserId(r.Context(), getIDFromContext(r))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, debtors); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	if err := app.service.Debtors.Delete(r.Context(), getIDFromContext(r)); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
