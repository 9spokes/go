package indexer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/9spokes/go/http"
	"github.com/9spokes/go/types"
	goLogging "github.com/op/go-logging"
)

// Context represents a connection object into the token service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *goLogging.Logger
}

// GetIndex returns a connection by ID from the designated indexer service instance
func (ctx Context) GetIndex(company, osp, datasource, cycle string) (*types.IndexerDatasource, error) {

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%s", r)
			if ctx.Logger != nil {
				ctx.Logger.Errorf("an error occured parsing the response from the Indexer service: %s", err)
			}
		}
	}()

	url := fmt.Sprintf("%s/%s/%s/%s?cycle=%s", ctx.URL, company, osp, datasource, cycle)

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Indexer service at: %s", url)
	}

	response, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
	}.Get()

	if err != nil {
		return nil, fmt.Errorf("error invoking Indexer service at: %s: %s", url, err.Error())
	}

	var parsed struct {
		Status  string                  `json:"status"`
		Message string                  `json:"message"`
		Details types.IndexerDatasource `json:"details"`
	}

	if err := json.Unmarshal(response.Body, &parsed); err != nil {
		return nil, fmt.Errorf("error parsing response from Indexer service: %s", err.Error())
	}

	if parsed.Status != "ok" {
		return nil, fmt.Errorf("non-OK response received from Indexer service: %s", parsed.Message)
	}

	switch parsed.Details.Type {
	case "rolling":
		data := make([]types.IndexerDatasourceRolling, len(parsed.Details.Data.([]interface{})))

		for i, e := range parsed.Details.Data.([]interface{}) {
			updated, _ := time.Parse(time.RFC3339, e.(map[string]interface{})["updated"].(string))
			data[i] = types.IndexerDatasourceRolling{
				Index:   e.(map[string]interface{})["index"].(string),
				Outcome: e.(map[string]interface{})["outcome"].(string),
				Period:  e.(map[string]interface{})["period"].(string),
				Retry:   e.(map[string]interface{})["retry"].(bool),
				Status:  e.(map[string]interface{})["status"].(string),
				Updated: updated,
			}
		}
		parsed.Details.Data = data

	case "absolute":
		e := parsed.Details.Data.(interface{})
		updated, _ := time.Parse(time.RFC3339, e.(map[string]interface{})["updated"].(string))
		expires, _ := time.Parse(time.RFC3339, e.(map[string]interface{})["expires"].(string))

		parsed.Details.Data = types.IndexerDatasourceAbsolute{
			Index:   e.(map[string]interface{})["index"].(string),
			Outcome: e.(map[string]interface{})["outcome"].(string),
			Expires: expires,
			Retry:   e.(map[string]interface{})["retry"].(bool),
			Status:  e.(map[string]interface{})["status"].(string),
			Updated: updated,
		}
	}
	return &parsed.Details, nil
}

// UpdateIndex updates an entry with the data provided
func (ctx Context) UpdateIndex(company, osp, datasource, cycle, index, outcome string, ok, retry bool) error {

	location := fmt.Sprintf("%s/%s/%s/%s?cycle=%s&index=%s", ctx.URL, company, osp, datasource, cycle, index)

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Indexer service at: %s", location)
	}

	status := "ok"
	if !ok {
		status = "err"
	}

	params := url.Values{}
	params.Add("outcome", outcome)
	params.Add("status", status)
	params.Add("retry", fmt.Sprintf("%t", retry))

	response, err := http.Request{
		URL: location,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/x-www-form-urlencoded",
		Body:        []byte(params.Encode()),
	}.Put()

	if err != nil && response.StatusCode < 399 {
		return fmt.Errorf("error invoking Indexer service at: %s: [%d] %s", location, response.StatusCode, err.Error())
	}

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(response.Body, &parsed); err != nil {
		return fmt.Errorf("error parsing response from Indexer service: %s", err.Error())
	}

	if parsed.Status != "ok" {
		return fmt.Errorf("non-OK response received from Indexer service: %s", parsed.Message)
	}

	return nil
}
