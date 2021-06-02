package auth_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/9spokes/go/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorize(t *testing.T) {

	ts := func(statusCode int, response map[string]string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(statusCode)

			if response != nil {
				defer r.Body.Close()
				body, _ := ioutil.ReadAll(r.Body)
				response["query"] = r.URL.RawQuery
				response["auth"] = r.Header.Get("Authorization")
				response["body"] = string(body)
				ret, _ := json.Marshal(response)
				w.Write(ret)
			}
		}))
	}

	tests := []struct {
		name     string
		ctx      *auth.OAuth2
		errorMsg string
		server   *httptest.Server
		options  auth.Options
	}{
		{
			name: "client_id is missing",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				Code:          "bogus",
			},
			errorMsg: "client_id",
		},
		{
			name: "non-OK response",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				Code:          "bogus",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			errorMsg: "non-OK response",
			server:   ts(404, nil),
		},
		{
			name: "Success",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				Code:          "bogus",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			server: ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
		{
			name: "Auth in header",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				Code:          "bogus",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			options: auth.Options{AuthInHeader: true},
			server:  ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
		{
			name: "Data in query",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				Code:          "bogus",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			options: auth.Options{DataInQuery: true},
			server:  ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
		{
			name: "Auth in header, Data in query",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				Code:          "bogus",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			options: auth.Options{AuthInHeader: true, DataInQuery: true},
			server:  ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			if tt.server != nil {
				defer tt.server.Close()
				tt.ctx.TokenEndpoint = tt.server.URL
			}
			ret, err := tt.ctx.Authorize(tt.options)

			if tt.errorMsg == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tt.errorMsg)
				return
			}

			var data url.Values
			if tt.options.DataInQuery {
				data, _ = url.ParseQuery(ret["query"].(string))
			} else {
				data, _ = url.ParseQuery(ret["body"].(string))

				if !tt.options.AuthInHeader {
					assert.Equal(tt.ctx.ClientSecret, data.Get("client_secret"), "incorrect client_secret sent as auth")
				}
			}

			assert.Equal(tt.ctx.Code, data.Get("code"), "incorrect code sent in auth request")

			assert.Equal(tt.ctx.RedirectURI, data.Get("redirect_uri"), "incorrect redirect_uri sent in auth request")

			assert.Equal("authorization_code", data.Get("grant_type"), "incorrect grant_type sent in auth request")

			if tt.options.AuthInHeader {
				authHeader, ok := ret["auth"].(string)
				require.True(ok, "Authorization header must be string")
				require.NotEmpty(authHeader, "Authorization header")

				assert.Regexp("^basic ", authHeader, "incorrect auth header")
				got, _ := base64.StdEncoding.DecodeString(strings.Split(authHeader, " ")[1])
				want := fmt.Sprintf(
					"%s:%s",
					tt.ctx.ClientID,
					tt.ctx.ClientSecret)
				assert.EqualValues(want, got, "incorrect auth")

			} else {
				assert.Equal(tt.ctx.ClientID, data.Get("client_id"), "incorrect client_id sent as auth")
			}
		})
	}
}

