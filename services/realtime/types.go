package realtime

type BankingAccount struct {
	AccountID               string  `json:"account_id,omitempty"`
	AccountActive           bool    `json:"account_active,omitempty"`
	AccountBalance          float64 `json:"account_balance,omitempty"`
	AccountCurrency         string  `json:"account_currency,omitempty"`
	AccountName             string  `json:"account_name,omitempty"`
	AccountType             string  `json:"account_type,omitempty"`
	AccountCategory         string  `json:"account_category,omitempty"`
	AccountStatus           string  `json:"account_status,omitempty"`
	AccountNumber           string  `json:"account_number,omitempty"`
	AccountNumberDisplay    string  `json:"account_number_display,omitempty"`
	AccountNextPayDate      string  `json:"account_next_pay_date,omitempty"`
	AccountMinPayAmount     float64 `json:"account_min_pay_amount,omitempty"`
	AccountLastPayAmount    float64 `json:"account_last_pay_amount,omitempty"`
	AccountCashAdvanceLimit float64 `json:"account_cash_advance_limit,omitempty"`
	AccountPointsAccrued    float64 `json:"account_points_accrued,omitempty"`
	AccountRewardsBalance   float64 `json:"account_rewards_balance,omitempty"`
}
