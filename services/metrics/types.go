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
type InnerMetric struct {
	Name 	string `json:"name" bson:"name"`
	Osps 	string `json:"osps" bson:"osps"`
	Periods string `json:"periods" bson:"periods"`
	Unit 	string `json:"unit" bson:"unit"`
}

type ResponseDetails struct {
	Category 	string `json:"category" bson:"category"`
	Metrics []InnerMetric `json:"metrics" bson:"metrics"`
}

type MetricServiceResponse struct {
	Details []ResponseDetails `json:"details" bson:"details"`
}