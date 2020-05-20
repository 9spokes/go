package misc

import (
	"regexp"

	"github.com/satori/go.uuid"
)

// GenUUIDv4 generates a new universal user ID (UUID) version 4
func GenUUIDv4() string {
	id := uuid.NewV4()
	return id.String()
}

// IsUUIDv4 determins whether a string is a valid UUID form.  Returns a boolean
func IsUUIDv4(id string) bool {

	if match, err := regexp.Match(`^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$`, []byte(id)); !match || err != nil {
		return false
	}

	return true

}
