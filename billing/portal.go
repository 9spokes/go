package billing

type APIResponsePortalURL struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	Configuration string `json:"configuration"`
	Created       int    `json:"created"`
	Customer      string `json:"customer"`
	Livemode      bool   `json:"livemode"`
	ReturnURL     string `json:"return_url"`
	URL           string `json:"url"`

	Error `json:"error"`
}
