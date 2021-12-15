package throttler

import (
	"errors"
	"fmt"
	"time"

	"github.com/9spokes/go/http"

	"github.com/9spokes/go/logging/v3"
)

// Context represents a connection object into the token service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *logging.Logger
}

var (
	ErrTooManyRequests = errors.New("failed to acquire token")
)

// GetToken is a helper function that retrieves a token from the OSP rate limiting service `etl-throttler`.
//
// You must call this function before initiating any outbound connection which enforces a rate limiting policy.
//
// The function will block until a valid token is returned from the throttling service hence you may wish to implement a timeout
func (ctx Context) GetToken(osp string, retries int) error {

	logging.Debugf("[%s] Getting rate-limiting token", osp)

	url := fmt.Sprintf("%s/token/%s", ctx.URL, osp)
	logging.Debugf("[%s] URL is: %s", osp, ctx.URL)

	request := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		ContentType: "application/json",
	}

	for i := 0; i < retries; i++ {
		logging.Debugf("[%s] Attempt #%d", osp, i+1)
		response, _ := request.Get()
		if response.StatusCode == 200 {
			logging.Debugf("[%s] Successfully acquired rate-limiting token", osp)
			return nil
		}

		if response.StatusCode != 429 {
			logging.Debugf("[%s] Unexpected response from throttling service: %d", osp, response.StatusCode)
			return fmt.Errorf("unexpected response from throttling service: %d", response.StatusCode)
		}

		time.Sleep(5 * time.Second)
	}

	logging.Errorf("[%s] Failed to acquire rate-limtiing token after %d attempts", osp, retries)
	return ErrTooManyRequests
}
