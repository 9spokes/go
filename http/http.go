package http

import (
	"bytes"
	"encoding/json"
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

// ErrorResponse sends a standardised error message body to the caller
// in the form of { "status": "err", "message": "<your error message>"}
func ErrorResponse(w http.ResponseWriter, msg string, code int) {

	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("AAccess-Control-Request-Headers", "*")

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{"err", msg}

	e, _ := json.Marshal(response)
	w.WriteHeader(code)
	w.Write(e)
	return
}

// SuccessResponse sends a standardised success message body to the caller
// in the form of { "status": "ok", "details": "<your data>"}
func SuccessResponse(w http.ResponseWriter, data interface{}, code int) {

	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Request-Headers", "*")
	w.Header().Set("Access-Control-Allow-Method", "*")

	response := struct {
		Status  string      `json:"status"`
		Details interface{} `json:"details"`
	}{"ok", data}

	o, _ := json.Marshal(response)
	w.Write(o)
	return
}
