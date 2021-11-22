package billing

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/9spokes/go/http"
)

type Context struct {
	APIKey      string
	CallbackURL string
}

type Profile struct {
	Name  string
	Email string
	ID    string
}

func (ctx *Context) GetPortalURL(user string) (string, error) {

	body := url.Values{}
	body.Add("customer", user)
	body.Add("return_url", ctx.CallbackURL)

	ret, _ := http.Request{
		URL:            "https://api.stripe.com/v1/billing_portal/sessions",
		ContentType:    "application/x-www-form-urlencoded",
		Authentication: http.Authentication{Scheme: "Basic", Username: ctx.APIKey, Password: ""},
		Body:           []byte(body.Encode()),
	}.Post()

	var response APIResponsePortalURL

	err := json.Unmarshal(ret.Body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response from billing engine: %s", ret.Body)
	}

	if response.Error.Message != "" {
		return "", errors.New(response.Error.Message)
	}

	if response.URL == "" {
		return "", fmt.Errorf("the response did not include the billing URL: %s", ret.Body)
	}

	return response.URL, nil

}

func (ctx *Context) CreateUser(p *Profile) (string, error) {

	if ctx.APIKey == "" {
		return "", fmt.Errorf("API key is required")
	}

	body := url.Values{}
	body.Add("email", p.Email)
	body.Add("name", p.Name)
	body.Add("metadata[id]", p.ID)

	ret, _ := http.Request{
		URL:            "https://api.stripe.com/v1/customers",
		Authentication: http.Authentication{Username: ctx.APIKey, Password: "", Scheme: "basic"},
		ContentType:    "application/x-www-form-urlencoded",
		Body:           []byte(body.Encode()),
	}.Post()

	var response APIResponseCreateUser

	err := json.Unmarshal(ret.Body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response from billing engine: %s", ret.Body)
	}

	if response.Error.Message != "" {
		return "", errors.New(response.Error.Message)
	}

	if response.ID == "" {
		return "", fmt.Errorf("the response did not include the billing ID: %s", ret.Body)
	}

	return response.ID, nil
}

func New(key string, cb string) (*Context, error) {
	if key == "" {
		return nil, fmt.Errorf("the API key is required")
	}

	if cb == "" {
		return nil, fmt.Errorf("the callback URL is required")
	}

	return &Context{
		APIKey:      key,
		CallbackURL: cb,
	}, nil
}
