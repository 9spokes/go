package events

import (
	"encoding/json"
	"fmt"

	"github.com/9spokes/go/http"
)

// Context represents the coordinates of the event service
type Context struct {
	url string
}

// New returns a new event context object
func New(url string) (*Context, error) {

	svc := &Context{
		url: url + "/events",
	}

	return svc, nil
}

// Post is a method that writes a new user-based event to a central logging database.
func (svc *Context) Post(event Event) error {

	encoded, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("while marshaling fields: %s", err.Error())
	}

	_, err = http.Request{
		Body:        encoded,
		ContentType: "application/json",
		URL:         svc.url,
	}.Post()

	return err
}
