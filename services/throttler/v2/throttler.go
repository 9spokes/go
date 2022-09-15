package throttler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// Context represents a handle into the throttler service
type Context struct {
	URL string
}

func connect(url string) (net.Conn, error) {
	for {
		c, err := net.Dial("tcp", url)
		if err != nil {
			if strings.Contains(err.Error(), "too many open files") {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return nil, err
		}
		return c, err
	}
}

func get(c net.Conn, req Request) (*Response, error) {
	marshalled, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintln(c, string(marshalled)); err != nil {
		return nil, err
	}

	data, err := bufio.NewReader(c).ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	resp := &Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Status == "err" {
		if resp.Retry.IsZero() {
			resp.Retry = time.Now().Add(time.Second)
		}

		return resp, fmt.Errorf(resp.Message)
	}

	return resp, nil
}

// Asks the Throttler Svc for permission to call an external API that is
// protected by the rate limits specified in the request. It returns a
// connection which should be closed by the caller after the ticket is used.
// Closing the connection notifies the Throttler Svc that it can recycle the
// ticket.
//
// The following call requests a ticket for BAC and is willing to wait up to
// 5 minutes for that ticket in case one is not available right away.
//
//	ctx.GetTicket(Request{
//		Osp: "bac",
//		CID: "7e3af66b-51ce-4012-8c0a-57827c886981",
//	}, ThrottlerOptions{MaxWait: time.Minute * 5})
//
// The following call requests a ticket for ZohoBooks specifying the rate limits
// which should be observed. It also says that the caller is willing to wait up
// to 2 mintues for that ticket.
//
//	ctx.GetTicket(Request{
//		Osp: "zohobooks",
//		CID: "051a1952-a10d-4ade-9e2f-92d8f2c4f390",
//		Limits: map[string]string{
//			"views-per-day": "051a1952-a10d-4ade-9e2f-92d8f2c4f390",
//			"views-per-min": "132a531c-bf79-42ec-8a6b-a08c684a8e44",
//	}}, ThrottlerOptions{MaxWait: time.Minute * 2})
func (ctx Context) GetTicket(req Request, opt ThrottlerOptions) (net.Conn, error) {

	for deadline := time.Now().Add(opt.MaxWait); time.Until(deadline) >= 0; {
		c, err := connect(ctx.URL)
		if err != nil {
			return nil, err
		}

		resp, err := get(c, req)
		if err != nil {
			if resp.Retry.Before(deadline) {
				time.Sleep(time.Until(resp.Retry))
				continue
			}

			c.Close()
			return c, fmt.Errorf(resp.Message)
		}

		return c, nil
	}

	return nil, fmt.Errorf("reached deadline")
}
