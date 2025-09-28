package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

type TransactionPayload struct {
	MarkedServiceFee    string                    `json:"marked_service_fee"`
	DeliveredServiceFee *int64                    `json:"delivered_service_fee"`
	ReceivedIncomes     []types.ReceivedIncomes   `json:"received_incomes"`
	DeliveredOutcomes   []types.DeliveredOutcomes `json:"delivered_outcomes"`
	ReceivedCompanyId   int64                     `json:"received_company_id"`
	DeliveredCompanyId  int64                     `json:"delivered_company_id"`
	ReceivedUserId      int64                     `json:"received_user_id"`
	DeliveredUserId     *int64                    `json:"delivered_user_id"`
	Phone               string                    `json:"phone"`
	Details             string                    `json:"details"`
	Type                int64                     `json:"type"`
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
		MarkedServiceFee:    payload.MarkedServiceFee,
		DeliveredServiceFee: payload.DeliveredServiceFee,
		ReceivedIncomes:     payload.ReceivedIncomes,
		DeliveredOutcomes:   payload.DeliveredOutcomes,
		ReceivedCompanyId:   payload.ReceivedCompanyId,
		DeliveredCompanyId:  payload.DeliveredCompanyId,
		ReceivedUserId:      payload.ReceivedUserId,
		DeliveredUserId:     payload.DeliveredUserId,
		Phone:               payload.Phone,
		Details:             payload.Details,
		Type:                payload.Type,
		Status:              1,
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
		ID:                  getIDFromContext(r),
		MarkedServiceFee:    payload.MarkedServiceFee,
		DeliveredServiceFee: payload.DeliveredServiceFee,
		ReceivedIncomes:     payload.ReceivedIncomes,
		DeliveredOutcomes:   payload.DeliveredOutcomes,
		ReceivedCompanyId:   payload.ReceivedCompanyId,
		DeliveredCompanyId:  payload.DeliveredCompanyId,
		ReceivedUserId:      payload.ReceivedUserId,
		DeliveredUserId:     payload.DeliveredUserId,
		Phone:               payload.Phone,
		Details:             payload.Details,
		Type:                payload.Type,
		Status:              1,
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
	var payload types.TransactionComplete
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.service.Transactions.CompleteTransaction(r.Context(), payload); err != nil {
		if err == sql.ErrNoRows {
			app.badRequestResponse(w, r, fmt.Errorf("BUYURTMA ALLAQACHON YAKUNLANGAN"))
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "SUCCESS"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetTransactionsByFieldHandler(w http.ResponseWriter, r *http.Request) {
	app.LoadPaginationInfo(r, r.Context())
	var payload FieldRequestPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transactions, err := app.service.Transactions.GetByField(r.Context(), payload.FieldName, payload.FieldValue, app.Pagination)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetTransactionsCompanyIdHandler(w http.ResponseWriter, r *http.Request) {
	app.LoadPaginationInfo(r, r.Context())
	transactions, err := app.service.Transactions.GetByCompanyId(r.Context(), getIDFromContext(r), app.Pagination)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetInfosByCompanyIdHandler(w http.ResponseWriter, r *http.Request) {
	app.LoadPaginationInfo(r, r.Context())
	transactions, err := app.service.Transactions.GetInfos(r.Context(), getIDFromContext(r), app.Pagination)
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
	app.LoadPaginationInfo(r, r.Context())
	var payload FieldRequestPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transactions, err := app.store.Transactions.GetByFieldAndDate(r.Context(), payload.FieldName, *payload.From, *payload.To, payload.FieldValue, app.Pagination)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) ArchiveTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserKey).(int64)
	println("userId", userId)

	user, err := app.store.Users.GetById(r.Context(), &userId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.store.Transactions.Archive(r.Context(), user.CompanyId); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "ARCHIVED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) ArchivedTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	app.LoadPaginationInfo(r, r.Context())
	transactions, err := app.service.Transactions.Archived(r.Context(), app.Pagination)
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
