package types

import "fmt"

const (
	ErrSeverityFatal string = "fatal"
	ErrSeverityWarn  string = "warn"
	ErrSeverityError string = "error"
)

const (
	ErrAccessTokenExpired string = "access token expired"
	ErrTooManyRequests    string = "rate limit threshold exceeded"
	ErrOK                 string = "ok"

	ErrMissingClientID          string = "missing client id"
	ErrMissingContentTypeHeader string = "content-type header missing"
	ErrMissingAuthorizationCode string = "authorization code missing"
	ErrMissingRefreshToken      string = "refresh token missing"

	ErrDeserialiseFailed     string = "deserialise failed"
	ErrResponseFormatUnknown string = "response format unknown"
	ErrError                 string = "error"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	Severity   string `json:"severity"`
	ID         string `json:"id"`
	HTTPStatus int
}

func (e *ErrorResponse) Error() string {
	m := fmt.Sprintf("%s: %s", e.Message, e.ID)
	if e.HTTPStatus > 0 {
		return fmt.Sprintf("[HTTP %d] %s", e.HTTPStatus, m)
	}
	return m
}

func (e *ErrorResponse) IsFatal() bool {
	return e.Severity == ErrSeverityFatal
}
