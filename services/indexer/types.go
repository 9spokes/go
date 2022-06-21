package indexer

import (
	"time"
)

// DatasourceRolling is an indexing data entry for a "rolling" datasource
type DatasourceRolling struct {
	Period  string    `json:"period"`
	Status  string    `json:"status"`
	Retry   bool      `json:"retry"`
	Updated time.Time `json:"updated,omit_empty"`
	Outcome string    `json:"outcome"`
	Index   string    `json:"index"`
}

// DatasourceAbsolute is an indexing data entry for an "absolute" datasource
type DatasourceAbsolute struct {
	Status  string    `json:"status"`
	Retry   bool      `json:"retry"`
	Updated time.Time `json:"updated,omit_empty"`
	Outcome string    `json:"outcome"`
	Index   string    `json:"index"`
	Expires time.Time `json:"expires"`
}

// Index is an index entry used to create a new Indexer document
type Index struct {
	Count        int64       `json:"count"`
	Cycle        string      `json:"cycle"`
	Connection   string      `json:"connection,omitempty"`
	Datasource   string      `json:"datasource"`
	Webhooks     bool        `json:"webhooks,omitempty"`
	Notify       bool        `json:"notify,omitempty"`
	OSP          string      `json:"osp"`
	Status       string      `json:"status,omitempty"`
	Storage      string      `json:"storage"`
	Type         string      `json:"type"`
	Dependencies []string    `json:"depends,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// IndexStatus is the response type for getting index status
type IndexStatus struct {
	Datasource  string    `json:"datasource"`
	Status      string    `json:"status"`
	Total       int       `json:"total"`
	Completed   int       `json:"completed"`
	Percent     float64   `json:"percent"`
	LastUpdated time.Time `json:"last_updated"`
}

type ETLMessage struct {
	Connection   string `json:"connection"`
	Datasource   string `json:"datasource"`
	Cycle        string `json:"cycle"`
	Index        string `json:"index"`
	Type         string `json:"type"`
	OSP          string `json:"osp"`
	Outcome      string `json:"outcome"`
	Status       string `json:"status"`
	Retry        string `json:"retry"`
	ImmediateETL bool   `json:"immediate_etl"`
}
