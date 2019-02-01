package crypto

import (
	"sort"
	"strconv"
	"time"

	"github.com/9spokes/go/misc"
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
		"oauth_nonce":            misc.GenerateNonce(),
		"oauth_version":          "1.0",
		"oauth_signature_method": "RSA-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
	}
	METHOD := "GET&"
	if input.Refresh {
		Auth["oauth_session_handle"] = input.SessionHandle
		METHOD = "POST&"
	} else {
		if input.Query != "" {
			Auth["where"] = input.Query
		}
	}
	sortedAuthString := SortAuth(Auth)
	signatureText := METHOD + misc.OauthEscape(input.BaseURL) + "&" + misc.OauthEscape(sortedAuthString)
	signature := SignRSA([]byte(signatureText), input.PrivateKeyPath)
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
		sortedAuthString = sortedAuthString + K + "=" + misc.OauthEscape(Auth[K]) + "&"
	}
	sortedAuthString = sortedAuthString[:len(sortedAuthString)-1]
	return sortedAuthString
}
