package recoverer

import (
	"fmt"
	HTTP "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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

func TestGoroutineRecoverer(t *testing.T) {

	tests := []struct {
		Name        string
		Fn          func(func(interface{}), func())
		CallCleanup bool
		CallFinally bool
	}{
		{
			Name: "log only",
			Fn: func(cleanup func(interface{}), finally func()) {
				defer RecoverGoroutinePanic("log only", nil, nil)
				panic("test panic")
			},
			CallCleanup: false,
			CallFinally: false,
		},
		{
			Name: "cleanup only",
			Fn: func(cleanup func(interface{}), finally func()) {
				defer RecoverGoroutinePanic("cleanup only", cleanup, nil)
				panic("test panic")
			},
			CallCleanup: true,
			CallFinally: false,
		},
		{
			Name: "finally only",
			Fn: func(cleanup func(interface{}), finally func()) {
				defer RecoverGoroutinePanic("finally only", nil, finally)
				panic("test panic")
			},
			CallCleanup: false,
			CallFinally: true,
		},
		{
			Name: "cleanup and finally",
			Fn: func(cleanup func(interface{}), finally func()) {
				defer RecoverGoroutinePanic("cleanup and finally", cleanup, finally)
				panic("test panic")
			},
			CallCleanup: true,
			CallFinally: true,
		},
		{
			Name: "without panic",
			Fn: func(cleanup func(interface{}), finally func()) {
				defer RecoverGoroutinePanic("without panic", cleanup, finally)
				fmt.Println("run without panic")
			},
			CallCleanup: false,
			CallFinally: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)

			cleanCount := 0
			cleanup := func(interface{}) {
				cleanCount++
			}

			finallyCount := 0
			finally := func() {
				finallyCount++
			}
			defer func() {
				if err := recover(); err != nil {
					assert.FailNow("panic not recovered", "err: %v", err)
				}
			}()

			test.Fn(cleanup, finally)

			if test.CallCleanup {
				assert.Equal(1, cleanCount, "clean should be called once")
			} else {
				assert.Equal(0, cleanCount, "clean should not be called")
			}

			if test.CallFinally {
				assert.Equal(1, finallyCount, "finally should be called once")
			} else {
				assert.Equal(0, finallyCount, "finally should not be called")
			}

		})
	}
}
