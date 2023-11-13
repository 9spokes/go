package throttler

import (
	"fmt"
	"regexp"
	"testing"
)

func TestParseUrl(t *testing.T) {
	cases := []struct {
		description string
		url         string
		e           string
		scheme      string
		host        string
		port        string
	}{
		{
			description: "a URL with scheme and default port 80",
			url:         "tcp://throttlerng",
			scheme:      "tcp",
			host:        "throttlerng:80",
			port:        "80",
		},
		{
			description: "a URL without scheme",
			url:         "throttlerng:80",
			scheme:      "tcp",
			host:        "throttlerng:80",
			port:        "80",
		},
		{
			description: "an invalid URL",
			url:         " tcp://throttlerng",
			e:           "parse",
		},
	}

	for _, c := range cases {
		fmt.Printf("Testing case: %s\n", c.description)

		ctx, err := New(c.url)
		if err != nil && c.e != "" && !regexp.MustCompile(c.e).MatchString(err.Error()) {
			t.Fatalf("Test case failed. Expected error to contain '%s', got '%s'", c.e, err.Error())
		}
		if ctx != nil {
			if ctx.url.Scheme != c.scheme {
				t.Fatalf("Test case failed. Expected scheme '%s', got '%s'", c.scheme, ctx.url.Scheme)
			}
			if ctx.url.Host != c.host {
				t.Fatalf("Test case failed. Expected host '%s', got '%s'", c.host, ctx.url.Host)
			}
			if ctx.url.Port() != c.port {
				t.Fatalf("Test case failed. Expected port '%s', got '%s'", c.port, ctx.url.Port())
			}
		}
	}
}
