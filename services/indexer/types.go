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

// Datasource is a new data index structure
type Datasource struct {
	Datasource string      `json:"datasource"`
	Type       string      `json:"type"`
	OSP        string      `json:"osp"`
	Count      int64       `json:"count"`
	Cycle      string      `json:"cycle"`
	Storage    string      `json:"storage"`
	Data       interface{} `json:"data"`
}

// Index is an index entry used to create a new Indexer document
type Index struct {
	Connection string `json:"connection"`
	OSP        string `json:"osp"`
	Datasource string `json:"string"`
	Count      int    `json:"count"`
	Cycle      string `json:"cycle"`
	Type       string `json:"type"`
	Storage    string `json:"storage"`
}
