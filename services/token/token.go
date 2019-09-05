package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Document is a short-hand for map[string]interface{}
type Document map[string]interface{}

// Connection represents a Connection document object as stored in the database
type Connection struct {
	ID            string    `json:"id" bson:"connection"`
	Credentials   Document  `bson:"credentials,omitempty" json:"credentials"`
	Token         Document  `bson:"token,omitempty" json:"token"`
	Settings      Document  `bson:"settings,omitempty" json:"settings"`
	User          string    `bson:"user,omitempty" json:"user"`
	Configuration Document  `bson:"config,omitempty" json:"config"`
	OSP           string    `bson:"osp,omitempty" json:"osp"`
	Created       time.Time `bson:"created,omitempty" json:"created"`
	Modified      time.Time `bson:"modified,omitempty" json:"modified"`
	Status        string    `bson:"status,omitempty" default:"NEW" json:"status"`
}

//GetConnections returns a list of documents from the Token service that match the criteria set forth in "filter"
func GetConnections(url, filter, limit, selector string) ([]Connection, error) {

	res, err := http.Get(fmt.Sprintf("%s/connections?filter=%s&limit=%s&selector=%s",
		url,
		filter,
		limit,
		selector,
	))
	if err != nil {
		return nil, fmt.Errorf("while interacting with token services: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading message body: %s", err.Error())
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received a non-OK response from token service: [%s] %s", res.Status, body)
	}

	var ret struct {
		Status  string       `json:"status"`
		Details []Connection `json:"details"`
		Message string       `json:"message"`
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details, nil
}

//CreateConnection creates a new connection document and optionally runs the "Initiate" state to prime the document
func CreateConnection(url, user, app string) (Connection, error) {

	res, err := http.Post(
		fmt.Sprintf("%s/connections", url),
		"application/x-www-form-urlencoded",
		bytes.NewReader([]byte(fmt.Sprintf("osp=%s&user=%s", app, user))),
	)
	if err != nil {
		return Connection{}, fmt.Errorf("while interacting with token services: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Connection{}, fmt.Errorf("while reading message body: %s", err.Error())
	}

	if res.StatusCode != 200 {
		return Connection{}, fmt.Errorf("received a non-OK response from token service: [%s] %s", res.Status, body)
	}

	var ret struct {
		Status  string     `json:"status"`
		Details Connection `json:"details"`
		Message string     `json:"message"`
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return Connection{}, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details, nil
}

//ActivateConnection transitions a connection from authorized to active by swapping the authoriztion code for an access token
func ActivateConnection(url, user, app string, form url.Values) (Connection, error) {

	res, err := http.Get(fmt.Sprintf("%s/connections/%s?action=authorize&%s",
		url,
		app,
		form.Encode(),
	))
	if err != nil {
		return Connection{}, fmt.Errorf("while interacting with token services: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Connection{}, fmt.Errorf("while reading message body: %s", err.Error())
	}

	if res.StatusCode != 200 {
		return Connection{}, fmt.Errorf("received a non-OK response from token service: [%s] %s", res.Status, body)
	}

	var ret struct {
		Status  string     `json:"status"`
		Details Connection `json:"details"`
		Message string     `json:"message"`
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return Connection{}, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details, nil
}

//GetOSP returns an OSP definition from the Token service
func GetOSP(url, osp string) (Document, error) {

	res, err := http.Get(fmt.Sprintf("%s/osp/%s", url, osp))
	if err != nil {
		return nil, fmt.Errorf("while interacting with token services: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading message body: %s", err.Error())
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received a non-OK response from token service: [%s] %s", res.Status, body)
	}

	var ret struct {
		Status  string   `json:"status"`
		Details Document `json:"details"`
		Message string   `json:"message"`
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details, nil
}
