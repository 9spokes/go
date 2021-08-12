package types

var (
	SeverityFatal string = "fatal"
	SeverityWarn  string = "warn"
	SeverityError string = "error"
)

var (
	DataCodeAccessTokenExpired string = "access token expired"
	DataCodeTooManyRequest     string = "too many request"
	DataCodeOK                 string = "ok"

	DataCodeMissingClientID          string = "missing client id"
	DataCodeMissingContentTypeHeader string = "content-type header missing"
	DataCodeMissingAuthorizationCode string = "authorization code missing"
	DataCodeMissingRefreshToken      string = "refresh token missing"

	DataCodeDeserialiseFailed     string = "deserialise failed"
	DataCodeResponseFormatUnknown string = "response format unknown"
	DataCodeError                 string = "error"
)

type ErrorResponse struct {
	Err      error  `json:"message"`
	Severity string `json:"severity"`
	DataCode string `json:"dataCode"`
	HttpCode int
}

func (e *ErrorResponse) Error() string {
	return e.Err.Error()
}

func (e *ErrorResponse) Fatal() bool {
	return e.Severity == SeverityFatal
}
