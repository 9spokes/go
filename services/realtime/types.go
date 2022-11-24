package realtime

type BankingAccount struct {
	ID               string  `json:"account_id,omitempty"`
	Active           bool    `json:"account_active,omitempty"`
	CurrentBalance   float64 `json:"account_current_balance,omitempty"`
	AvailableBalance float64 `json:"account_available_balance,omitempty"`
	Currency         string  `json:"account_currency,omitempty"`
	Name             string  `json:"account_name,omitempty"`
	Type             string  `json:"account_type,omitempty"`
	Category         string  `json:"account_category,omitempty"`
	Status           string  `json:"account_status,omitempty"`
	Number           string  `json:"account_number,omitempty"`
	NumberDisplay    string  `json:"account_number_display,omitempty"`
	NextPayDate      string  `json:"account_next_pay_date,omitempty"`
	MinPayAmount     float64 `json:"account_min_pay_amount,omitempty"`
	LastPayAmount    float64 `json:"account_last_pay_amount,omitempty"`
	CashAdvanceLimit float64 `json:"account_cash_advance_limit,omitempty"`
	PointsAccrued    float64 `json:"account_points_accrued,omitempty"`
	RewardsBalance   float64 `json:"account_rewards_balance,omitempty"`
}
