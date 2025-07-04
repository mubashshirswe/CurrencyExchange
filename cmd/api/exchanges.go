package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type ExchangePayload struct {
	ID               *int64 `json:"id"`
	ReceivedMoney    int64  `json:"received_money"`
	ReceivedCurrency string `json:"received_currency"`
	SelledMoney      int64  `json:"selled_money"`
	SelledCurrency   string `json:"selled_currency"`
	UserId           int64  `json:"user_id"`
	CompanyID        *int64 `json:"company_id"`
	Details          string `json:"details"`
}

func (app *application) CreateExchangeHandler(w http.ResponseWriter, r *http.Request) {
	var payload ExchangePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	exchange := &store.Exchange{
		ReceivedMoney:    payload.ReceivedMoney,
		ReceivedCurrency: payload.ReceivedCurrency,
		SelledMoney:      payload.SelledMoney,
		SelledCurrency:   payload.SelledCurrency,
		UserId:           payload.UserId,
		Details:          &payload.Details,
	}

	if err := app.service.Exchanges.Create(r.Context(), exchange); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetExchangesHandler(w http.ResponseWriter, r *http.Request) {
	var payload FieldRequestPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	records, err := app.store.Exchanges.GetByField(r.Context(), payload.FieldName, payload.FieldValue)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, records); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateExchangeHandler(w http.ResponseWriter, r *http.Request) {
	var payload ExchangePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	exchange := &store.Exchange{
		ID:               getIDFromContext(r),
		ReceivedMoney:    payload.ReceivedMoney,
		ReceivedCurrency: payload.ReceivedCurrency,
		SelledMoney:      payload.SelledMoney,
		SelledCurrency:   payload.SelledCurrency,
		UserId:           payload.UserId,
		Details:          &payload.Details,
	}

	if err := app.service.Exchanges.Update(r.Context(), exchange); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteExchangeHandler(w http.ResponseWriter, r *http.Request) {
	id := getIDFromContext(r)

	if err := app.service.Exchanges.Delete(r.Context(), id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
