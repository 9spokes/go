package webhooks

import (
	"fmt"

	"github.com/9spokes/go/http"
	goLogging "github.com/op/go-logging"
)

// Context represents a connection object into the webhooks service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *goLogging.Logger
}

// CreateWebhook creates a new Webhook listener
func (ctx *Context) CreateWebhook(osp string, connection string) error {

	url := fmt.Sprintf("%s/%s/connections/%s", ctx.URL, osp, connection)
	_, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"internal-service": "true",
		},
		ContentType: "application/json",
	}.Post()

	return err
}

// DeleteWebook removes a webhook listener
func (ctx *Context) DeleteWebhook(osp string, connection string) error {

	url := fmt.Sprintf("%s/%s/connections/%s", ctx.URL, osp, connection)
	_, err := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"internal-service": "true",
		},
		ContentType: "application/json",
	}.Delete()

	return err
}
