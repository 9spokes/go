package network

import (
	"fmt"
	"net"
	"time"
)

// Dial attempts to connect to a TCP or UDP destination until timeout value is reached.
func Dial(proto string, port int, timeout int) error {

	go func() {
		conn, err := net.Dial("tcp", ":1123")
		if err != nil {
			time.Sleep(1000 * time.Second)
		}
		defer conn.Close()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("Timeout connecting to %s/%d after %d seconds", proto, port, timeout)
	default:
		return nil
	}
}
