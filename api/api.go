package api

import (
	"encoding/json"
	"net/http"
)

//Response is an API response message envelope
type Response struct {
	Status        string      `json:"status,omitempty"`
	CorrelationId string      `json:"correlationId,omitempty"`
	Message       string      `json:"message,omitempty"`
	Details       interface{} `json:"details,omitempty"`
}

// ErrorResponse sends a standardised error message body to the caller
// in the form of { "status": "err", "message": "<your error message>"}
// code should be a valid 4xx-5xx error status code
func ErrorResponse(w http.ResponseWriter, msg string, code int) {
	ErrorResponseWithCorrelation(w, msg, "", code)
}

// Sends a standardised error message body to the caller
// in the form of { "status": "err", "message": "<your error message>" correlationId: "<corrId>"}.
// code should be a valid 4xx-5xx error status code
func ErrorResponseWithCorrelation(w http.ResponseWriter, msg, corrId string, code int) {

	w.Header().Set("Content-type", "application/json")

	response := Response{Status: "err", Message: msg, CorrelationId: corrId}

	e, _ := json.Marshal(response)
	w.WriteHeader(code)
	w.Write(e)
}

// SuccessResponse sends a standardised success message body to the caller
// in the form of { "status": "ok", "details": "<your data>"}
func SuccessResponse(w http.ResponseWriter, data interface{}, code int) {
	SuccessResponseWithCorrelationAndStatus(w, data, "", code)
}

// SuccessResponse sends a standardised success message body to the caller
// in the form of { "status": "ok", "details": "<your data>"}.
// Http status code will be 200 OK
func SuccessResponseWithCorrelation(w http.ResponseWriter, data interface{}, corrId string) {
	SuccessResponseWithCorrelationAndStatus(w, data, corrId, http.StatusOK)
}

// SuccessResponse sends a standardised success message body to the caller
// in the form of { "status": "ok", "details": "<your data>"}
// statusCode should be a valid 2xx http status
func SuccessResponseWithCorrelationAndStatus(w http.ResponseWriter, data interface{}, corrId string, statusCode int) {

	w.Header().Set("Content-type", "application/json")

	response := Response{Status: "ok", Details: data, CorrelationId: corrId}

	o, _ := json.Marshal(response)
	w.WriteHeader(statusCode)
	w.Write(o)
}
