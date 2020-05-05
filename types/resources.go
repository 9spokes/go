package types

import (
	"crypto/x509"

	"github.com/9spokes/go/db"
	event "github.com/9spokes/go/services/events"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
)

// Resources is a struct that contains all the resources that an application could consume
type Resources struct {
	Mongo    db.MongoDB
	RedisDB  *redis.Client
	OSPs     []Document
	Store    *sessions.CookieStore
	Events   *event.Context
	Clients  map[string]string
	Creds    Credentials
	Keystore []x509.Certificate
}
