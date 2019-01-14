package network

import (
	"fmt"
	"net"
	"time"
)

// Dial attempts to connect to a TCP or UDP destination until timeout value is reached.
func Dial(proto string, dest string, timeout int) error {

	ch := make(chan string, 1)
	go func() {
		for {
			_, err := net.Dial(proto, dest)
			if err == nil {
				ch <- "connected"
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("Timeout connecting to %s/%s after %d seconds", proto, dest, timeout)
	case <-ch:
		return nil
	}
}
