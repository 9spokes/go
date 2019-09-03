package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Connection is a short-form connection object as returned by the token service, it excludes sensitive info and is meant as a summary
type Connection struct {
	ID           string    `json:"id,omitempty"`
	Created      time.Time `json:"created" bson:"created,omitempty"`
	Modified     time.Time `json:"modified" bson:"modified,omitempty"`
	OSP          string    `json:"osp" bson:"osp,omitempty"`
	Status       string    `json:"status" bson:"status,omitempty"`
	AuthorizeURL string    `json:"authorize_url,omitempty"`
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
		Status  string `json:"status"`
		Details struct {
			Connections []Connection `json:"connections"`
		} `json:"details"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details.Connections, nil
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
		Status  string `json:"status"`
		Details struct {
			ID          string                 `json:"id"`
			Created     time.Time              `json:"created"`
			Modified    time.Time              `json:"modified"`
			OSP         string                 `json:"osp"`
			Status      string                 `json:"status"`
			Credentials map[string]interface{} `json:"credentials"`
		} `json:"details"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return Connection{}, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	// created, _ := time.Parse("2006-01-02T15:04:05.000Z", ret.Details["created"].(string))
	// modified, _ := time.Parse("2006-01-02T15:04:05.000Z", ret.Details["modified"].(string))

	return Connection{
		AuthorizeURL: ret.Details.Credentials["authorize_url"].(string),
		ID:           ret.Details.ID,
		Created:      ret.Details.Created,
		Modified:     ret.Details.Modified,
		OSP:          ret.Details.OSP,
		Status:       ret.Details.Status,
	}, nil
}
