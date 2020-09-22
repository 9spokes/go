package misc

import (
	"net/url"
	"strings"
)

//OauthEscape creates a safe text to use in URL's
func OauthEscape(value string) string {
	escapedQuery := url.QueryEscape(value)
	toReturn := strings.Replace(escapedQuery, "+", "%20", -1)
	return toReturn
}
