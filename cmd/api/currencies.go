package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CreateCurrencyPayload struct {
	Name      string `json:"name"`
	Sell      *int64 `json:"sell"`
	Buy       *int64 `json:"buy"`
	CompanyId int64  `json:"company_id"`
}

type UpdateCurrencyPayload struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Sell      *int64 `json:"sell"`
	Buy       *int64 `json:"buy"`
	CompanyId int64  `json:"company_id"`
}

func (app *application) CreateCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCurrencyPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	currency := &store.Currency{
		Name:      payload.Name,
		Sell:      payload.Sell,
		Buy:       payload.Buy,
		CompanyId: payload.CompanyId,
	}

	if err := app.store.Currencies.Create(r.Context(), currency); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, currency); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	currencies, err := app.store.Currencies.GetAll(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, currencies); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)
	var payload UpdateCurrencyPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	currency := &store.Currency{
		ID:        id,
		Name:      payload.Name,
		Sell:      payload.Sell,
		Buy:       payload.Buy,
		CompanyId: payload.CompanyId,
	}

	if err := app.store.Currencies.Update(r.Context(), currency); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, currency); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	if err := app.store.Currencies.Delete(r.Context(), &id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
