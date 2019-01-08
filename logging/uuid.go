package logging

import uuid "github.com/satori/go.uuid"

// GenUUIDv4 generates a new universal user ID (UUID) version 4
func GenUUIDv4() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return id.String()
}
