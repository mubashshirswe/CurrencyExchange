package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CreateBalancePayload struct {
	Balance   int64 `json:"balance"`
	UserId    int64 `json:"user_id"`
	InOutLay  int64 `json:"in_out_lay"`
	OutInLay  int64 `json:"out_in_lay"`
	CompanyId int64 `json:"company_id"`
}

type UpdateBalancePayload struct {
	Balance  int64 `json:"balance"`
	UserId   int64 `json:"user_id"`
	InOutLay int64 `json:"in_out_lay"`
	OutInLay int64 `json:"out_in_lay"`
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
		Balance:   payload.Balance,
		UserId:    payload.UserId,
		InOutLay:  payload.InOutLay,
		OutInLay:  payload.OutInLay,
		CompanyId: payload.CompanyId,
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
	id := GetIdFromContext(r)
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
	id := GetIdFromContext(r)
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
