package token

import (
	"encoding/json"
	"fmt"

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

// GetConnection returns a connection by ID from the designated Token service instance
func (ctx Context) GetConnection(id string) (types.Connection, error) {

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
		e := fmt.Sprintf("Error invoking token service at: %s: %s", url, err.Error())
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return types.Connection{}, fmt.Errorf(e)
	}

	var parsed struct {
		Status  string           `json:"status"`
		Message string           `json:"message"`
		Details types.Connection `json:"details"`
	}

	if json.Unmarshal(response.Body, &parsed); err != nil {
		e := fmt.Sprintf("Error parsing response from Token service: %s", err.Error())
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return types.Connection{}, fmt.Errorf(e)
	}

	if parsed.Status != "ok" {
		e := fmt.Sprintf("Non-OK response received from Token service: %s", parsed.Message)
		if ctx.Logger != nil {
			ctx.Logger.Error(e)
		}
		return types.Connection{}, fmt.Errorf(e)
	}

	return parsed.Details, nil
}
