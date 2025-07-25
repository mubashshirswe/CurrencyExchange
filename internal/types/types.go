package types

type BalanceRecordPayload struct {
	ID               int64  `json:"id"`
	ReceivedMoney    int64  `json:"received_money"`
	ReceivedCurrency string `json:"received_currency"`
	SelledMoney      int64  `json:"selled_money"`
	SelledCurrency   string `json:"selled_currency"`
	UserId           int64  `json:"user_id"`
	CompanyID        *int64 `json:"company_id"`
	Details          string `json:"details"`
}

type TransactionComplete struct {
	TransactionID      int64 `json:"transactionID"`
	DeliveredUserId    int64 `json:"delivered_user_id"`
	RecievedServiceFee int64 `json:"received_service_fee"`
}

type ReceivedIncomes struct {
	ReceivedAmount   int64  `json:"received_amount"`
	ReceivedCurrency string `json:"received_currency"`
}

type DeliveredOutcomes struct {
	DeliveredAmount   int64  `json:"delivered_amount"`
	DeliveredCurrency string `json:"delivered_currency"`
}

type Pagination struct {
	Page                int         `json:"page"`
	Limit               int         `json:"limit"`
	Data                interface{} `json:"data"`
	Offset              int         `json:"offset"`
	TaskId              int         `json:"taskId"`
	ProductCollectionId int         `json:"productCollectionId"`
	UserId              int64       `json:"user_id"`
	Language            *string     `json:"language"`
}
