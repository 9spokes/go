package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Service represents the coordinates of the event service
type Service struct {
	url string
}

// New returns a new event service object
func New(url string) (*Service, error) {

	svc := &Service{
		url: url,
	}

	return svc, nil
}

// Post is a method that writes a new user-based event to a central logging database.  The target collection is determined by the `col` argument.  All other `fields` are written as-is to the doucment.
func (svc *Service) Post(fields map[string]interface{}) error {

	encoded, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("while marshaling fields: %s", err.Error())
	}

	_, err = http.Post(svc.url, "application/json", bytes.NewReader(encoded))
	return err
}
