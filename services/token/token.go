package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Connection is a short-form connection object as returned by the token service, it excludes sensitive info and is meant as a summary
type Connection struct {
	ID       string    `json:"id,omitempty"`
	Created  time.Time `json:"created" bson:"created,omitempty"`
	Modified time.Time `json:"modified" bson:"modified,omitempty"`
	OSP      string    `json:"osp" bson:"osp,omitempty"`
	Status   string    `json:"status" bson:"status,omitempty"`
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
