package billing

type APIResponseSubscription struct {
	ID                    string      `json:"id"`
	Object                string      `json:"object"`
	ApplicationFeePercent interface{} `json:"application_fee_percent"`
	AutomaticTax          struct {
		Enabled bool `json:"enabled"`
	} `json:"automatic_tax"`
	BillingCycleAnchor   int           `json:"billing_cycle_anchor"`
	BillingThresholds    interface{}   `json:"billing_thresholds"`
	CancelAt             interface{}   `json:"cancel_at"`
	CancelAtPeriodEnd    bool          `json:"cancel_at_period_end"`
	CanceledAt           int           `json:"canceled_at"`
	CollectionMethod     string        `json:"collection_method"`
	Created              int           `json:"created"`
	CurrentPeriodEnd     int64         `json:"current_period_end"`
	CurrentPeriodStart   int64         `json:"current_period_start"`
	Customer             string        `json:"customer"`
	DaysUntilDue         int           `json:"days_until_due"`
	DefaultPaymentMethod string        `json:"default_payment_method"`
	DefaultSource        string        `json:"default_source"`
	DefaultTaxRates      []interface{} `json:"default_tax_rates"`
	Discount             int           `json:"discount"`
	EndedAt              int           `json:"ended_at"`
	Items                struct {
		Object  string                        `json:"object"`
		Data    []APIResponseSubscriptionItem `json:"data"`
		HasMore bool                          `json:"has_more"`
		URL     string                        `json:"url"`
	} `json:"items"`
	LatestInvoice                 string            `json:"latest_invoice"`
	Livemode                      bool              `json:"livemode"`
	Metadata                      map[string]string `json:"metadata"`
	NextPendingInvoiceItemInvoice interface{}       `json:"next_pending_invoice_item_invoice"`
	PauseCollection               bool              `json:"pause_collection"`
	PaymentSettings               struct {
		PaymentMethodOptions interface{} `json:"payment_method_options"`
		PaymentMethodTypes   interface{} `json:"payment_method_types"`
	} `json:"payment_settings"`
	PendingInvoiceItemInterval interface{}     `json:"pending_invoice_item_interval"`
	PendingSetupIntent         interface{}     `json:"pending_setup_intent"`
	PendingUpdate              interface{}     `json:"pending_update"`
	Plan                       APIResponsePlan `json:"plan"`
	Quantity                   int             `json:"quantity"`
	Schedule                   interface{}     `json:"schedule"`
	StartDate                  int             `json:"start_date"`
	Status                     string          `json:"status"`
	TransferData               interface{}     `json:"transfer_data"`
	TrialEnd                   int             `json:"trial_end"`
	TrialStart                 int             `json:"trial_start"`

	Error `json:"error"`
}

type APIResponseSubscriptionItem struct {
	ID                string            `json:"id"`
	Object            string            `json:"object"`
	BillingThresholds interface{}       `json:"billing_thresholds"`
	Created           int               `json:"created"`
	Metadata          map[string]string `json:"metadata"`
	Price             APIResponsePrice  `json:"price"`
	Quantity          int               `json:"quantity"`
	Subscription      string            `json:"subscription"`
	TaxRates          []interface{}     `json:"tax_rates"`

	Error `json:"error"`
}
