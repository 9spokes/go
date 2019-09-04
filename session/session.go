package session

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

//Validate is used to ensure a bearer token represents a valid session in the cache.  If so, the user ID is retrieved and returned to the caller
func Validate(redisdb *redis.Client, auth string) (string, error) {

	// Extract the JWT token from the Authorization header
	tokenStr, err := parseAuthHeader(auth)
	if err != nil {
		return "", fmt.Errorf("while parsing authorization header: %s", err.Error())
	}

	// Validate token and extract the subject
	sub, err := validateToken(tokenStr)
	if err != nil {
		return "", fmt.Errorf("while validating JWT token: %s", err.Error())
	}

	// Lookup the session in Redis
	user, err := Get(redisdb, sub, "remote")
	if err != nil {
		return "", fmt.Errorf("while retrieving user ID from session: %s", err.Error())
	}

	return user, nil

}

func validateToken(token string) (string, error) {
	type claims struct {
		jwt.StandardClaims
	}

	// Validate that the JWT is valid
	ret, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "", fmt.Errorf("Malformed JWT token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return "", fmt.Errorf("Expired token")
		}
	} else if !ret.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	return ret.Claims.(*claims).Subject, nil
}

func parseAuthHeader(header string) (string, error) {
	// Check that we have an Authorization header
	if header == "" {
		return "", fmt.Errorf("Missing authorization header")
	}

	// Item to tokenize the header and check that the type is "Bearer"
	items := strings.Fields(header)
	if strings.ToLower(items[0]) != "bearer" {
		return "", fmt.Errorf("Invalid authorization header type")
	}

	// Verify that the 2nd token (the JWT token) is present
	if len(items) < 2 {
		return "", fmt.Errorf("Invalid authorization header")
	}

	return items[1], nil
}
