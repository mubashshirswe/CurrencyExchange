package main

import (
	"net/http"

	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type CreateCityPayload struct {
	Name      string `json:"name"`
	ParentId  *int64 `json:"parent_id"`
	CompanyId int64  `json:"company_id"`
}

type UpdateCityPayload struct {
	Name     string `json:"name"`
	ParentId *int64 `json:"parent_id"`
}

func (app *application) CreateCityHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCityPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	city := &store.City{
		Name:      payload.Name,
		ParentId:  payload.ParentId,
		CompanyId: payload.CompanyId,
	}

	if err := app.store.Cities.Create(r.Context(), city); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, city); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetAllCityHandler(w http.ResponseWriter, r *http.Request) {
	cities, err := app.store.Cities.GetAll(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, cities); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) UpdateCityHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateCityPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	city := &store.City{
		Name:     payload.Name,
		ParentId: payload.ParentId,
	}

	if err := app.store.Cities.Update(r.Context(), city); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, city); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) DeleteCityHandler(w http.ResponseWriter, r *http.Request) {
	id := GetIdFromContext(r)
	if err := app.store.Cities.Delete(r.Context(), &id); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
