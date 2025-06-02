package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type Debtors struct {
	UserID          int64  `json:"user_id"`
	Amount          int64  `json:"amount"`
	SerialNo        string `json:"serial_no"`
	BalanceId       int64  `json:"balance_id"`
	CompanyId       int64  `json:"company_id"`
	Details         string `json:"details"`
	DebtorsName     string `json:"debtors_name"`
	DebtorsPhone    string `json:"debtors_phone"`
	CurrencyId      int64  `json:"currency_id"`
	CurrencyType    string `json:"currency_type"`
	Type            int    `json:"type"`
	IsBalanceEffect int    `json:"is_balance_effect"`
}

func (app *application) CreateDebtorsHandler(w http.ResponseWriter, r *http.Request) {
	var payload Debtors
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	debtor := &store.Debtors{
		Type:            payload.Type,
		CurrencyType:    payload.CurrencyType,
		UserID:          payload.UserID,
		Amount:          payload.Amount,
		BalanceId:       payload.BalanceId,
		CompanyId:       payload.CompanyId,
		Details:         payload.Details,
		DebtorsName:     payload.DebtorsName,
		DebtorsPhone:    payload.DebtorsPhone,
		CurrencyId:      payload.CurrencyId,
		IsBalanceEffect: payload.IsBalanceEffect,
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
		ID:              getIDFromContext(r),
		Type:            payload.Type,
		CurrencyType:    payload.CurrencyType,
		UserID:          payload.UserID,
		SerialNo:        payload.SerialNo,
		Amount:          payload.Amount,
		BalanceId:       payload.BalanceId,
		CompanyId:       payload.CompanyId,
		Details:         payload.Details,
		DebtorsName:     payload.DebtorsName,
		DebtorsPhone:    payload.DebtorsPhone,
		CurrencyId:      payload.CurrencyId,
		IsBalanceEffect: payload.IsBalanceEffect,
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

func (app *application) ReceivedDebtHandler(w http.ResponseWriter, r *http.Request) {
	err := app.service.Debtors.ReceivedDebt(r.Context(), getIDFromContext(r))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "SUCCESS"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetDebtorsByIdHandler(w http.ResponseWriter, r *http.Request) {
	debtors, err := app.store.Debtors.GetById(r.Context(), getIDFromContext(r))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, debtors); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteDebtorsHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	if err := app.service.Debtors.Delete(r.Context(), id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
