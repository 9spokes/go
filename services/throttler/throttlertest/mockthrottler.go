package throttlertest

import (
	"fmt"

	"github.com/9spokes/go/services/throttler"
)

// Mock Throttler service that always returns success
type MockThrottlerSuccess struct{}

func (ctx MockThrottlerSuccess) GetTicket(req throttler.Request, opt throttler.ThrottlerOptions) (*throttler.Ticket, error) {
	return &throttler.Ticket{Conn: nil}, nil
}

// Mock Throttler service that always returns error
type MockThrottlerErr struct{}

func (ctx MockThrottlerErr) GetTicket(req throttler.Request, opt throttler.ThrottlerOptions) (*throttler.Ticket, error) {
	return nil, fmt.Errorf("no ticket available")
}
