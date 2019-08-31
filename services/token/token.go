package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GetConnections returns a list of documents from the Token service that match the criteria set forth in "filter"
func GetConnections(url, filter, limit, selector string) ([]Document, error) {

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

	connections := make([]Document, 0)
	if err := json.Unmarshal(body, &connections); err != nil {
		return nil, fmt.Errorf("while unmarshalling messaging: %s", err.Error())
	}

	return connections, nil
}
