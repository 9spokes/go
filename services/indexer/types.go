package indexer

import "time"

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
	Notify       bool        `json:"notify,omitempty"`
	OSP          string      `json:"osp"`
	Status       string      `json:"status,omitempty"`
	Storage      string      `json:"storage"`
	Type         string      `json:"type"`
	Dependencies []string    `json:"depends,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}
