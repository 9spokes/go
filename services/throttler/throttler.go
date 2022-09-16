package throttler

import (
	"net/url"
	"time"

	v1 "github.com/9spokes/go/services/throttler/v1"
	v2 "github.com/9spokes/go/services/throttler/v2"
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
	MaxWait       time.Duration
}

var ErrTooManyRequests = v1.ErrTooManyRequests

func (ctx Context) GetToken(osp string, opt ThrottlerOptions) (*v2.Ticket, error) {

	if ctx.URL == "http://throttlerng" {
		u, err := url.Parse(osp)
		if err != nil {
			return nil, err
		}

		limits := map[string]string{}
		q := u.Query()
		for k := range q {
			limits[k] = q.Get(k)
		}

		return v2.Context{
			URL: ctx.URL,
		}.GetTicket(v2.Request{
			Osp:    u.Path,
			Limits: limits,
			CID:    opt.CorrelationID,
		}, v2.ThrottlerOptions{MaxWait: opt.MaxWait})

	} else {

		err := v1.Context{
			URL:          ctx.URL,
			ClientID:     ctx.ClientID,
			ClientSecret: ctx.ClientSecret,
		}.GetToken(osp, v1.ThrottlerOptions{
			Retries:       opt.Retries,
			CorrelationID: opt.CorrelationID,
		})

		return &v2.Ticket{Conn: nil}, err
	}
}
