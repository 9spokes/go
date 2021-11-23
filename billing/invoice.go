package billing

type APIResponseInvoice struct {
	ID                   string      `json:"id"`
	Object               string      `json:"object"`
	AccountCountry       string      `json:"account_country"`
	AccountName          string      `json:"account_name"`
	AccountTaxIds        interface{} `json:"account_tax_ids"`
	AmountDue            int         `json:"amount_due"`
	AmountPaid           int         `json:"amount_paid"`
	AmountRemaining      int         `json:"amount_remaining"`
	ApplicationFeeAmount interface{} `json:"application_fee_amount"`
	AttemptCount         int         `json:"attempt_count"`
	Attempted            bool        `json:"attempted"`
	AutoAdvance          bool        `json:"auto_advance"`
	AutomaticTax         struct {
		Enabled bool   `json:"enabled"`
		Status  string `json:"status"`
	} `json:"automatic_tax"`
	BillingReason         string        `json:"billing_reason"`
	Charge                interface{}   `json:"charge"`
	CollectionMethod      string        `json:"collection_method"`
	Created               int           `json:"created"`
	Currency              string        `json:"currency"`
	CustomFields          interface{}   `json:"custom_fields"`
	Customer              string        `json:"customer"`
	CustomerAddress       interface{}   `json:"customer_address"`
	CustomerEmail         string        `json:"customer_email"`
	CustomerName          string        `json:"customer_name"`
	CustomerPhone         string        `json:"customer_phone"`
	CustomerShipping      interface{}   `json:"customer_shipping"`
	CustomerTaxExempt     string        `json:"customer_tax_exempt"`
	CustomerTaxIds        []interface{} `json:"customer_tax_ids"`
	DefaultPaymentMethod  interface{}   `json:"default_payment_method"`
	DefaultSource         interface{}   `json:"default_source"`
	DefaultTaxRates       []interface{} `json:"default_tax_rates"`
	Description           string        `json:"description"`
	Discount              interface{}   `json:"discount"`
	Discounts             []interface{} `json:"discounts"`
	DueDate               int           `json:"due_date"`
	EndingBalance         int           `json:"ending_balance"`
	Footer                interface{}   `json:"footer"`
	HostedInvoiceURL      string        `json:"hosted_invoice_url"`
	InvoicePdf            string        `json:"invoice_pdf"`
	LastFinalizationError interface{}   `json:"last_finalization_error"`
	Lines                 struct {
		Object     string                `json:"object"`
		Data       []APIResponseLineItem `json:"data"`
		HasMore    bool                  `json:"has_more"`
		TotalCount int                   `json:"total_count"`
		URL        string                `json:"url"`
	} `json:"lines"`
	Livemode bool `json:"livemode"`
	Metadata struct {
	} `json:"metadata"`
	NextPaymentAttempt interface{} `json:"next_payment_attempt"`
	Number             string      `json:"number"`
	OnBehalfOf         interface{} `json:"on_behalf_of"`
	Paid               bool        `json:"paid"`
	PaymentIntent      string      `json:"payment_intent"`
	PaymentSettings    struct {
		PaymentMethodOptions interface{} `json:"payment_method_options"`
		PaymentMethodTypes   interface{} `json:"payment_method_types"`
	} `json:"payment_settings"`
	PeriodEnd                    int         `json:"period_end"`
	PeriodStart                  int         `json:"period_start"`
	PostPaymentCreditNotesAmount int         `json:"post_payment_credit_notes_amount"`
	PrePaymentCreditNotesAmount  int         `json:"pre_payment_credit_notes_amount"`
	Quote                        interface{} `json:"quote"`
	ReceiptNumber                interface{} `json:"receipt_number"`
	StartingBalance              int         `json:"starting_balance"`
	StatementDescriptor          interface{} `json:"statement_descriptor"`
	Status                       string      `json:"status"`
	StatusTransitions            struct {
		FinalizedAt           int         `json:"finalized_at"`
		MarkedUncollectibleAt interface{} `json:"marked_uncollectible_at"`
		PaidAt                interface{} `json:"paid_at"`
		VoidedAt              interface{} `json:"voided_at"`
	} `json:"status_transitions"`
	Subscription         string        `json:"subscription"`
	Subtotal             int           `json:"subtotal"`
	Tax                  interface{}   `json:"tax"`
	TaxPercent           interface{}   `json:"tax_percent"`
	Total                int           `json:"total"`
	TotalDiscountAmounts []interface{} `json:"total_discount_amounts"`
	TotalTaxAmounts      []interface{} `json:"total_tax_amounts"`
	TransferData         interface{}   `json:"transfer_data"`
	WebhooksDeliveredAt  int           `json:"webhooks_delivered_at"`

	Error `json:"error"`
}

type APIResponseLineItem struct {
	ID              string        `json:"id"`
	Object          string        `json:"object"`
	Amount          int           `json:"amount"`
	Currency        string        `json:"currency"`
	Description     string        `json:"description"`
	DiscountAmounts []interface{} `json:"discount_amounts"`
	Discountable    bool          `json:"discountable"`
	Discounts       []interface{} `json:"discounts"`
	Livemode        bool          `json:"livemode"`
	Metadata        struct {
	} `json:"metadata"`
	Period struct {
		End   int `json:"end"`
		Start int `json:"start"`
	} `json:"period"`
	Plan             APIResponsePlan  `json:"plan"`
	Price            APIResponsePrice `json:"price"`
	Proration        bool             `json:"proration"`
	Quantity         int              `json:"quantity"`
	Subscription     string           `json:"subscription"`
	SubscriptionItem string           `json:"subscription_item"`
	TaxAmounts       []interface{}    `json:"tax_amounts"`
	TaxRates         []interface{}    `json:"tax_rates"`
	Type             string           `json:"type"`

	Error `json:"error"`
}
