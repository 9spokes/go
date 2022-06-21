package index

import "time"

type Index struct {
	Count      int64        `json:"count"`
	Cycle      string       `json:"cycle"`
	Connection string       `json:"connection,omitempty"`
	Datasource string       `json:"datasource"`
	Webhooks   bool         `json:"webhooks,omitempty"`
	Notify     bool         `json:"notify,omitempty"`
	OSP        string       `json:"osp"`
	Status     string       `json:"status,omitempty"`
	Storage    string       `json:"storage"`
	Type       string       `json:"type"`
	NewETL     bool         `json:"new_etl"`
	Data       []IndexEntry `json:"data,omitempty"`
}

type IndexEntry struct {
	Period  string    `json:"period"`
	Status  string    `json:"status"`
	Retry   bool      `json:"retry"`
	Updated time.Time `json:"updated,omitempty"`
	Outcome string    `json:"outcome"`
	Index   string    `json:"index"`
	Cycle   string    `json:"cycle"`
}

type UpdateBody struct {
	Connection string       `json:"connection"`
	Datasource string       `json:"datasource"`
	Cycle      string       `json:"cycle"`
	Data       []IndexEntry `json:"data"`
}
