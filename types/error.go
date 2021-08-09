package types

type ErrorResponse struct {
	Error error
	Fatal bool
	Code  int
}
