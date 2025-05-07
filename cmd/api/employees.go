package main

// import (
// 	"net/http"

// 	"github.com/mubashshir3767/currencyExchange/internal/env"
// 	"github.com/mubashshir3767/currencyExchange/internal/store"
// )

// type EmployeePayload struct {
// 	Username  string `json:"username"`
// 	Phone     string `json:"phone"`
// 	Role      int64  `json:"role"`
// 	Password  string `json:"password"`
// 	CompanyId int64  `json:"company_id"`
// }

// type LoginEmployeePayload struct {
// 	Phone    string `json:"phone"`
// 	Password string `json:"password"`
// }

// func (app *application) CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload UserPayload
// 	if err := readJSON(w, r, &payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	employee := &store.Employee{
// 		Username: payload.Username,
// 		Phone:    payload.Phone,
// 		Role:     payload.Role,
// 		Password: payload.Password,
// 	}

// 	if err := app.store.Employees.Create(r.Context(), employee); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	if err := app.writeResponse(w, http.StatusOK, employee); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// }

// func (app *application) LoginEmployeeHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload EmployeePayload
// 	if err := readJSON(w, r, &payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	employee := &store.Employee{
// 		Phone:    payload.Phone,
// 		Password: payload.Password,
// 	}

// 	if err := app.store.Employees.Login(r.Context(), employee); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	token, err := JWTCreate([]byte(env.GetString("JWTSECRET", "secret")), employee.ID)
// 	if err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	if err := app.writeResponse(w, http.StatusOK, token); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// }

// func (app *application) GetAllEmployeeHandler(w http.ResponseWriter, r *http.Request) {
// 	employees, err := app.store.Employees.GetAll(r.Context())
// 	if err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	if err := app.writeResponse(w, http.StatusOK, employees); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// }

// func (app *application) UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload EmployeePayload
// 	if err := readJSON(w, r, &payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	employee := &store.Employee{
// 		Username: payload.Username,
// 		Phone:    payload.Phone,
// 		Role:     payload.Role,
// 		Password: payload.Password,
// 	}

// 	if err := app.store.Employees.Update(r.Context(), employee); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	if err := app.writeResponse(w, http.StatusOK, employee); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// }

// func (app *application) DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
// 	id := GetIdFromContext(r)

// 	if err := app.store.Employees.Delete(r.Context(), &id); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	if err := app.writeResponse(w, http.StatusOK, "DELETED"); err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}
// }
