package recoverer

import (
	"net/http"
	"runtime/debug"

	"github.com/9spokes/go/logging/v3"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil && err != http.ErrAbortHandler {
				if r != nil {
					logging.Errorf("Unhandled panic: %v, \nstack: %s \nrequest: %s %v, \nheaders: %v", err, debug.Stack(), r.Method, r.URL, r.Header)
				} else {
					logging.Errorf("Unhandled panic with nil http request object: %s", debug.Stack())
				}
				if w != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
