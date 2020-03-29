package events

import (
	"encoding/json"
	"fmt"

	"github.com/9spokes/go/http"
	"github.com/9spokes/go/types"
	goLogging "github.com/op/go-logging"
)

// Context represents a connection object into the token service
type Context struct {
	URL          string
	ClientID     string
	ClientSecret string
	Logger       *goLogging.Logger
}

// Post is a method that writes a new user-based event to a central logging database.  The target collection is determined by the `col` argument.  All other `fields` are written as-is to the doucment.
func (ctx Context) Post(fields types.Document) error {

	encoded, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("while marshaling fields: %s", err.Error())
	}

	_, err = http.Request{
		URL:         ctx.URL,
		ContentType: "application/json",
		Body:        encoded,
	}.Post()

	return err
}
