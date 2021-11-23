package billing

type APIResponseCreateUser struct {
	ID            string      `json:"id"`
	Object        string      `json:"object"`
	Address       string      `json:"address"`
	Balance       int         `json:"balance"`
	Created       int         `json:"created"`
	Currency      string      `json:"currency"`
	Delinquent    bool        `json:"delinquent"`
	Description   string      `json:"description"`
	Discount      interface{} `json:"discount"`
	Email         string      `json:"email"`
	InvoicePrefix string      `json:"invoice_prefix"`
	Livemode      bool        `json:"livemode"`
	Metadata      struct {
		ID string `json:"id"`
	} `json:"metadata"`
	Name                string   `json:"name"`
	NextInvoiceSequence int      `json:"next_invoice_sequence"`
	Phone               string   `json:"phone"`
	PreferredLocales    []string `json:"preferred_locales"`
	Sources             struct {
		Object     string        `json:"object"`
		Data       []interface{} `json:"data"`
		HasMore    bool          `json:"has_more"`
		TotalCount int           `json:"total_count"`
		URL        string        `json:"url"`
	} `json:"sources"`
	Subscriptions struct {
		Object     string        `json:"object"`
		Data       []interface{} `json:"data"`
		HasMore    bool          `json:"has_more"`
		TotalCount int           `json:"total_count"`
		URL        string        `json:"url"`
	} `json:"subscriptions"`
	TaxExempt string `json:"tax_exempt"`
	TaxIds    struct {
		Object     string        `json:"object"`
		Data       []interface{} `json:"data"`
		HasMore    bool          `json:"has_more"`
		TotalCount int           `json:"total_count"`
		URL        string        `json:"url"`
	} `json:"tax_ids"`
	Error `json:"error"`
}
