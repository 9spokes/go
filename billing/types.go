package billing

import "time"

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
	DocURL  string `json:"doc_url"`
	Param   string `json:"param"`
}

type Invoice struct {
	Date         time.Time `json:"date"`
	Price        int       `json:"price"`
	Subscription string    `json:"subscription"`
}

type Subscription struct {
	Name     string    `json:"name"`
	Price    int       `json:"price"`
	Cycle    string    `json:"cycle"`
	Currency string    `json:"currency"`
	Renew    time.Time `json:"renew"`
	Invoices []Invoice `json:"invoices"`
}
