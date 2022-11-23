package index

import (
	"encoding/json"
	"errors"
	"fmt"

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

// GetIndexes returns a connection by ID from the designated indexer service instance
func (ctx *Context) GetIndexes(conn string) (map[string][]IndexEntry, error) {

	indexes := make(map[string][]IndexEntry)

	defer func() {
		if r := recover(); r != nil {
			logging.Errorf("An error occured parsing the response from the Indexer service: %s", r)
		}

	}()

	url := fmt.Sprintf("%s/connections/%s", ctx.URL, conn)

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
		return indexes, ErrNotFound
	}

	if err != nil {
		return indexes, fmt.Errorf("error invoking Indexer service at: %s: %s", url, err.Error())
	}

	var parsed struct {
		Status  string                  `json:"status"`
		Message string                  `json:"message"`
		Details map[string][]IndexEntry `json:"details"`
	}

	if err := json.Unmarshal(response.Body, &parsed); err != nil {
		return indexes, fmt.Errorf("error parsing response from Indexer service: %s", err.Error())
	}

	if parsed.Status != "ok" {
		return indexes, fmt.Errorf("non-OK response received from Indexer service: %s", parsed.Message)
	}

	indexes = parsed.Details

	return indexes, nil
}

// UpdateIndex updates an entry with the data provided
func (ctx *Context) UpdateIndex(connection, osp, datasource, cycle string, updatedIndexes []IndexEntry) error {

	location := fmt.Sprintf("%s/connections", ctx.URL)

	logging.Debugf("Invoking Indexer service at: %s", location)

	update := UpdateBody{
		OSP:        osp,
		Connection: connection,
		Datasource: datasource,
		Cycle:      cycle,
		Data:       updatedIndexes,
	}

	body, _ := json.Marshal(update)

	response, err := http.Request{
		URL: location,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/x-www-form-urlencoded",
		Body:        body,
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

// NewIndex creates a new index for a given connection and datasource.  It returns the new index document
func (ctx *Context) NewIndex(index *Index) (*Index, error) {

	// New post-Indexer message
	raw, err := http.Request{
		ContentType: "application/x-www-form-urlencoded",
		URL:         ctx.URL + "/connections",
		Body: []byte(fmt.Sprintf("connection=%s&datasource=%s&count=%d&type=%s&storage=%s&cycle=%s&osp=%s&notify=%t&new_etl=%t&webhooks=%t",
			index.Connection,
			index.Datasource,
			index.Count,
			index.Type,
			index.Storage,
			index.Cycle,
			index.OSP,
			index.Notify,
			index.NewETL,
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
		return nil, fmt.Errorf("error parsing response from indexer service: %s", err.Error())
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf("received an error response from the indexer service: %s", response.Message)
	}

	return &response.Details, nil
}
