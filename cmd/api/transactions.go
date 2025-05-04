package main

import "net/http"

type CreateTransactionPayload struct {
	Amount             int64  `json:"amount"`
	ServiceFee         int64  `json:"service_fee"`
	FromCurrencyTypeId int64  `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64  `json:"to_currency_type_id"`
	SenderId           int64  `json:"sender_id"`
	FromCityId         int64  `json:"from_city_id"`
	ToCityId           int64  `json:"to_city_id"`
	ReceiverName       string `json:"receiver_name"`
	ReceiverPhone      string `json:"receiver_phone"`
	Details            string `json:"details"`
	Type               int64  `json:"type"`
	CompanyId          int64  `json:"company_id"`
	BalanceId          int64  `json:"balance_id"`
}

type UpdateTransactionPayload struct {
	Amount             int64  `json:"amount"`
	ServiceFee         int64  `json:"service_fee"`
	FromCurrencyTypeId int64  `json:"from_currency_type_id"`
	ToCurrencyTypeId   int64  `json:"to_currency_type_id"`
	SenderId           int64  `json:"sender_id"`
	FromCityId         int64  `json:"from_city_id"`
	ToCityId           int64  `json:"to_city_id"`
	ReceiverName       string `json:"receiver_name"`
	ReceiverPhone      string `json:"receiver_phone"`
	Details            string `json:"details"`
	Type               int64  `json:"type"`
}

func (app *application) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllTransactionHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetAllTransactionByDateHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) GetTransactionByIdHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {

}
