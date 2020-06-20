package types

import "time"

// IndexerDatasourceRolling is an indexing data entry for a "rolling" datasource
type IndexerDatasourceRolling struct {
	Period  string    `json:"period"`
	Status  string    `json:"status"`
	Retry   bool      `json:"retry"`
	Updated time.Time `json:"updated,omit_empty"`
	Outcome string    `json:"outcome"`
	Index   string    `json:"index"`
}

// IndexerDatasourceAbsolute is an indexing data entry for an "absolute" datasource
type IndexerDatasourceAbsolute struct {
	Status  string    `json:"status"`
	Retry   bool      `json:"retry"`
	Updated time.Time `json:"updated,omit_empty"`
	Outcome string    `json:"outcome"`
	Index   string    `json:"index"`
	Expires time.Time `json:"expires"`
}

// IndexerDatasource is a new data index structure
type IndexerDatasource struct {
	Datasource string      `json:"datasource"`
	Company    string      `json:"company"`
	Type       string      `json:"type"`
	OSP        string      `json:"osp"`
	Count      int64       `json:"count"`
	Cycle      string      `json:"cycle"`
	Storage    string      `json:"storage"`
	Data       interface{} `json:"data"`
}
