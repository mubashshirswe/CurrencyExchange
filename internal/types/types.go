package types

type BalanceRecordPayload struct {
	ID               int64  `json:"id"`
	ReceivedMoney    int64  `json:"received_money"`
	ReceivedCurrency string `json:"received_currency"`
	SelledMoney      int64  `json:"selled_money"`
	SelledCurrency   string `json:"selled_currency"`
	UserId           int64  `json:"user_id"`
	CompanyID        int64  `json:"company_id"`
	Details          string `json:"details"`
}
