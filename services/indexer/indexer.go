package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/9spokes/go/http"
	"github.com/9spokes/go/logging/v3"
)

// Context represents a connection object into the token service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *logging.Logger
}

var (
	ErrNotFound = errors.New("not found")
)

// NewIndex creates a new index for a given connection and datasource.  It returns the new index document
func (ctx *Context) NewIndex(index *Index) (*Index, error) {

	// New post-Indexer message
	raw, err := http.Request{
		ContentType: "application/x-www-form-urlencoded",
		URL:         ctx.URL + "/connections",
		Body: []byte(fmt.Sprintf("connection=%s&datasource=%s&count=%d&type=%s&storage=%s&cycle=%s&osp=%s&notify=%t&depends=%s&webhooks=%t",
			index.Connection,
			index.Datasource,
			index.Count,
			index.Type,
			index.Storage,
			index.Cycle,
			index.OSP,
			index.Notify,
			strings.Join(index.Dependencies, ","),
			index.Webhooks,
		)),
		Authentication: http.Authentication{Scheme: "basic", Username: ctx.ClientID, Password: ctx.ClientSecret},
	}.Post()
	if err != nil {
		return nil, fmt.Errorf("failed to create new index: %s", err.Error())
	}

	var response struct {
		Status  string `json:"status,omitempty"`
		Message string `json:"message,omitempty"`
		Details Index  `json:"details,omitempty"`
	}
	if err := json.Unmarshal(raw.Body, &response); err != nil {
		return nil, fmt.Errorf("Error parsing response from indexer service: %s", err.Error())
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf("Received an error response from the indexer service: %s", response.Message)
	}

	return response.Details.updateData()
}

func (ds *Index) updateData() (*Index, error) {

	switch ds.Type {
	case "rolling":
		data := make([]DatasourceRolling, len(ds.Data.([]interface{})))

		for i, e := range ds.Data.([]interface{}) {

			skip := false

			for _, key := range []string{"index", "outcome", "period", "retry", "status", "updated"} {
				if _, ok := e.(map[string]interface{})[key]; !ok {
					//logging.Errorf("Failed to parsed '%s' as a string", key)
					skip = true
				}
			}

			if skip {
				continue
			}

			updated, _ := time.Parse(time.RFC3339, e.(map[string]interface{})["updated"].(string))
			data[i] = DatasourceRolling{
				Index:   e.(map[string]interface{})["index"].(string),
				Outcome: e.(map[string]interface{})["outcome"].(string),
				Period:  e.(map[string]interface{})["period"].(string),
				Retry:   e.(map[string]interface{})["retry"].(bool),
				Status:  e.(map[string]interface{})["status"].(string),
				Updated: updated,
			}
		}
		ds.Data = data

	case "absolute":
		e := ds.Data.(interface{})
		updated, _ := time.Parse(time.RFC3339, e.(map[string]interface{})["updated"].(string))
		expires, _ := time.Parse(time.RFC3339, e.(map[string]interface{})["expires"].(string))

		ds.Data = DatasourceAbsolute{
			Index:   e.(map[string]interface{})["index"].(string),
			Outcome: e.(map[string]interface{})["outcome"].(string),
			Expires: expires,
			Retry:   e.(map[string]interface{})["retry"].(bool),
			Status:  e.(map[string]interface{})["status"].(string),
			Updated: updated,
		}
	}
	return ds, nil
}

// GetIndex returns a connection by ID from the designated indexer service instance
func (ctx *Context) GetIndex(conn, datasource, cycle string) (*Index, error) {

	defer func() {
		if r := recover(); r != nil {
			logging.Errorf("An error occured parsing the response from the Indexer service: %s", r)
		}

	}()

	url := fmt.Sprintf("%s/connections/%s/%s?cycle=%s", ctx.URL, conn, datasource, cycle)

	logging.Debugf("Invoking Indexer service at: %s", url)

	response, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
	}.Get()

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error invoking Indexer service at: %s: %s", url, err.Error())
	}

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Details Index  `json:"details"`
	}

	if err := json.Unmarshal(response.Body, &parsed); err != nil {
		return nil, fmt.Errorf("error parsing response from Indexer service: %s", err.Error())
	}

	if parsed.Status != "ok" {
		return nil, fmt.Errorf("non-OK response received from Indexer service: %s", parsed.Message)
	}

	return parsed.Details.updateData()
}

// UpdateIndex updates an entry with the data provided
func (ctx *Context) UpdateIndex(conn, datasource, cycle, index, outcome string, ok, retry bool) error {

	location := fmt.Sprintf("%s/connections/%s/%s?cycle=%s&index=%s", ctx.URL, conn, datasource, cycle, index)

	logging.Debugf("Invoking Indexer service at: %s", location)

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

	if err != nil {
		return fmt.Errorf("error invoking Indexer service at: %s: %s", location, err.Error())
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

// GetIndex returns the datasource status from the designated indexer service
// instance. If the datasource is not specified the statuses of all datasources
// belonging to the specified connection are returned.
func (ctx *Context) GetDatasourceStatus(conn, datasource string) ([]*IndexStatus, error) {

	var url string
	if datasource == "" {
		url = fmt.Sprintf("%s/connections/%s/status", ctx.URL, conn)
	} else {
		url = fmt.Sprintf("%s/connections/%s/%s/status", ctx.URL, conn, datasource)
	}

	logging.Debugf("Invoking Indexer service at: %s", url)

	response, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
	}.Get()

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error invoking Indexer service at: %s: %s", url, err.Error())
	}

	if datasource != "" {
		var parsed struct {
			Status  string       `json:"status"`
			Message string       `json:"message"`
			Details *IndexStatus `json:"details"`
		}

		if err := json.Unmarshal(response.Body, &parsed); err != nil {
			return nil, fmt.Errorf("error parsing response from Indexer service: %s", err.Error())
		}

		if parsed.Status != "ok" {
			return nil, fmt.Errorf("non-OK response received from Indexer service: %s", parsed.Message)
		}

		return []*IndexStatus{parsed.Details}, nil
	}

	var parsed struct {
		Status  string         `json:"status"`
		Message string         `json:"message"`
		Details []*IndexStatus `json:"details"`
	}

	if err := json.Unmarshal(response.Body, &parsed); err != nil {
		return nil, fmt.Errorf("error parsing response from Indexer service: %s", err.Error())
	}

	if parsed.Status != "ok" {
		return nil, fmt.Errorf("non-OK response received from Indexer service: %s", parsed.Message)
	}

	return parsed.Details, nil
}
