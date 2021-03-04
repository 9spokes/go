package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/9spokes/go/http"
)

// Context represents a company object into the profile service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
}

//GetTimeSeriesMetric calculates a metric based on the supplied criteria and returns an answer.
func (ctx Context) GetTimeSeriesMetric(q Query) (*TimeSeries, error) {

	res, err := http.Request{
		URL: ctx.URL,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
	}.Get()

	if err != nil {
		return nil, err
	}

	type response struct {
		Status  string     `json:"status"`
		Message string     `json:"message"`
		Details TimeSeries `json:"details"`
	}

	var ret response
	if err := json.Unmarshal(res.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling response: %s", err.Error())
	}

	if ret.Status != "ok" {
		return nil, fmt.Errorf(ret.Message)
	}

	return &ret.Details, nil
}
