package types

import (
	"github.com/9spokes/go/db"
	"github.com/9spokes/go/event"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"github.com/streadway/amqp"
)

// Resources is a struct that contains all the resources that an application could consume
type Resources struct {
	Mongo   db.MongoDB
	RedisDB *redis.Client
	AMQP    *amqp.Channel
	OSPs    []Document
	Store   *sessions.CookieStore
	Events  *event.Service
	Clients map[string]string
	Creds   Credentials
}
