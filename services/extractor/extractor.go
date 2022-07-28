package extractor

import (
	"encoding/json"
	"fmt"

	"github.com/9spokes/go/http"
	"github.com/9spokes/go/logging/v2"
)

type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *logging.Logger
}

// ImmediateETL takes a connection ID and notifies the extractor to kick off the connection based Immediate ETL.
func (ctx Context) ImmediateETL(conn string) error {

	response, err := http.Request{
		URL: ctx.URL,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
		Body:        []byte(fmt.Sprintf("{\"connection\":\"%s\"}", conn)),
	}.Post()

	if err != nil {
		return fmt.Errorf("while interacting with extractor-ng service: %s", err.Error())
	}

	var ret struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	if ret.Status != "ok" {
		return fmt.Errorf(ret.Message)
	}

	return nil
}
