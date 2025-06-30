package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
	"github.com/mubashshir3767/currencyExchange/internal/types"
)

/*
   ReceivedMoney + currency
   SelledMondey  + currency
   UserID
   CompanyID
   Details
*/

type FieldRequestPayload struct {
	From       *string `json:"from"`
	To         *string `json:"to"`
	FieldName  string  `json:"field_name"`
	FieldValue any     `json:"field_value"`
}

func (app *application) CreateBalanceRecordHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.BalanceRecordPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.service.BalanceRecords.PerformBalanceRecord(r.Context(), payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetBalanceRecordsByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	var payload FieldRequestPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	records, err := app.store.BalanceRecords.GetByField(r.Context(), payload.FieldName, payload.FieldValue)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, records); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetBalanceRecordsByBalanceIdHandler(w http.ResponseWriter, r *http.Request) {
	var payload FieldRequestPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	records, err := app.store.BalanceRecords.GetByFieldAndDate(r.Context(), payload.FieldName, *payload.From, *payload.To, payload.FieldValue)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, records); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateBalanceRecordHandler(w http.ResponseWriter, r *http.Request) {
	var payload store.BalanceRecord
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.service.BalanceRecords.UpdateRecord(r.Context(), payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteBalanceRecordHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	if err := app.service.BalanceRecords.RollbackBalanceRecord(r.Context(), id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
