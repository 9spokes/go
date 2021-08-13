package types

import "fmt"

const (
	ErrSeverityFatal string = "fatal"
	ErrSeverityWarn  string = "warn"
	ErrSeverityError string = "error"
)

const (
	ErrDataCodeAccessTokenExpired string = "access token expired"
	ErrDataCodeTooManyRequest     string = "too many request"
	ErrDataCodeOK                 string = "ok"

	ErrDataCodeMissingClientID          string = "missing client id"
	ErrDataCodeMissingContentTypeHeader string = "content-type header missing"
	ErrDataCodeMissingAuthorizationCode string = "authorization code missing"
	ErrDataCodeMissingRefreshToken      string = "refresh token missing"

	ErrDataCodeDeserialiseFailed     string = "deserialise failed"
	ErrDataCodeResponseFormatUnknown string = "response format unknown"
	ErrDataCodeError                 string = "error"
)

type ErrorResponse struct {
	Err      error  `json:"message"`
	Severity string `json:"severity"`
	DataCode string `json:"dataCode"`
	HttpCode int
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("data-code :%s, http-code: %d, err %s", e.DataCode, e.HttpCode, e.Err.Error())
}

func (e *ErrorResponse) Fatal() bool {
	return e.Severity == ErrSeverityFatal
}
