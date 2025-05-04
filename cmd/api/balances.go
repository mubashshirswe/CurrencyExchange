package main

import "net/http"

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

}

func (app *application) GetBalanceByIdHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetBalanceByUserIdHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllBalanceHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateBalanceHandler(w http.ResponseWriter, r *http.Request) {

}