func TestRefresh(t *testing.T) {

	ts := func(statusCode int, response map[string]string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")

			w.WriteHeader(statusCode)

			if response != nil {
				defer r.Body.Close()
				body, _ := ioutil.ReadAll(r.Body)
				response["query"] = r.URL.RawQuery
				response["auth"] = r.Header.Get("Authorization")
				response["body"] = string(body)
				ret, _ := json.Marshal(response)
				w.Write(ret)
			}
		}))
	}

	tests := []struct {
		name     string
		ctx      *auth.OAuth2
		errorMsg string
		server   *httptest.Server
		options  auth.Options
	}{
		{
			name: "client_id is missing",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				RefreshToken:  "bogus-refresh-token",
			},
			errorMsg: "client_id",
		},
		{
			name: "The refresh_token is missing",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			errorMsg: "the refresh token is missing",
		},
		{
			name: "non-OK response",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				RefreshToken:  "bogus-refresh-token",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			errorMsg: "non-OK response",
			server:   ts(404, nil),
		},
		{
			name: "Success",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				RefreshToken:  "bogus-refresh-token",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			server: ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
		{
			name: "Success/Auth in header",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				RefreshToken:  "bogus-refresh-token",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			options: auth.Options{AuthInHeader: true},
			server:  ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
		{
			name: "Success/Data in query",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				RefreshToken:  "bogus-refresh-token",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			options: auth.Options{DataInQuery: true},
			server:  ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
		{
			name: "Success/Auth in header, Data in query",
			ctx: &auth.OAuth2{
				TokenEndpoint: "http://bogus",
				RefreshToken:  "bogus-refresh-token",
				ClientID:      "bogus",
				ClientSecret:  "bogus",
			},
			options: auth.Options{AuthInHeader: true, DataInQuery: true},
			server:  ts(200, map[string]string{"access_token": "bogus", "refresh_token": "bogus-refresh-token"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.server != nil {
				defer tt.server.Close()
				tt.ctx.TokenEndpoint = tt.server.URL
			}
			ret, err := tt.ctx.Refresh(tt.options)

			if err != nil {
				if tt.errorMsg == "" {
					t.Errorf("unexpected error: %s", err.Error())
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("incorrect error, got=(%s), want=(%s)", err.Error(), tt.errorMsg)
				}
				return
			} else if tt.errorMsg != "" {
				t.Errorf("no error, expected error: %s", tt.errorMsg)
				return
			}

			var data url.Values
			if tt.options.DataInQuery {
				data, _ = url.ParseQuery(ret["query"].(string))
			} else {
				data, _ = url.ParseQuery(ret["body"].(string))

				if !tt.options.AuthInHeader {
					wantClientSecret := tt.ctx.ClientSecret
					if clientSecret := data.Get("client_secret"); clientSecret != wantClientSecret {
						t.Errorf("incorrect client_secret sent as auth got=(%s) want=(%s)", clientSecret, wantClientSecret)
					}
				}
			}

			wantRefreshToken := tt.ctx.RefreshToken
			if refreshToken := data.Get("refresh_token"); refreshToken != wantRefreshToken {
				t.Errorf("incorrect refresh_token  sent in auth request got=(%s) want=(%s)", refreshToken, wantRefreshToken)
			}
			wantRedirect := tt.ctx.RedirectURI
			if redirect := data.Get("redirect_uri"); redirect != wantRedirect {
				t.Errorf("incorrect redirect_uri sent in auth request got=(%s) want=(%s)", redirect, wantRedirect)
			}
			wantGrantType := "refresh_token"
			if grantType := data.Get("grant_type"); grantType != wantGrantType {
				t.Errorf("incorrect grant_type sent in auth request got=(%s) want=(%s)", grantType, wantGrantType)
			}

			if tt.options.AuthInHeader {
				if authHeader, _ := ret["auth"].(string); authHeader == "" {
					t.Errorf("Authorization header not set")
				} else {
					if !strings.HasPrefix(authHeader, "basic ") {
						t.Errorf("incorrect auth header, got=(%s) wanted with pefix=(basic )", authHeader)
					}
					got, _ := base64.StdEncoding.DecodeString(strings.Split(authHeader, " ")[1])
					want := fmt.Sprintf(
						"%s:%s",
						tt.ctx.ClientID,
						tt.ctx.ClientSecret)
					if string(got) != want {
						t.Errorf("incorrect auth, want=(%s) got=(%s)", want, got)
					}
				}
			} else {
				wantClientId := tt.ctx.ClientID
				if clientId := data.Get("client_id"); clientId != wantClientId {
					t.Errorf("incorrect client_id sent as auth got=(%s) want=(%s)", clientId, wantClientId)
				}
			}
		})
	}
}
