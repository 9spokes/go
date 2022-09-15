package throttler

import (
	"net"
	"time"
)

type ThrottlerOptions struct {
	MaxWait time.Duration
}

type Request struct {
	Osp    string            `json:"osp"`
	Limits map[string]string `json:"limits"`
	CID    string            `json:"cid"`
}

type Response struct {
	Status  string    `json:"status,omitempty"`
	Message string    `json:"message,omitempty"`
	Retry   time.Time `json:"retry,omitempty"`
	CID     string    `json:"correlationId,omitempty"`
}

type Ticket struct {
	CID  string
	Conn net.Conn
}
