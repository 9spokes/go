package api

import (
	"encoding/json"
	"net/http"
)

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
