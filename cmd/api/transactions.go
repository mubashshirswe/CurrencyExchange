package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CreateTransactionPayload struct {
	Amount             int64  `json:"amount"`
	ServiceFee         int64  `json:"service_fee"`
	FromCurrencyTypeId int64  `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64  `json:"to_currency_type_id"`
	ReceiverId         int64  `json:"receiver_id"`
	SenderId           int64  `json:"sender_id"`
	FromCityId         int64  `json:"from_city_id"`
	ToCityId           int64  `json:"to_city_id"`
	ReceiverName       string `json:"receiver_name"`
	ReceiverPhone      string `json:"receiver_phone"`
	Details            string `json:"details"`
	SerialNo           string `json:"serial_no"`
	CompanyId          int64  `json:"company_id"`
	BalanceId          int64  `json:"balance_id"`
}

type UpdateTransactionPayload struct {
	Amount             int64  `json:"amount"`
	ServiceFee         int64  `json:"service_fee"`
	FromCurrencyTypeId int64  `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64  `json:"to_currency_type_id"`
	ReceiverId         int64  `json:"receiver_id"`
	FromCityId         int64  `json:"from_city_id"`
	ToCityId           int64  `json:"to_city_id"`
	ReceiverName       string `json:"receiver_name"`
	ReceiverPhone      string `json:"receiver_phone"`
	Details            string `json:"details"`
	SerialNo           string `json:"serial_no"`
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
		ReceiverId:         payload.ReceiverId,
		FromCityId:         payload.FromCityId,
		ToCityId:           payload.ToCityId,
		ReceiverName:       payload.ReceiverName,
		ReceiverPhone:      payload.ReceiverPhone,
		Details:            payload.Details,
		CompanyId:          payload.CompanyId,
		BalanceId:          payload.BalanceId,
	}

	if err := app.service.Transactions.PerformTransaction(r.Context(), transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) CompleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	serialNo := chi.URLParam(r, "serial_no")

	if err := app.service.Transactions.CompleteTransaction(r.Context(), serialNo); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, serialNo); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllTransactionByBalanceIdHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) GetAllTransactionByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	transactions, err := app.store.Transactions.GetAllByUserId(r.Context(), &id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllTransactionByReceiverIdHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	transactions, err := app.store.Transactions.GetAllByReceiverId(r.Context(), &id)
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

func (app *application) GetAllActiveTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")
	var value int64
	if status == "1" {
		value = 1
	} else if status == "2" {
		value = 2
	} else {
		app.badRequestResponse(w, r, fmt.Errorf("STATUS NOT GIVEN"))
	}

	transactions, err := app.store.Transactions.GetAllByStatus(r.Context(), value)
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
	senderID, _ := r.Context().Value(UserKey).(int)

	transaction := &store.Transaction{
		SerialNo:           payload.SerialNo,
		Amount:             payload.Amount,
		ServiceFee:         payload.ServiceFee,
		FromCurrencyTypeId: payload.FromCurrencyTypeId,
		ToCurrencyTypeId:   payload.ToCurrencyTypeId,
		SenderId:           int64(senderID),
		ReceiverId:         payload.ReceiverId,
		FromCityId:         payload.FromCityId,
		ToCityId:           payload.ToCityId,
		ReceiverName:       payload.ReceiverName,
		ReceiverPhone:      payload.ReceiverPhone,
		Details:            payload.Details,
	}

	if err := app.service.Transactions.Update(r.Context(), transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transaction); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	if err := app.service.Transactions.Delete(r.Context(), &id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
