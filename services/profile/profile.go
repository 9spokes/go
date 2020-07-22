package profile

import (
	"encoding/json"
	"fmt"

	"github.com/9spokes/go/http"
	"github.com/9spokes/go/types"
	goLogging "github.com/op/go-logging"
)

// Context represents a company object into the profile service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *goLogging.Logger
}

//GetCompanies returns a list of companies that the user belongs to
func (ctx Context) GetCompanies(user string) ([]types.Company, error) {

	response, err := http.Request{
		URL: fmt.Sprintf("%s/companies", ctx.URL),
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"x-9sp-user": user,
		},
		ContentType: "application/json",
	}.Get()

	if err != nil {
		return nil, fmt.Errorf("while interacting with Profile service: %s", err.Error())
	}

	var ret struct {
		Status  string          `json:"status"`
		Details []types.Company `json:"details"`
		Message string          `json:"message"`
	}

	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling message: %s", err.Error())
	}

	return ret.Details, nil
}

// GetOptions retrieves all options for the specified user
func (ctx Context) GetOptions(user string) (map[string]interface{}, error) {

	optionsURL := fmt.Sprintf("%s/options", ctx.URL)

	req, err := http.Request{
		URL:     optionsURL,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"x-9sp-user": user,
		},
		ContentType: "application/json",
	}.Get()

	if err != nil {
		return nil, fmt.Errorf("Error getting user profile data %s", err.Error())
	}
	
	var profileOptions map[string]interface{}
	if err := json.Unmarshal(req.Body, &profileOptions); err != nil {
		return nil, fmt.Errorf("while parsing response from profile service options: %s", err.Error())
	}

	return profileOptions, nil
}