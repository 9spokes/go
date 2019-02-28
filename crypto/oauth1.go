package crypto

import (
	"sort"
	"strings"

	"github.com/9spokes/go/misc"
)

//Signer generates a signature for a xero request
func Signer(auth map[string]string) (string, error) {

	url := getKey(auth, "BaseURL")
	secret := getKey(auth, "Secret")
	method := getKey(auth, "Method")

	sortedAuthString := func(Auth map[string]string) string {
		keys := make([]string, 0, len(Auth))
		for K := range Auth {
			keys = append(keys, K)
		}
		sort.Strings(keys)
		var sortedAuthString strings.Builder
		for _, K := range keys {
			sortedAuthString.WriteString(K + "=" + misc.OauthEscape(Auth[K]) + "&")
		}
		return sortedAuthString.String()[:len(sortedAuthString.String())-1]
	}(auth)

	signature := Sign([]byte(strings.Join([]string{method, misc.OauthEscape(url), misc.OauthEscape(sortedAuthString)}, "&")), auth["oauth_signature_method"], secret)
	return signature, nil
}

func getKey(input map[string]string, key string) string {

	if _, ok := input[key]; !ok {
		panic("Missing " + key)
	}

	ret := input[key]
	delete(input, key)
	return ret
}
