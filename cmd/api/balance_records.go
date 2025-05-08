package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CreateBalanceRecordPayload struct {
	Amount     int64  `json:"amount"`
	UserID     int64  `json:"user_id"`
	BalanceID  int64  `json:"balance_id"`
	CompanyID  int64  `json:"company_id"`
	Details    string `json:"details"`
	CurrenctID int64  `json:"currency_id"`
	Type       int64  `json:"type"`
}

type UpdateBalanceRecordPayload struct {
	Amount     int64  `json:"amount"`
	Details    string `json:"details"`
	CurrenctID int64  `json:"currency_id"`
	Type       int64  `json:"type"`
}

func (app *application) CreateBalanceRecordHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateBalanceRecordPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	balanceRecord := &store.BalanceRecord{
		Amount:     payload.Amount,
		UserID:     payload.UserID,
		BalanceID:  payload.BalanceID,
		CompanyID:  payload.CompanyID,
		Details:    payload.Details,
		CurrenctID: payload.CurrenctID,
		Type:       payload.Type,
	}

	if err := app.service.BalanceRecords.PerformBalanceRecord(r.Context(), balanceRecord); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, balanceRecord); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetBalanceRecordsByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	userID := getIDFromContext(r)

	records, err := app.store.BalanceRecords.GetByUserId(r.Context(), userID)
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
	balanceID := getIDFromContext(r)

	records, err := app.store.BalanceRecords.GetByBalanceId(r.Context(), balanceID)
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
	var payload UpdateBalanceRecordPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	balanceRecord := &store.BalanceRecord{
		Amount:     payload.Amount,
		Details:    payload.Details,
		CurrenctID: payload.CurrenctID,
		Type:       payload.Type,
	}

	if err := app.store.BalanceRecords.Update(r.Context(), balanceRecord); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, balanceRecord); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteBalanceRecordHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	if err := app.store.BalanceRecords.Delete(r.Context(), id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
