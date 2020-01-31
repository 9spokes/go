package status

import (
	"github.com/9spokes/go/db"
)

//ValidateMongo validates an Mongo connection
func ValidateMongo(mongo db.MongoDB) bool {
	return mongo.Session.Ping() == nil
}
