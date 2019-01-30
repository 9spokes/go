package crypto

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

//Xero is the structure required for input to the xero function
type Xero struct {
	AccessToken    string
	ConsumerKey    string
	TokenSecret    string
	SessionHandle  string
	Refresh        bool
	Query          string
	BaseURL        string
	PrivateKeyPath string
}

//XeroSigner generates a signature for a xero request
func XeroSigner(input Xero) (string, error) {
	Auth := map[string]string{
		"oauth_token":            input.AccessToken,
		"oauth_consumer_key":     input.ConsumerKey,
		"oauth_nonce":            strconv.FormatInt(GenerateNonce(), 10),
		"oauth_version":          "1.0",
		"oauth_signature_method": "RSA-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
	}
	URL := input.BaseURL
	METHOD := "GET&"
	if input.Refresh {
		Auth["oauth_session_handle"] = input.SessionHandle
		METHOD = "POST&"
	} else {
		if input.Query != "" {
			Auth["where"] = input.Query
			URL = URL + "?where=" + OauthEscape(input.Query)
		}
	}
	sortedAuthString := SortAuth(Auth)
	signatureText := METHOD + OauthEscape(URL) + "&" + OauthEscape(sortedAuthString)
	signatureByte := []byte(signatureText)
	signature := GenerateSignature(signatureByte, input.PrivateKeyPath)
	return signature, nil
}

//SortAuth creates a sorted Authentication string. The string is sorted Lexographically.
func SortAuth(Auth map[string]string) string {
	var keys []string
	for K := range Auth {
		keys = append(keys, K)
	}
	sort.Strings(keys)
	sortedAuthString := ""
	for _, K := range keys {
		sortedAuthString = sortedAuthString + K + "=" + OauthEscape(Auth[K]) + "&"
	}
	sortedAuthString = sortedAuthString[:len(sortedAuthString)-1]
	return sortedAuthString
}

//GenerateNonce creates a sudo random integer to use as a nonce value.
func GenerateNonce() int64 {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(9999999999))
	if err != nil {
		panic(err)
	}
	nonce := bigInt.Int64()
	return nonce
}

//OauthEscape creates a safe text to use in URL's
func OauthEscape(value string) string {
	escapedQuery := url.QueryEscape(value)
	toReturn := strings.Replace(escapedQuery, "+", "%20", -1)
	return toReturn
}
