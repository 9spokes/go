package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	Http "github.com/9spokes/go/http"
	"github.com/9spokes/go/services/throttler"
	"github.com/9spokes/go/types"
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
	DataInQuery  bool `default:"false"`
}

func (params OAuth2) oauthRequest(opt Options, data url.Values) (map[string]interface{}, error) {
	if params.ClientID == "" {
		return nil, &types.ErrorResponse{Severity: types.SeverityFatal, DataCode: types.DataCodeMissingClientID, Err: fmt.Errorf("client_id cannot be empty")}
	}

	var auth Http.Authentication
	if opt.AuthInHeader {
		auth = Http.Authentication{
			Scheme:   "Basic",
			Username: params.ClientID,
			Password: params.ClientSecret,
		}
	} else {
		data.Set("client_id", params.ClientID)
		if !opt.DataInQuery && params.ClientSecret != "" {
			data.Set("client_secret", params.ClientSecret)
		}
	}

	var body string
	requestUrl := params.TokenEndpoint

	if opt.DataInQuery {
		u, err := url.Parse(params.TokenEndpoint)
		if err != nil {
			return nil, &types.ErrorResponse{Err: err, DataCode: types.DataCodeError}
		}
		rawQuery := u.RawQuery
		if rawQuery != "" {
			rawQuery += "&"
		}
		rawQuery += data.Encode()
		u.RawQuery = rawQuery
		requestUrl = u.String()
	} else {
		body = data.Encode()
	}

	request := Http.Request{
		Client:         params.Client,
		URL:            requestUrl,
		Body:           []byte(body),
		Authentication: auth,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	if !opt.DataInQuery {
		request.Headers["Content-type"] = "application/x-www-form-urlencoded"
	}

	for k, v := range params.Headers {
		request.Headers[k] = v
	}

	response, err := request.Post()

	if err != nil {
		return nil, &types.ErrorResponse{HttpCode: response.StatusCode, DataCode: types.DataCodeError, Err: fmt.Errorf("error while connecting to %s: %w", params.TokenEndpoint, err)}
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return nil, &types.ErrorResponse{HttpCode: response.StatusCode, DataCode: types.DataCodeTooManyRequest, Err: throttler.ErrTooManyRequests}
	}

	if response.Headers["Content-Type"] == nil {
		return nil, &types.ErrorResponse{HttpCode: response.StatusCode, DataCode: types.DataCodeMissingContentTypeHeader, Err: fmt.Errorf("content-type header missing in response: %s", response.Body)}
	}

	contentType := response.Headers["Content-Type"][0]

	if strings.Contains(contentType, "application/json") {
		parsed, ok := response.JSON.(map[string]interface{})
		if !ok {
			return nil, &types.ErrorResponse{HttpCode: response.StatusCode, DataCode: types.DataCodeDeserialiseFailed, Err: fmt.Errorf("failed to deserialise the response: %s", response.Body)}
		}
		return parsed, nil
	}

	if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "text/html") {

		m, _ := url.ParseQuery(string(response.Body))
		ret := make(map[string]interface{})
		for k, v := range m {
			ret[k] = v[0]
		}

		return ret, nil
	}

	return nil, &types.ErrorResponse{HttpCode: response.StatusCode, DataCode: types.DataCodeResponseFormatUnknown, Err: fmt.Errorf("could not determine content type encoding from response")}
}

// Authorize implements an OAuth2 authorization using the parameters defined in the OAuth2 struct
func (params OAuth2) Authorize(opt Options) (map[string]interface{}, error) {

	if params.Code == "" {
		return nil, &types.ErrorResponse{Severity: types.SeverityFatal, DataCode: types.DataCodeMissingAuthorizationCode, Err: fmt.Errorf("the authorization code is missing")}
	}

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {params.Code},
		"redirect_uri": {params.RedirectURI},
	}

	for k, v := range params.Extras {
		data.Set(k, v)
	}

	return params.oauthRequest(opt, data)
}

// Refresh implements an OAuth2 token refresh methods.  Parameters are sent via the OAuth2 struct
func (params OAuth2) Refresh(opt Options) (map[string]interface{}, error) {

	if params.RefreshToken == "" {
		return nil, &types.ErrorResponse{Severity: types.SeverityFatal, DataCode: types.DataCodeMissingRefreshToken, Err: fmt.Errorf("the refresh token is missing")}
	}

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {params.RefreshToken},
	}

	return params.oauthRequest(opt, data)
}
