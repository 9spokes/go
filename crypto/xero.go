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
	SessionHandle  string
	Refresh        bool
	Query          string
	BaseURL        string
	PrivateKeyPath string
}

//Sign generates a signature for a xero request
func Sign(input Xero) (string, error) {
	auth := map[string]string{
		"oauth_token":            input.AccessToken,
		"oauth_consumer_key":     input.ConsumerKey,
		"oauth_nonce":            misc.GenerateNonce(),
		"oauth_version":          "1.0",
		"oauth_signature_method": "RSA-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
	}
	METHOD := "GET&"
	if input.Refresh {
		auth["oauth_session_handle"] = input.SessionHandle
		METHOD = "POST&"
	} else {
		if input.Query != "" {
			auth["where"] = input.Query
		}
	}
	sortedAuthString := sortAuth(auth)
	signatureText := METHOD + misc.OauthEscape(input.BaseURL) + "&" + misc.OauthEscape(sortedAuthString)
	signature := SignRSA([]byte(signatureText), input.PrivateKeyPath)
	var authHeader string
	if input.Refresh {
		authHeader = "OAuth oauth_consumer_key=\\\"" + auth["oauth_consumer_key"] + "\\\",oauth_nonce=\\\"" + auth["oauth_nonce"] + "\\\",oauth_session_handle=\\\"" + auth["oauth_session_handle"] + "\\\",oauth_signature_method=\\\"" + auth["oauth_signature_method"] + "\\\",oauth_timestamp=\\\"" + auth["oauth_timestamp"] + "\\\",oauth_token=\\\"" + auth["oauth_token"] + "\\\",oauth_version=\\\"" + auth["oauth_version"] + "\\\",oauth_signature=\\\"" + misc.OauthEscape(signature) + "\\\""
	} else {
		authHeader = "OAuth oauth_consumer_key=\\\"" + auth["oauth_consumer_key"] + "\\\", oauth_token=\\\"" + auth["oauth_token"] + "\\\", oauth_signature_method=\\\"" + auth["oauth_signature_method"] + "\\\", oauth_timestamp=\\\"" + auth["oauth_timestamp"] + "\\\", oauth_nonce=\\\"" + auth["oauth_nonce"] + "\\\", oauth_version=\\\"" + auth["oauth_version"] + "\\\", oauth_signature=\\\"" + misc.OauthEscape(signature) + "\\\""
	}
	return authHeader, nil
}

//sortAuth creates a sorted Authentication string. The string is sorted Lexographically.
func sortAuth(auth map[string]string) string {
	var keys []string
	for k := range auth {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedAuthString := ""
	for _, k := range keys {
		sortedAuthString = sortedAuthString + k + "=" + misc.OauthEscape(auth[k]) + "&"
	}
	sortedAuthString = sortedAuthString[:len(sortedAuthString)-1]
	return sortedAuthString
}
