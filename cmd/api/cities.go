package main

import "net/http"

type CreateCityPayload struct {
	Name      string  `json:"name"`
	SubName   *string `json:"sub_name"`
	CompanyId int64   `json:"company_id"`
}

type UpdateCityPayload struct {
	Name    string  `json:"name"`
	SubName *string `json:"sub_name"`
}

func (app *application) CreateCityHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllCityHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateCityHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DeleteCityHandler(w http.ResponseWriter, r *http.Request) {

}
