package producer

import (
	"encoding/json"
	"fmt"

	"github.com/9spokes/go/http"
	goLogging "github.com/op/go-logging"
)

// Context represents a connection object into the token service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *goLogging.Logger
}

// ImmediateETL takes a connection ID and notifies the producer to kick off the Immediate ETL cycle for that particular connection.
func (ctx Context) ImmediateETL(conn string) error {

	response, err := http.Request{
		URL: ctx.URL,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
		Body:        []byte(fmt.Sprintf("{\"connectionId\":\"%s\"}", conn)),
	}.Post()

	if err != nil {
		return fmt.Errorf("while interacting with Producer service: %s", err.Error())
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
