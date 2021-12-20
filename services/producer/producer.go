package producer

import (
	"encoding/json"
	"fmt"

	"github.com/9spokes/go/api"
	"github.com/9spokes/go/http"
	"github.com/9spokes/go/logging/v3"
	"github.com/9spokes/go/types"
)

// Context represents a connection object into the token service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *logging.Logger
}

// ImmediateETL takes a connection ID and notifies the producer to kick off the Immediate ETL cycle for that particular connection.
func (ctx Context) ImmediateETL(conn, osp string) error {

	response, err := http.Request{
		URL: ctx.URL,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
		Body:        []byte(fmt.Sprintf("{\"connection\":\"%s\", \"osp\":\"%s\"}", conn, osp)),
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

func (ctx Context) GetSchedules(organization string, app, datasource *string) ([]types.Schedule, error) {
	logging.Infof("Getting schedules organization %q, app %v, datasource %v", organization, app, datasource)

	query := make(map[string]string)
	if app != nil {
		query["app"] = *app
	}
	if datasource != nil {
		query["datasource"] = *datasource
	}

	response, err := http.Request{
		URL: fmt.Sprintf("%s/organizations/%s/schedule", ctx.URL, organization),
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Query:       query,
		ContentType: "application/json",
	}.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to send request to producer for getting schedules: %s", err.Error())
	}

	var res api.Response
	if err := json.Unmarshal(response.Body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json into response: %s", err.Error())
	}

	schedules, ok := res.Details.([]types.Schedule)
	if !ok {
		return nil, fmt.Errorf("details(%T) from response is not type of []Schedule", res.Details)
	}

	logging.Infof("Fetched %d schedules", len(schedules))
	return schedules, nil
}
