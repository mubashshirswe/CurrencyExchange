package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CreateBalancePayload struct {
	ID         int64 `json:"id"`
	Balance    int64 `json:"balance"`
	UserId     int64 `json:"user_id"`
	InOutLay   int64 `json:"in_out_lay"`
	CurrencyId int64 `json:"currency_id"`
	OutInLay   int64 `json:"out_in_lay"`
	CompanyId  int64 `json:"company_id"`
}

type UpdateBalancePayload struct {
	ID         int64 `json:"id"`
	Balance    int64 `json:"balance"`
	UserId     int64 `json:"user_id"`
	InOutLay   int64 `json:"in_out_lay"`
	OutInLay   int64 `json:"out_in_lay"`
	CurrencyId int64 `json:"currency_id"`
}

func (app *application) CreateBalanceHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateBalancePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	balance := &store.Balance{
		Balance:    payload.Balance,
		UserId:     payload.UserId,
		InOutLay:   payload.InOutLay,
		OutInLay:   payload.OutInLay,
		CurrencyId: payload.CurrencyId,
		CompanyId:  payload.CompanyId,
	}

	if err := app.store.Balances.Create(r.Context(), balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetBalanceByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	balance, err := app.store.Balances.GetById(r.Context(), &id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetBalanceByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	balance, err := app.store.Balances.GetByUserId(r.Context(), &id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllBalanceHandler(w http.ResponseWriter, r *http.Request) {
	balance, err := app.store.Balances.GetAll(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateBalanceHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateBalancePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	balance := &store.Balance{
		Balance:  payload.Balance,
		UserId:   payload.UserId,
		InOutLay: payload.InOutLay,
		OutInLay: payload.OutInLay,
	}

	if err := app.store.Balances.Update(r.Context(), balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, balance); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteBalanceHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	if err := app.store.Balances.Delete(r.Context(), id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
