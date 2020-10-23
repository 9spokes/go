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

// Response is an HTTP response, optionally parsed if the `json=true` parameter is supplied to httpGet()
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
