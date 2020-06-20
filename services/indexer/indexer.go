package indexer

import (
	"encoding/json"
	"fmt"
	"net/url"

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
		e := fmt.Sprintf("Error invoking Indexer service at: %s: %s", url, err.Error())
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return nil, fmt.Errorf(e)
	}

	var parsed struct {
		Status  string                  `json:"status"`
		Message string                  `json:"message"`
		Details types.IndexerDatasource `json:"details"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		e := fmt.Sprintf("Error parsing response from Indexer service: %s", err.Error())
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return nil, fmt.Errorf(e)
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Token service: %s", parsed.Message)
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return nil, fmt.Errorf(e)
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

	if err != nil {
		e := fmt.Sprintf("Error invoking Indexer service at: %s: [%d] %s", location, response.StatusCode, err.Error())
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return fmt.Errorf(e)
	}

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		e := fmt.Sprintf("Error parsing response from Indexer service: %s", err.Error())
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return fmt.Errorf(e)
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Indexer service: %s", parsed.Message)
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return fmt.Errorf(e)
	}

	return nil
}
