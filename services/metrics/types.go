package metrics

import (
	"time"
)

// Query struct
type Query struct {
	Connection string `json:"$connection,omitempty"`
	Filters    string `json:"$filters,omitempty"`
	Sort       string `json:"$sort,omitempty"`
	Limit      int64  `json:"$limit,omitempty"`
	Offset     int64  `json:"$offset,omitempty"`
	Fields     string `json:"$fields,omitempty"`
}

// TimeSeries is a time series metrics response
type TimeSeries struct {
	Query Query       `json:"query"`
	Data  []TimeSeriesDatapoint `json:"data"`
}

// TimeSeriesDatapoint is a single entry in a time series
type TimeSeriesDatapoint struct {
	Time time.Time `json:"time"`
	Value float64 `json:"value"`
}