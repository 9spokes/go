package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// A middleware is a function that receives the Request as input and returns
// the Response. It can modify the request before passing control to the
// next middleware in the call chain. It can also modify the response before
// returning it.
type Middleware func(context.Context, *Request) (*Response, error)

// A middleware function is a function that receives a middleware instance and
// produces a middleware instance. The input is the next middlware in the
// execution chain.
type MiddlewareFunc func(Middleware) Middleware

// Can be used to specify middlewares for the request. These blocks of code
// will be executed in the order in which they are added.
func (r *Request) Use(m MiddlewareFunc) {
	r.Middlewares = append(r.Middlewares, m)
}

// Response is a wrapper around http.Response and besides all the attributes of
// an http.Response it contains the response payload as a byte slice.
type Response struct {
	*http.Response
	Payload []byte
	Date    time.Time
}

// Authorization contains the keys needed to build an authorization header.
type Authorization struct {
	Scheme   string
	Username string
	Password string
	Token    string
}

// Request represents an HTTP request and is the starting point for any HTTP
// network call. It's attributes can be grouped into 3 categries:
//
//   - core: URL, Method, Headers, Query and Body
//   - security: Authorization and Client
//   - middleware: Middlewares and Context
//
// The name of the attributes should be self descriptive but there are some
// mentions to be made.
//
// `Method` should not be set directly. It will be set by the HTTP verb
// methods: Get, Post, etc.
//
// `Authorization` is the preferred way for setting the request's Authorization
// header and if the authorization details are specified both using the Headers
// and Authorization attributes, the latter will take precedence.
//
// `Middlewares` and `Use` can be used to specify middlewares for the request.
// These blocks of code will be executed in the order in which they are added.
//
// `Context` can be used to store additional information about the context in
// which the request is being made. It is useful for supplying extra data to
// the middlewares.
type Request struct {
	URL     string
	Method  string
	Headers map[string]string
	Query   map[string]string
	Body    []byte

	Authorization Authorization
	Client        *http.Client

	Middlewares []MiddlewareFunc
}

// Send an HTTP POST request. The context can be used to set a timeout or
// cancel the request.
func (request *Request) Post(ctx context.Context) (*Response, error) {
	request.Method = "POST"
	return request.httpWithMiddleware(ctx)
}

// Send an HTTP GET request. The context can be used to set a timeout or
// cancel the request.
func (request *Request) Get(ctx context.Context) (*Response, error) {
	request.Method = "GET"
	return request.httpWithMiddleware(ctx)
}

// Send an HTTP PUT request. The context can be used to set a timeout or
// cancel the request.
func (request *Request) Put(ctx context.Context) (*Response, error) {
	request.Method = "PUT"
	return request.httpWithMiddleware(ctx)
}

// Send an HTTP PATCH request. The context can be used to set a timeout or
// cancel the request.
func (request *Request) Patch(ctx context.Context) (*Response, error) {
	request.Method = "PATCH"
	return request.httpWithMiddleware(ctx)
}

// Send an HTTP DELETE request. The context can be used to set a timeout or
// cancel the request.
func (request *Request) Delete(ctx context.Context) (*Response, error) {
	request.Method = "DELETE"
	return request.httpWithMiddleware(ctx)
}

// Wraps the request making `http` method with the available middlewares and
// starts the execution chain.
func (request *Request) httpWithMiddleware(ctx context.Context) (*Response, error) {
	var next Middleware = func(context.Context, *Request) (*Response, error) {
		return request.http(ctx)
	}

	for i := len(request.Middlewares) - 1; i >= 0; i-- {
		next = request.Middlewares[i](next)
	}

	return next(ctx, request)
}

// Makes an HTTP request and, if successful, reads the response body.
func (request *Request) http(ctx context.Context) (*Response, error) {

	// Use default HTTP client if not set.
	if request.Client == nil {
		request.Client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	// Set authorization
	switch strings.ToLower(request.Authorization.Scheme) {
	case "basic":
		req.SetBasicAuth(request.Authorization.Username, request.Authorization.Password)
	case "bearer":
		req.Header.Set("Authorization", request.Authorization.Scheme+" "+request.Authorization.Token)
	}

	// Set query params
	req.URL.RawQuery = request.QueryString()

	// Call endpoint
	resp, err := request.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed while calling %v: %w ", req, err)
	}
	defer resp.Body.Close()

	// Read payload
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed while reading response from %v: %w ", req, err)
	}

	return &Response{resp, payload, time.Now().UTC()}, nil
}

// Returns the details of the HTTP request as a string. Useful for debugging or
// logging errors.
func (r *Request) String() string {
	return fmt.Sprintf("URL: %s (Params: %v)", r.URL, r.QueryString())
}

// Returns query string
func (r *Request) QueryString() string {
	q := url.Values{}
	for k, v := range r.Query {
		q.Add(k, v)
	}
	return q.Encode()
}
