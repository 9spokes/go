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
	URL           string
	ClientID      string
	ClientSecret  string
	CorrelationID string
}

var (
	ErrTooManyRequests = errors.New("failed to acquire token")
)

func (ctx Context) WithCID(correlationID string) Context {
	return Context{
		URL:           ctx.URL,
		ClientID:      ctx.ClientID,
		ClientSecret:  ctx.ClientSecret,
		CorrelationID: correlationID,
	}
}

// GetToken is a helper function that retrieves a token from the OSP rate limiting service `etl-throttler`.
//
// You must call this function before initiating any outbound connection which enforces a rate limiting policy.
//
// The function will block until a valid token is returned from the throttling service hence you may wish to implement a timeout
func (ctx Context) GetToken(osp string, retries int) error {

	logging.Debugf("[CID: %s][%s] Getting rate-limiting token", ctx.CorrelationID, osp)

	url := fmt.Sprintf("%s/token/%s", ctx.URL, osp)
	logging.Debugf("[CID: %s][%s] URL is: %s", ctx.CorrelationID, osp, ctx.URL)

	request := http.Request{
		URL: url,
		Authentication: http.Authentication{
			Scheme:   "basic",
			Username: ctx.ClientID,
			Password: ctx.ClientSecret,
		},
		Headers: map[string]string{
			"x-correlation-id": ctx.CorrelationID,
		},
		ContentType: "application/json",
	}

	for i := 0; i < retries; i++ {
		logging.Debugf("[CID: %s][%s] Attempt #%d", ctx.CorrelationID, osp, i+1)
		response, err := request.Get()
		if err != nil {
			return fmt.Errorf("failed to issue request: %w", err)
		}

		if response.StatusCode == 200 {
			logging.Debugf("[CID: %s][%s] Successfully acquired rate-limiting token", ctx.CorrelationID, osp)
			return nil
		}

		if response.StatusCode != 429 {
			logging.Debugf("[CID: %s][%s] Unexpected response from throttling service: %d", ctx.CorrelationID, osp, response.StatusCode)
			return fmt.Errorf("unexpected response from throttling service: %d", response.StatusCode)
		}

		time.Sleep(5 * time.Second)
	}

	logging.Errorf("[CID: %s][%s] Failed to acquire rate-limtiing token after %d attempts", ctx.CorrelationID, osp, retries)
	return ErrTooManyRequests
}
