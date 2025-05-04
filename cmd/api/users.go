package main

import "net/http"

type CreateUserPayload struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Role     int64  `json:"role"`
	Password string `json:"password"`
}

type UpdateUserPayload struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Role     int64  `json:"role"`
	Password string `json:"password"`
}

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

}
