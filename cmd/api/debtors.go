package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type DebtorPayload struct {
	FullName         string  `json:"full_name"`
	ReceivedAmount   int64   `json:"received_amount"`
	ReceivedCurrency string  `json:"received_currency"`
	DebtedAmount     int64   `json:"debted_amount"`
	DebtedCurrency   string  `json:"debted_currency"`
	UserID           int64   `json:"user_id"`
	Details          *string `json:"details"`
	Phone            *string `json:"phone"`
	IsBalanceEffect  int     `json:"is_balance_effect"`
	Type             int     `json:"type"`
}

func (app *application) CreateDebtorsHandler(w http.ResponseWriter, r *http.Request) {
	var payload DebtorPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	debtor := &store.Debts{
		FullName:         payload.FullName,
		ReceivedAmount:   payload.ReceivedAmount,
		ReceivedCurrency: payload.ReceivedCurrency,
		DebtedAmount:     payload.DebtedAmount,
		DebtedCurrency:   payload.DebtedCurrency,
		UserID:           payload.UserID,
		Details:          payload.Details,
		Phone:            payload.Phone,
		IsBalanceEffect:  payload.IsBalanceEffect,
		Type:             payload.Type,
	}

	if err := app.service.Debts.Create(r.Context(), debtor); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, debtor); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) CreateDebtorTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var payload *store.Debts
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.service.Debts.Transaction(r.Context(), payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateDebtsHandler(w http.ResponseWriter, r *http.Request) {
	var payload *store.Debts
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	debt := &store.Debts{
		ID:               getIDFromContext(r),
		FullName:         payload.FullName,
		ReceivedAmount:   payload.ReceivedAmount,
		ReceivedCurrency: payload.ReceivedCurrency,
		DebtedAmount:     payload.DebtedAmount,
		DebtedCurrency:   payload.DebtedCurrency,
		UserID:           payload.UserID,
		Details:          payload.Details,
		Phone:            payload.Phone,
		IsBalanceEffect:  payload.IsBalanceEffect,
		Type:             payload.Type,
		DebtorId:         payload.DebtorId,
	}

	if err := app.service.Debts.Update(r.Context(), debt); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetDebtorsByCompanyIdHandler(w http.ResponseWriter, r *http.Request) {
	debtors, err := app.store.Debtors.GetByCompanyId(r.Context(), getIDFromContext(r))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, debtors); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetDebtsByDebtorIdHandler(w http.ResponseWriter, r *http.Request) {
	debtors, err := app.store.Debts.GetByDebtorId(r.Context(), getIDFromContext(r))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, debtors); err != nil {
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

func (app *application) DeleteDebtsHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	if err := app.service.Debts.Delete(r.Context(), id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteDebtorsHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	if err := app.store.Debtors.Delete(r.Context(), id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
