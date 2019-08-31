package api

import (
	"encoding/json"
	"net/http"
)

//Response is an API response message envelope
type Response struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorResponse sends a standardised error message body to the caller
// in the form of { "status": "err", "message": "<your error message>"}
func ErrorResponse(w http.ResponseWriter, msg string, code int) {

	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("AAccess-Control-Request-Headers", "*")

	response := Response{Status: "err", Message: msg}

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

	response := Response{Status: "ok", Details: data, Message: ""}

	o, _ := json.Marshal(response)
	w.Write(o)
	return
}
