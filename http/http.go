package http

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Response is an HTTP response, optionally parsed if the `json=true` parameter is supplied to httpGet()
type Response struct {
	StatusCode int
	Body       []byte
}

// Request is an HTTP request struct
type Request struct {
	URL     string `default:"http://localhost/"`
	Method  string `default:"GET"`
	Headers map[string]string
	Query   map[string]string
	Body    string `default:""`
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

	client := &http.Client{}

	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer([]byte(request.Body)))

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	q := req.URL.Query()
	for k, v := range request.Query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.New("Failed to read the response body: " + err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{resp.StatusCode, body}, errors.New("Failed to read the response body: " + err.Error())
	}

	return &Response{resp.StatusCode, body}, nil
}

// ValidateBasicAuthCreds parses an HTTP Basic authoriation header and validates the credentials contained therein against
// a the map of credentials supplied as the second argument
// Returns the client ID (username) that matched on success or empty string if no match.  Returns an error if the parsing failed
func ValidateBasicAuthCreds(header string, creds map[string]string) (string, error) {

	if header == "" {
		return "", fmt.Errorf("Authorization header missing")
	}

	// Parse Authorization header (eg: Basic a2V5OnZhbHVlCg==)
	fields := strings.Split(header, " ")
	if len(fields) != 2 || strings.ToLower(fields[0]) != "basic" {
		return "", fmt.Errorf("Invalid authorization type '%s'", fields[0])
	}

	// Base64 decode
	encoded, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		return "", fmt.Errorf("Unable to parse authorization header '%s': %s", fields[0], err.Error())
	}

	// Extract client_id & client_secret
	val := strings.Split(strings.TrimSpace(string(encoded)), ":")
	if len(val) != 2 {
		return "", fmt.Errorf("Malformed authorization header: %s", fields[1])
	}

	for id, secret := range creds {
		if strings.ToLower(id) == strings.ToLower(val[0]) && secret == val[1] {
			return id, nil

		}
	}

	return "", fmt.Errorf("Invalid client_id/client_secret combination")
}
