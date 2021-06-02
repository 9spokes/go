package profile

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/9spokes/go/api"
	"github.com/9spokes/go/http"
	"github.com/9spokes/go/types"
)

// Context represents a company object into the profile service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	User		 string
}

//UpdateProfile updates the user profile based on the contents of the session
func (ctx Context) UpdateProfile(user string, form *url.Values) error {

	response, err := http.Request{
		URL:         ctx.URL,
		ContentType: "application/x-www-form-urlencoded",
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{"x-9sp-user": user},
		Body:    []byte(form.Encode()),
	}.Put()

	if err != nil {
		return err
	}

	var ret api.Response
	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return fmt.Errorf("while unmarshalling response: %s", err.Error())
	}

	if ret.Status != "ok" {
		return fmt.Errorf(ret.Message)
	}

	return nil
}

// GetOptions retrieves all options for the specified user
func (ctx Context) GetOptions(user string) (map[string]interface{}, error) {

	optionsURL := fmt.Sprintf("%s/options", ctx.URL)

	response, err := http.Request{
		URL: optionsURL,
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
		return nil, fmt.Errorf("Error getting user options %s", err.Error())
	}

	var ret api.Response
	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling options: %s", err.Error())
	}

	if ret.Status != "ok" {
		return nil, fmt.Errorf(ret.Message)
	}

	if _, ok := ret.Details.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("return type is %T, not a map", ret.Details)
	}

	return ret.Details.(map[string]interface{}), nil
}

// GetOption looksup one user option by key
func (ctx Context) GetOption(user string, option string) (interface{}, error) {

	optionsURL := fmt.Sprintf("%s/options/%s", ctx.URL, option)

	response, err := http.Request{
		URL: optionsURL,
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
		return nil, fmt.Errorf("Error getting user option %s: %s", option, err.Error())
	}

	var ret api.Response
	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling option %s: %s", option, err.Error())
	}

	if ret.Status != "ok" {
		return nil, fmt.Errorf(ret.Message)
	}

	return ret.Details, nil
}

//GetProfile get user's profile by userId
func (ctx Context) GetProfile(user string) (interface{}, error) {
	profileURL := fmt.Sprintf("%s/users/%s", ctx.URL, user)

	response, err := http.Request{
		URL: profileURL,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"x-9sp-user": ctx.User,
		},
		ContentType: "application/json",
	}.Get()

	if err != nil {
		return nil, fmt.Errorf("error getting user profile %s: %s", user, err.Error())
	}

	var ret api.Response
	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling user profile %s: %s", user, err.Error())
	}

	if ret.Status != "ok" {
		return nil, fmt.Errorf(ret.Message)
	}

	return ret.Details, nil
}

// New takes as parameters and returns a context for subsequent method calls.
func New(url string, creds *types.Credentials) *Context {

	return &Context{
		URL:      url,
		ClientID: creds.Username,
		ClientSecret: creds.Password,
	}
}