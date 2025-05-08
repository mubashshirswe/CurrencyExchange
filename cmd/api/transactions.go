package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CreateTransactionPayload struct {
	Amount             int64  `json:"amount"`
	ServiceFee         int64  `json:"service_fee"`
	FromCurrencyTypeId int64  `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64  `json:"to_currency_type_id"`
	SenderId           int64  `json:"sender_id"`
	FromCityId         int64  `json:"from_city_id"`
	ToCityId           int64  `json:"to_city_id"`
	ReceiverName       string `json:"receiver_name"`
	ReceiverPhone      string `json:"receiver_phone"`
	Details            string `json:"details"`
	Type               int64  `json:"type"`
	CompanyId          int64  `json:"company_id"`
	BalanceId          int64  `json:"balance_id"`
}

type UpdateTransactionPayload struct {
	Amount             int64  `json:"amount"`
	ServiceFee         int64  `json:"service_fee"`
	FromCurrencyTypeId int64  `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64  `json:"to_currency_type_id"`
	SenderId           int64  `json:"sender_id"`
	FromCityId         int64  `json:"from_city_id"`
	ToCityId           int64  `json:"to_city_id"`
	ReceiverName       string `json:"receiver_name"`
	ReceiverPhone      string `json:"receiver_phone"`
	Details            string `json:"details"`
	Type               int64  `json:"type"`
}

type DateTransactionPayload struct {
	From      string `json:"from"`
	To        string `json:"to"`
	BalanceId int64  `json:"balance_id"`
}

func (app *application) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateTransactionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transaction := &store.Transaction{
		Amount:             payload.Amount,
		ServiceFee:         payload.ServiceFee,
		FromCurrencyTypeId: payload.FromCurrencyTypeId,
		ToCurrencyTypeId:   payload.ToCurrencyTypeId,
		SenderId:           payload.SenderId,
		FromCityId:         payload.FromCityId,
		ToCityId:           payload.ToCityId,
		ReceiverName:       payload.ReceiverName,
		ReceiverPhone:      payload.ReceiverPhone,
		Details:            payload.Details,
		Type:               payload.Type,
		CompanyId:          payload.CompanyId,
		BalanceId:          payload.BalanceId,
	}

	if err := app.store.Transactions.Create(r.Context(), transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	transactions, err := app.store.Transactions.GetAllByBalanceId(r.Context(), &id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllTransactionByDateHandler(w http.ResponseWriter, r *http.Request) {
	var payload DateTransactionPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transactions, err := app.store.Transactions.GetAllByDate(r.Context(), payload.From, payload.To, &payload.BalanceId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetTransactionByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	transaction, err := app.store.Transactions.GetById(r.Context(), &id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateTransactionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transaction := &store.Transaction{
		Amount:             payload.Amount,
		ServiceFee:         payload.ServiceFee,
		FromCurrencyTypeId: payload.FromCurrencyTypeId,
		ToCurrencyTypeId:   payload.ToCurrencyTypeId,
		SenderId:           payload.SenderId,
		FromCityId:         payload.FromCityId,
		ToCityId:           payload.ToCityId,
		ReceiverName:       payload.ReceiverName,
		ReceiverPhone:      payload.ReceiverPhone,
		Details:            payload.Details,
		Type:               payload.Type,
	}

	if err := app.store.Transactions.Update(r.Context(), transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
