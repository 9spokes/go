package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
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

// Get is an HTTP GET Method
func Get(request Request) (*Response, error) {

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
