package throttler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	goLogging "github.com/op/go-logging"
)

// Context represents a connection object into the token service
type Context struct {
	URL    string
	Logger *goLogging.Logger
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

	ctx.Logger.Debugf("[%s] Getting rate-limiting token", osp)
	client := &http.Client{}

	url := fmt.Sprintf("%s/token/%s", ctx.URL, osp)
	ctx.Logger.Debugf("[%s] URL is: %s", osp, ctx.URL)

	for i := 0; i < retries; i++ {
		ctx.Logger.Debugf("[%s] Attempt #%d", osp, i+1)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			// ctx.Logger.Errorf("[%s] Error retrieving the list of documents from %s: %s", osp, url, err.Error())
			return fmt.Errorf("error retrieving the list of documents from %s: %s", url, err.Error())
		}

		response, err := client.Do(request)
		if err != nil {
			// ctx.Logger.Errorf("[%s] Error retrieving the list of documents from %s: %s", osp, url, err.Error())
			return fmt.Errorf("error retrieving the list of documents from %s: %s", url, err.Error())
		}

		if response.StatusCode == 200 {
			ctx.Logger.Debugf("[%s] Successfully acquired rate-limiting token", osp)
			return nil
		}

		if response.StatusCode != 429 {
			ctx.Logger.Debugf("[%s] Unexpected response from throttling service: %d", osp, response.StatusCode)
			return fmt.Errorf("unexpected response from throttling service: %d", response.StatusCode)
		}

		time.Sleep(5 * time.Second)
	}

	ctx.Logger.Errorf("[%s] Failed to acquire rate-limtiing token after %d attempts", osp, retries)
	return ErrTooManyRequests
}
