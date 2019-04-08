package misc

import (
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

//GenerateNonce creates a pseudo-random integer to use as a nonce value.
func GenerateNonce() string {
	rand.Seed(time.Now().UnixNano())
	toRet := make([]byte, 32)
	length := len(alphanumeric)
	for i := range toRet {
		toRet[i] = alphanumeric[rand.Int63()%int64(length)]
	}
	return string(toRet)
}

//OauthEscape creates a safe text to use in URL's
func OauthEscape(value string) string {
	escapedQuery := url.QueryEscape(value)
	toReturn := strings.Replace(escapedQuery, "+", "%20", -1)
	return toReturn
}
