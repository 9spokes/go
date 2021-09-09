package token

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/9spokes/go/http"
	"github.com/9spokes/go/types"
	goLogging "github.com/op/go-logging"
)

// StatusActive is an ACTIVE connection document
const StatusActive = "ACTIVE"

// StatusNotConnected is an NOT_CONNECTED connection document
const StatusNotConnected = "NOT_CONNECTED"

// StatusNew is an NEW connection document
const StatusNew = "NEW"

// Context represents a connection object into the token service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *goLogging.Logger
}

func (ctx Context) InitiateETL(id string) error {

	url := fmt.Sprintf("%s/connections/%s?action=etl", ctx.URL, id)

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Token service at: %s", url)
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
		return err
	}

	var parsed struct {
		Status  string           `json:"status"`
		Message string           `json:"message"`
		Details types.Connection `json:"details"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		return err
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Token service: %s", parsed.Message)
		return fmt.Errorf(e)
	}

	return nil
}

// GetConnection returns a connection by ID from the designated Token service instance
func (ctx Context) GetConnection(id string) (*types.Connection, error) {

	url := fmt.Sprintf("%s/connections/%s", ctx.URL, id)

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Token service at: %s", url)
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
		return nil, err
	}

	var parsed struct {
		Status  string           `json:"status"`
		Message string           `json:"message"`
		Details types.Connection `json:"details"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		return nil, err
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Token service: %s", parsed.Message)
		return nil, fmt.Errorf(e)
	}

	return &parsed.Details, nil
}

//GetConnections returns a list of documents from the Token service that match the criteria set forth in "filter"
func (ctx Context) GetConnections(filter string, limit int, selector string) ([]types.Connection, error) {

	if selector == "" {
		selector = "osp"
	}

	url := fmt.Sprintf("%s/connections?filter=%s&limit=%d&selector=%s",
		ctx.URL,
		filter,
		limit,
		selector,
	)

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
		return nil, fmt.Errorf("while interacting with Token service: %s", err.Error())
	}

	var ret struct {
		Status  string             `json:"status"`
		Details []types.Connection `json:"details"`
		Message string             `json:"message"`
	}

	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details, nil
}

//GetOSP returns an OSP definition from the Token service
func (ctx Context) GetOSP(osp string) (types.Document, error) {

	url := fmt.Sprintf("%s/osp/%s", ctx.URL, osp)

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
		return nil, fmt.Errorf("while interacting with token services: %s", err.Error())
	}

	var ret struct {
		Status  string         `json:"status"`
		Details types.Document `json:"details"`
		Message string         `json:"message"`
	}

	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details, nil
}

// SetConnectionStatus returns a connection by ID from the designated Token service instance
func (ctx Context) SetConnectionStatus(id string, status string, reason string) error {

	if status != StatusNotConnected {
		return fmt.Errorf("cannot set status to %s. %s != %s", status, status, StatusNotConnected)
	}

	connectionUrl := fmt.Sprintf("%s/connections/%s/status", ctx.URL, id)

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Token service at: %s", connectionUrl)
	}

	body := url.Values{}
	body.Add("status", status)
	body.Add("reason", reason)

	response, err := http.Request{
		URL: connectionUrl,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/x-www-form-urlencoded",
		Body:        []byte(body.Encode()),
	}.Post()

	if err != nil {
		return err
	}

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		return err
	}

	if parsed.Status != "ok" {
		return fmt.Errorf("Non-OK response received from Token service: %s", parsed.Message)
	}

	return nil
}

// SetConnectionSetting updates connection setting by ID from the designated Token service instance
func (ctx Context) SetConnectionSetting(id string, settings types.Document) error {

	if settings == nil {
		return fmt.Errorf("The new settings provided is empty")
	}

	url := fmt.Sprintf("%s/connections/%s/settings", ctx.URL, id)

	newSettings, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	response, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
		Body:        newSettings,
	}.Post()

	if err != nil {
		return err
	}

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		return err
	}

	if parsed.Status != "ok" {
		return fmt.Errorf("Non-OK response received from Token service: %s", parsed.Message)
	}

	return nil
}

// CreateConnection requests the designated Token service instance to create a new connection
func (ctx Context) CreateConnection(form map[string]string) (*types.Connection, error) {

	// validate parameters
	if form["osp"] == "" {
		return nil, fmt.Errorf("the field osp is required")
	}
	if form["user"] == "" {
		return nil, fmt.Errorf("the field user is required")
	}

	// sanitizing parameters
	params := url.Values{}
	for _, p := range []string{"osp", "user", "company", "status", "action", "redirect_url", "client_id", "client_secret", "platform"} {
		if form[p] != "" {
			params.Add(p, form[p])
		}
	}

	url := fmt.Sprintf("%s/connections", ctx.URL)

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Token service at: %s", url)
	}

	response, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/x-www-form-urlencoded",
		Body:        []byte(params.Encode()),
	}.Post()

	if err != nil {
		return nil, err
	}

	var parsed struct {
		Status  string           `json:"status"`
		Message string           `json:"message"`
		Details types.Connection `json:"details"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		return nil, err
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Token service: %s", parsed.Message)
		return nil, fmt.Errorf(e)
	}

	return &parsed.Details, nil
}

// RemoveConnection requests the designated Token service instance to remove a connection
func (ctx Context) RemoveConnection(id string) error {

	url := fmt.Sprintf("%s/connections/%s", ctx.URL, id)

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Token service at: %s", url)
	}

	response, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
	}.Delete()

	if err != nil {
		return err
	}

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		return err
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Token service: %s", parsed.Message)
		return fmt.Errorf(e)
	}

	return nil
}

// ManageConnection asks the designated Token service instance to perform an action on the specified connection
func (ctx Context) ManageConnection(id string, action string, params map[string]string) error {

	if action == "" {
		return fmt.Errorf("the action must be specified")
	}

	rawurl := fmt.Sprintf("%s/connections/%s", ctx.URL, id)

	u, err := url.Parse(rawurl)
	if err != nil {
		return fmt.Errorf("Failed to parse Token Svc URL")
	}

	q := u.Query()
	q.Set("action", action)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	if ctx.Logger != nil {
		ctx.Logger.Debugf("Invoking Token service at: %s", u.String())
	}

	response, err := http.Request{
		URL: u.String(),
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
	}.Get()

	if err != nil {
		return err
	}

	var parsed struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Details string `json:"details"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		return err
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Token service: %s", parsed.Message)
		return fmt.Errorf(e)
	}

	return nil
}
