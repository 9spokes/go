package http

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusPartialContent       = 206 // RFC 7233, 4.1
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices  = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently = 301 // RFC 7231, 6.4.2
	StatusFound            = 302 // RFC 7231, 6.4.3
	StatusSeeOther         = 303 // RFC 7231, 6.4.4
	StatusNotModified      = 304 // RFC 7232, 4.1
	StatusUseProxy         = 305 // RFC 7231, 6.4.5

	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
	StatusTeapot                       = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusTooEarly                     = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

// Response is an HTTP response
type Response struct {
	StatusCode int
	Body       []byte
	JSON       interface{}
	Headers    http.Header
}

// Authentication contains the keys needed to build an authorization URL
type Authentication struct {
	Scheme   string
	Username string
	Password string
	Token    string
}

// Request is an HTTP request struct
type Request struct {
	URL            string `default:"http://localhost/"`
	Method         string `default:"GET"`
	Headers        map[string]string
	Query          map[string]string
	ContentType    string `default:"application/json"`
	Authentication Authentication
	Client         *http.Client
	Body           []byte
}

// Post isn an HTTP POST Method
func (request Request) Post() (*Response, error) {
	request.Method = "POST"
	return request.http()
}

// Get is an HTTP GET Method
func (request Request) Get() (*Response, error) {
	request.Method = "GET"
	return request.http()
}

// Put is an HTTP PUT Method
func (request Request) Put() (*Response, error) {
	request.Method = "PUT"
	return request.http()
}

// Patch is an HTTP PATCH Method
func (request Request) Patch() (*Response, error) {
	request.Method = "PATCH"
	return request.http()
}

// Delete is an HTTP DELETE Method
func (request Request) Delete() (*Response, error) {
	request.Method = "DELETE"
	return request.http()
}

func (request Request) http() (*Response, error) {

	var client *http.Client
	if request.Client != nil {
		client = request.Client
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, err
	}

	if request.ContentType == "" {
		req.Header.Set("Content-type", "application/json")
	} else {
		req.Header.Set("Content-type", request.ContentType)
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	switch strings.ToLower(request.Authentication.Scheme) {
	case "basic":
		req.Header.Set("Authorization", request.Authentication.Scheme+" "+base64.StdEncoding.EncodeToString([]byte(request.Authentication.Username+":"+request.Authentication.Password)))
	case "bearer":
		req.Header.Set("Authorization", request.Authentication.Scheme+" "+request.Authentication.Token)
	}

	q := req.URL.Query()
	for k, v := range request.Query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: " + err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{resp.StatusCode, body, "", resp.Header}, fmt.Errorf("failed to read the response body: " + err.Error())
	}

	var decoded interface{}
	json.Unmarshal(body, &decoded)

	if resp.StatusCode > 399 {
		return &Response{resp.StatusCode, body, decoded, resp.Header}, fmt.Errorf("non-OK response: %s", body)
	}

	return &Response{resp.StatusCode, body, decoded, resp.Header}, nil
}

// ValidateBasicAuthCreds parses an HTTP Basic authoriation header and validates the credentials contained therein against
// a the map of credentials supplied as the second argument
// Returns the client ID (username) that matched on success or empty string if no match.  Returns an error if the parsing failed
func ValidateBasicAuthCreds(header string, creds map[string]string) (string, error) {

	if header == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	// Parse Authorization header (eg: Basic a2V5OnZhbHVlCg==)
	fields := strings.Split(header, " ")
	if len(fields) != 2 || strings.ToLower(fields[0]) != "basic" {
		return "", fmt.Errorf("invalid authorization type '%s'", fields[0])
	}

	// Base64 decode
	encoded, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		return "", fmt.Errorf("unable to parse authorization header '%s': %s", fields[0], err.Error())
	}

	// Extract client_id & client_secret
	val := strings.Split(strings.TrimSpace(string(encoded)), ":")
	if len(val) != 2 {
		return "", fmt.Errorf("malformed authorization header: %s", fields[1])
	}

	h := sha256.New()
	h.Write([]byte(val[1]))
	hashed := hex.EncodeToString(h.Sum(nil))

	for id, secret := range creds {
		if strings.ToLower(id) == strings.ToLower(val[0]) && secret == hashed {
			return id, nil
		}
	}

	return "", fmt.Errorf("invalid client_id/client_secret combination")
}
