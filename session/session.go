package session

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Set stores a value into the user session
func Set(redisdb *redis.Client, sessionID, key string, value interface{}) error {

	if _, err := redisdb.HSet(sessionID, key, value).Result(); err != nil {
		return fmt.Errorf("failed to write %s to session: %s", key, err.Error())
	}

	return nil
}

// Get reads a value from the user session
func Get(redisdb *redis.Client, sessionID, key string) (string, error) {

	value, err := redisdb.HGet(sessionID, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to read %s from session: %s", key, err.Error())
	}

	return value, nil
}
