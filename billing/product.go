package billing

type APIResponseProduct struct {
	ID          string            `json:"id"`
	Object      string            `json:"object"`
	Active      bool              `json:"active"`
	Created     int               `json:"created"`
	Description string            `json:"description"`
	Livemode    bool              `json:"livemode"`
	Metadata    map[string]string `json:"metadata"`
	Name        string            `json:"name"`
	TaxCode     string            `json:"tax_code"`
	UnitLabel   string            `json:"unit_label"`
	Updated     int               `json:"updated"`
	URL         string            `json:"url"`

	Error `json:"error"`
}
