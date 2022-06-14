package timer

import (
	"net/http"
	"time"

	"github.com/9spokes/go/logging/v3"
)

// Timer returns a middleware function that checks the time it took to handle
// a request and logs a warning message if it is above the specifed threshhold.
//
// limit is the threshhold value in milliseconds
func Timer(limit int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			next.ServeHTTP(w, r)
			elapsed := time.Since(startTime)

			if elapsed.Milliseconds() > limit {
				logging.Warningf("Slow request [URI: %s, duration: %s]", r.RequestURI, elapsed)
			}
		})
	}
}
