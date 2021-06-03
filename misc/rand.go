package misc

import (
	"encoding/base32"
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers = "1234567890"
const alphanumeric = letters + numbers

func init() {
	rand.Seed(time.Now().UnixNano())
}

//GenerateNonce creates a pseudo-random alpanumeric code to use as a nonce value.
func GenerateNonce() string {
	return GenerateCode(32, "alphanumeric", "none")
}

// GenerateCode a random code having the specified lenght. Supported code types are:
// "numeric" for numeric only codes, "alpha" for alphabetic only codes and
// "alphanumeric" for alphanumeric codes. Supported encodings are "base32" and "none"
func GenerateCode(length uint, codeType string, encoding string) string {
	var runes = []rune(alphanumeric)
	switch codeType {
	case "numeric":
		runes = []rune(numbers)
	case "alpha":
		runes = []rune(letters)
	}

	b := make([]rune, length)
	len := len(runes)
	for i := range b {
		b[i] = runes[rand.Int63()%int64(len)]
	}

	switch encoding {
	case "base32":
		return base32.StdEncoding.EncodeToString([]byte(string(b)))
	default:
		return string(b)
	}
}
