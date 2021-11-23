package billing

type APIResponsePrice struct {
	ID            string      `json:"id"`
	Object        string      `json:"object"`
	Active        bool        `json:"active"`
	BillingScheme string      `json:"billing_scheme"`
	Created       int         `json:"created"`
	Currency      string      `json:"currency"`
	Livemode      bool        `json:"livemode"`
	LookupKey     interface{} `json:"lookup_key"`
	Metadata      struct {
	} `json:"metadata"`
	Nickname  interface{} `json:"nickname"`
	Product   string      `json:"product"`
	Recurring struct {
		AggregateUsage  interface{} `json:"aggregate_usage"`
		Interval        string      `json:"interval"`
		IntervalCount   int         `json:"interval_count"`
		TrialPeriodDays interface{} `json:"trial_period_days"`
		UsageType       string      `json:"usage_type"`
	} `json:"recurring"`
	TaxBehavior       string      `json:"tax_behavior"`
	TiersMode         interface{} `json:"tiers_mode"`
	TransformQuantity interface{} `json:"transform_quantity"`
	Type              string      `json:"type"`
	UnitAmount        int         `json:"unit_amount"`
	UnitAmountDecimal string      `json:"unit_amount_decimal"`

	Error `json:"error"`
}
