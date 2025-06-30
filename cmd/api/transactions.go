package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type TransactionPayload struct {
	MarkedServiceFee   *int64 `json:"marked_service_fee"`
	ReceivedServiceFee *int64 `json:"received_service_fee"`
	ReceivedAmount     int64  `json:"received_amount"`
	ReceivedCurrency   string `json:"received_currency"`
	DeliveredAmount    int64  `json:"delivered_amount"`
	DeliveredCurrency  string `json:"delivered_currency"`
	SenderCompanyId    int64  `json:"sender_company_id"`
	ReceiverCompanyId  int64  `json:"receiver_company_id"`
	ReceivedUserId     int64  `json:"received_user_id"`
	DeliveredUserId    *int64 `json:"delivered_user_id"`
	Phone              string `json:"phone"`
	Details            string `json:"details"`
	Status             int64  `json:"status"`
	Type               int64  `json:"type"`
}

func (app *application) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var payload TransactionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transaction := &store.Transaction{
		MarkedServiceFee:   payload.MarkedServiceFee,
		ReceivedServiceFee: payload.ReceivedServiceFee,
		ReceivedAmount:     payload.ReceivedAmount,
		ReceivedCurrency:   payload.ReceivedCurrency,
		DeliveredAmount:    payload.DeliveredAmount,
		DeliveredCurrency:  payload.DeliveredCurrency,
		SenderCompanyId:    payload.SenderCompanyId,
		ReceiverCompanyId:  payload.ReceiverCompanyId,
		ReceivedUserId:     payload.ReceivedUserId,
		DeliveredUserId:    payload.DeliveredUserId,
		Phone:              payload.Phone,
		Details:            payload.Details,
		Status:             1,
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

func (app *application) UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var payload TransactionPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transaction := &store.Transaction{
		ID:                 getIDFromContext(r),
		MarkedServiceFee:   payload.MarkedServiceFee,
		ReceivedServiceFee: payload.ReceivedServiceFee,
		ReceivedAmount:     payload.ReceivedAmount,
		ReceivedCurrency:   payload.ReceivedCurrency,
		DeliveredAmount:    payload.DeliveredAmount,
		DeliveredCurrency:  payload.DeliveredCurrency,
		SenderCompanyId:    payload.SenderCompanyId,
		ReceiverCompanyId:  payload.ReceiverCompanyId,
		ReceivedUserId:     payload.ReceivedUserId,
		DeliveredUserId:    payload.DeliveredUserId,
		Phone:              payload.Phone,
		Details:            payload.Details,
		Status:             1,
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

func (app *application) GetTransactionsByFieldHandler(w http.ResponseWriter, r *http.Request) {
	var payload FieldRequestPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transactions, err := app.store.Transactions.GetByField(r.Context(), payload.FieldName, payload.FieldValue)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetTransactionsByFieldAndDateHandler(w http.ResponseWriter, r *http.Request) {
	var payload FieldRequestPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transactions, err := app.store.Transactions.GetByFieldAndDate(r.Context(), payload.FieldName, *payload.From, *payload.To, payload.FieldValue)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
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
