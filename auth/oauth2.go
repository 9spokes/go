package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// OAuth2 represents the minimum fields required to perform an OAuth2 token exchange or token refresh.
type OAuth2 struct {
	ClientID      string
	ClientSecret  string
	RefreshToken  string
	RedirectURI   string
	Code          string
	TokenEndpoint string
	Extras        map[string]string
	Headers       map[string]string
	Client        *http.Client
}

// Options are a set of flags & modifiers to the OAuth2 implementation
type Options struct {
	AuthInHeader bool `default:"false"`
}

// Authorize implements an OAuth2 authorization using the parameters defined in the OAuth2 struct
func (params OAuth2) Authorize(opt Options) (map[string]interface{}, error) {

	if params.ClientID == "" || params.ClientSecret == "" {
		return nil, fmt.Errorf("client_id and client_secret cannot be empty")
	}

	if params.Code == "" {
		return nil, fmt.Errorf("the authorization code is missing")
	}

	var auth string

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("redirect_uri", params.RedirectURI)

	for k, v := range params.Extras {
		data.Set(k, v)
	}

	if opt.AuthInHeader {
		auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(params.ClientID+":"+params.ClientSecret))
	} else {
		data.Set("client_id", params.ClientID)
		data.Set("client_secret", params.ClientSecret)
	}

	request, err := http.NewRequest("POST", params.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create an HTTP client: %s", err.Error())
	}

	for k, v := range params.Headers {
		request.Header.Set(k, v)
	}

	if opt.AuthInHeader {
		request.Header.Set("Authorization", auth)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	var client *http.Client
	if params.Client != nil {
		client = params.Client
	} else {
		client = &http.Client{}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to %s: %s", params.TokenEndpoint, err.Error())
	}

	defer response.Body.Close()

	resp, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("non-OK response received: %s", resp)
	}

	if response.Header["Content-Type"] == nil {
		return nil, fmt.Errorf("content-type header missing in response: %s", resp)
	}

	contentType := response.Header["Content-Type"][0]

	if strings.Contains(contentType, "application/json") {
		var parsed map[string]interface{}
		if err = json.Unmarshal(resp, &parsed); err != nil {
			return nil, fmt.Errorf("failed to deserialise the response: %s (%s)", resp, err.Error())
		}
		return parsed, nil
	}

	if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "text/html") {

		m, _ := url.ParseQuery(string(resp))
		ret := make(map[string]interface{})
		for k, v := range m {
			ret[k] = v[0]
		}

		return ret, nil
	}

	return nil, fmt.Errorf("could not determine content type encoding from response")
}

// Refresh implements an OAuth2 token refresh methods.  Parameters are sent via the OAuth2 struct
func (params OAuth2) Refresh(opt Options) (map[string]interface{}, error) {

	var auth string

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", params.RefreshToken)

	for k, v := range params.Extras {
		data.Set(k, v)
	}

	if opt.AuthInHeader {
		auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(params.ClientID+":"+params.ClientSecret))
	} else {
		data.Set("client_id", params.ClientID)
		data.Set("client_secret", params.ClientSecret)
	}

	request, err := http.NewRequest("POST", params.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create an HTTP client: %s", err.Error())
	}

	if opt.AuthInHeader {
		request.Header.Set("Authorization", auth)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	var client *http.Client
	if params.Client != nil {
		client = params.Client
	} else {
		client = &http.Client{}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to %s: %s", params.TokenEndpoint, err.Error())
	}

	defer response.Body.Close()

	resp, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("non-OK response received: %s", resp)
	}

	if response.Header["Content-Type"] == nil {
		return nil, fmt.Errorf("content-type header missing in response: %s", resp)
	}

	contentType := response.Header["Content-Type"][0]

	if strings.Contains(contentType, "application/json") {
		var parsed map[string]interface{}
		if err = json.Unmarshal(resp, &parsed); err != nil {
			return nil, fmt.Errorf("failed to deserialise the response: %s (%s)", resp, err.Error())
		}
		return parsed, nil
	}

	if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "text/html") {

		m, _ := url.ParseQuery(string(resp))
		ret := make(map[string]interface{})
		for k, v := range m {
			ret[k] = v[0]
		}

		return ret, nil
	}

	return nil, fmt.Errorf("could not determine content type encoding from response")
}
