package profile

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/9spokes/go/api"
	"github.com/9spokes/go/http"
)

// Context represents a company object into the profile service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	User         string
}

//UpdateProfile updates the user profile based on the contents of the session
func (ctx Context) UpdateProfile(form *url.Values) error {

	response, err := http.Request{
		URL:         ctx.URL,
		ContentType: "application/x-www-form-urlencoded",
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{"x-9sp-user": ctx.User},
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

// GetOptions retrieves all options for the specified user. If a filter is specified
// only those options whose names match the filter get returned. The filter can be
// a plain string or a regex.
func (ctx Context) GetUserOptions(filter string) (map[string]interface{}, error) {

	optionsURL := fmt.Sprintf("%s/user/options", ctx.URL)
	
	request := http.Request{
		URL: optionsURL,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"x-9sp-user": ctx.User,
		},
		ContentType: "application/json",
	}
	
	if filter != "" {
		request.Query = map[string]string{
			"q": filter,
		}
	}

	response, err := request.Get()
	if err != nil {
		return nil, fmt.Errorf("error getting user options %s", err.Error())
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

// GetOption retrieves a single user option
func (ctx Context) GetUserOption(option string) (interface{}, error) {

	optionsURL := fmt.Sprintf("%s/user/options/%s", ctx.URL, option)

	response, err := http.Request{
		URL: optionsURL,
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
		return nil, fmt.Errorf("error getting user option %s: %s", option, err.Error())
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
func (ctx Context) GetProfile(user string) (*Profile, error) {
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
		return nil, fmt.Errorf("error getting user profile %s: %v", user, err)
	}

	var ret struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Details Profile
	}

	if err := json.Unmarshal(response.Body, &ret); err != nil {
		return nil, fmt.Errorf("while unmarshalling user profile %s: %v", user, err)
	}

	if ret.Status != "ok" {
		return nil, fmt.Errorf(ret.Message)
	}

	return &ret.Details, nil
}