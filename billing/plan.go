package billing

type APIResponsePlan struct {
	ID              string            `json:"id"`
	Object          string            `json:"object"`
	Active          bool              `json:"active"`
	AggregateUsage  interface{}       `json:"aggregate_usage"`
	Amount          int               `json:"amount"`
	AmountDecimal   string            `json:"amount_decimal"`
	BillingScheme   string            `json:"billing_scheme"`
	Created         int               `json:"created"`
	Currency        string            `json:"currency"`
	Interval        string            `json:"interval"`
	IntervalCount   int               `json:"interval_count"`
	Livemode        bool              `json:"livemode"`
	Metadata        map[string]string `json:"metadata"`
	Nickname        interface{}       `json:"nickname"`
	Product         string            `json:"product"`
	Tiers           interface{}       `json:"tiers"`
	TiersMode       interface{}       `json:"tiers_mode"`
	TransformUsage  interface{}       `json:"transform_usage"`
	TrialPeriodDays interface{}       `json:"trial_period_days"`
	UsageType       string            `json:"usage_type"`

	Error `json:"error"`
}
