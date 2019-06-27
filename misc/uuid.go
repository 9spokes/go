package misc

import (
	"regexp"

	uuid "github.com/satori/go.uuid"
)

// GenUUIDv4 generates a new universal user ID (UUID) version 4
func GenUUIDv4() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return id.String()
}

// IsUUIDv4 determins whether a string is a valid UUID form.  Returns a boolean
func IsUUIDv4(id string) bool {

	match, err := regexp.Match(`^\d{8}-\d{4}-\d{4}-\d{4}-\d{12}$`, []byte(id))

	if !match || err != nil {
		return false
	}

	return true

}
