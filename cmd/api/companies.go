package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CompanyPayload struct {
	Name     string `json:"name"`
	Details  string `json:"details"`
	Password string `json:"password"`
}

func (app *application) CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var payload CompanyPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	company := &store.Company{
		Name:     payload.Name,
		Details:  payload.Details,
		Password: payload.Password,
	}

	if err := app.store.Companies.Create(r.Context(), company); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, company); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllCompanyHandler(w http.ResponseWriter, r *http.Request) {
	companies, err := app.store.Companies.GetAll(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, companies); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetCompanyByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := GetIdFromContext(r)

	company, err := app.store.Companies.GetById(r.Context(), &id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, company); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var payload CompanyPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	company := &store.Company{
		Name:     payload.Name,
		Details:  payload.Details,
		Password: payload.Password,
	}

	if err := app.store.Companies.Update(r.Context(), company); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, company); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	id := GetIdFromContext(r)

	if err := app.store.Companies.Delete(r.Context(), &id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
