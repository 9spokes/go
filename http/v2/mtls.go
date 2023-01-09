package http

import (
	"net/http"

	v1 "github.com/9spokes/go/http"
)

type MTLSRequest v1.MTLSRequest

// Creates a new secure HTTP client.
func (req MTLSRequest) New() (*http.Client, error) {
	return v1.MTLSRequest(req).New()
}
