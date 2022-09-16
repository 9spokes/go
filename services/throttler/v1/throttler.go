package throttler

import (
	"errors"
	"fmt"
	"time"

	"github.com/9spokes/go/http"

	"github.com/9spokes/go/logging/v3"
)

// Context represents a connection object into the Throttler service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
}

type ThrottlerOptions struct {
	Retries       int
	CorrelationID string
}

var (
	ErrTooManyRequests = errors.New("failed to acquire token")
)

// GetToken is a helper function that retrieves a token from the OSP rate limiting service `etl-throttler`.
//
// You must call this function before initiating any outbound connection which enforces a rate limiting policy.
//
// The function will block until a valid token is returned from the throttling service hence you may wish to implement a timeout
func (ctx Context) GetToken(osp string, opt ThrottlerOptions) error {

	logging.Debugf("[CorrID:%s][%s] Getting rate-limiting token", opt.CorrelationID, osp)

	url := fmt.Sprintf("%s/token/%s", ctx.URL, osp)
	logging.Debugf("[CorrID:%s][%s] URL is: %s", opt.CorrelationID, osp, ctx.URL)

	request := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"x-correlation-id": opt.CorrelationID,
		},
		ContentType: "application/json",
	}

	for i := 0; i < opt.Retries; i++ {
		logging.Debugf("[CorrID:%s][%s] Attempt #%d", opt.CorrelationID, osp, i+1)
		response, err := request.Get()
		if err != nil {
			return fmt.Errorf("failed to issue request: %w", err)
		}

		if response.StatusCode == 200 {
			logging.Debugf("[CorrID:%s][%s] Successfully acquired rate-limiting token", opt.CorrelationID, osp)
			return nil
		}

		if response.StatusCode != 429 {
			logging.Debugf("[CorrID:%s][%s] Unexpected response from throttling service: %d", opt.CorrelationID, osp, response.StatusCode)
			return fmt.Errorf("unexpected response from throttling service: %d", response.StatusCode)
		}

		time.Sleep(5 * time.Second)
	}

	logging.Errorf("[CorrID:%s][%s] Failed to acquire rate-limtiing token after %d attempts", opt.CorrelationID, osp, opt.Retries)
	return ErrTooManyRequests
}
