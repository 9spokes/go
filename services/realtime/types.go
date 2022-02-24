package realtime

type BankingAccount struct {
	AccountID       string  `json:"account_id,omitempty"`
	AccountActive   bool    `json:"account_active,omitempty"`
	AccountBalance  float64 `json:"account_balance,omitempty"`
	AccountCurrency string  `json:"account_currency,omitempty"`
	AccountName     string  `json:"account_name,omitempty"`
	AccountType     string  `json:"account_type,omitempty"`
	AccountCategory string  `json:"account_category,omitempty"`
	AccountStatus   string  `json:"account_status,omitempty"`
}
