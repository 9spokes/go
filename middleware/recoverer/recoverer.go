package recover

import (
	"net/http"
	"runtime/debug"

	"github.com/9spokes/go/logging/v3"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil && err != http.ErrAbortHandler {
				logging.Errorf("Panic error: %s %s", err, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
