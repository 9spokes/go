package status

import (
	"github.com/go-redis/redis"
)

//ValidateRedis validates an Redis connection
func ValidateRedis(r *redis.Client) bool {
	_, err := r.Ping().Result()
	if err != nil {
		return false
	}
	return true

}
