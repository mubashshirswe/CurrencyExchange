package main

import "net/http"

type CreateCompanyPayload struct {
	Name     string `json:"name"`
	Details  string `json:"details"`
	Password string `json:"password"`
}

type UpdateCompanyPayload struct {
	Name     string `json:"name"`
	Details  string `json:"details"`
	Password string `json:"password"`
}

func (app *application) CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllCompanyHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetCompanyByIdHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {

}
