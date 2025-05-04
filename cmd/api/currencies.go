package main

import "net/http"

type CreateCurrencyPayload struct {
	Name      string `json:"name"`
	Sell      *int64 `json:"sell"`
	Buy       *int64 `json:"buy"`
	CompanyId int64  `json:"company_id"`
}

type UpdateCurrencyPayload struct {
	Name string `json:"name"`
	Sell *int64 `json:"sell"`
	Buy  *int64 `json:"buy"`
}

func (app *application) CreateCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DeleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {

}
