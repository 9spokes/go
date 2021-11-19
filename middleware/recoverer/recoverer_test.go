package recoverer

import (
	HTTP "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func Test_Recoverer(t *testing.T) {
	r := mux.NewRouter()
	r.Use(Recoverer)

	r.HandleFunc("/test", func(w HTTP.ResponseWriter, r *HTTP.Request) {
		panic("invalid memory address or nil pointer dereference")
	}).Methods("POST")

	rr := httptest.NewRecorder()
	req, err := HTTP.NewRequest("POST", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(rr, req)

	if rr.Code != HTTP.StatusInternalServerError {
		t.Error("panic error not recovered")
	}
}
