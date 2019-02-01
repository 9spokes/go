package misc

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"strconv"
	"strings"
)

//GenerateNonce creates a pseudo-random integer to use as a nonce value.
func GenerateNonce() string {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(9999999999))
	if err != nil {
		panic(err)
	}
	nonce := bigInt.Int64()
	return strconv.FormatInt(nonce, 10)
}

//OauthEscape creates a safe text to use in URL's
func OauthEscape(value string) string {
	escapedQuery := url.QueryEscape(value)
	toReturn := strings.Replace(escapedQuery, "+", "%20", -1)
	return toReturn
}
