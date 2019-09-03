package main

import (
	"fmt"
	"testing"

	"github.com/9spokes/go/jwt"
)

func TestMakeJWT(t *testing.T) {

	opt.CertificatePath = "certificate.pem"
	opt.KeyPath = "private.pem"
	opt.KeyPassword = "password"

	jws, err := jwt.MakeJWT(jwt.Params{Subject: "123456789", Claims: map[string]string{}})
	if err != nil {
		t.Fail()
	}

	fmt.Printf("JWS : [%s]\n", jws)
}

func TestValidateJWT_X5c(t *testing.T) {

	opt.CertificatePath = "certificate.pem"
	opt.KeyPath = "private.pem"
	opt.KeyPassword = "password"

	jws, err := jwt.MakeJWT(jwt.Params{Subject: "123456789", Claims: map[string]string{}})
	if err != nil {
		t.Fail()
	}

	token, err := jwt.ValidateJWT(jws)
	if err != nil {
		t.Fail()
	}

	fmt.Printf("TOKEN : [%v]\n", token.Claims)
}

func TestValidateJWT_KID(t *testing.T) {

	t.SkipNow() // Test intended to be used manually

	var tokenString = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6InJJcDhlZlE2MDZHZ1l6bmZhVEtQeEZOcllMLXhvLXlQczhKZ3NrMGMwRlEifQ.eyJleHAiOjE1NjU2Njg4MDIsIm5iZiI6MTU2NTY2NTIwMiwidmVyIjoiMS4wIiwiaXNzIjoiaHR0cHM6Ly85c3Bva2VzdGVzdC5iMmNsb2dpbi5jb20vMzA1NTM5MWMtMjAwNi00MDE5LWE0YTYtZTlmMTgyOGMxZDA1L3YyLjAvIiwic3ViIjoiZjU3NjdkZDEtMmEwNi00MGI4LWIzY2EtNTAwMGRmN2E4ZmZlIiwiYXVkIjoiNTNkZGExOWItYjZmOS00MjVhLWI5YmUtMjFiNDc5YWE2YmZlIiwiYWNyIjoiYjJjXzFhX3NpZ251cF9zaWduaW4iLCJub25jZSI6ImRlZmF1bHROb25jZSIsImlhdCI6MTU2NTY2NTIwMiwiYXV0aF90aW1lIjoxNTY1NjY1MjAyLCJnaXZlbl9uYW1lIjoiTXIiLCJmYW1pbHlfbmFtZSI6IlJhYmJpdCJ9.nj6oVTUTmr_KZ1fSIevU_efbNo-Oe203FuvFxe0UyKvtg_f3s6bnz1RyOuT8WlJ9T1_9dOpqKUCk2lGiAcI6CA5dUPswLOxxVCsJ45zzCTJU7Li7BnwfywPh-l56FoduYP3j7gnX4fwjCBpqTex936_lP-HDOejJqk64TwUCgn4pH2B4H-wAoUlEpjtfBgzi4ckQW26CiAW7loMkd6QPOWYVGbchDOB-aY3SOdfAwZuysoG9U5wp_gX4I45AL2P6CogfFilEuaA9iuOJnX40NAdYcpBxzOArWAdeTCN4JU1qtl6cJbIseOVabQMlXW3tgS0kLYcp383Xs4xxitfM8g"

	jwt.FetchJWKS("https://9spokestest.b2clogin.com/9spokestest.onmicrosoft.com/discovery/v2.0/keys?p=b2c_1a_signup_signin")

	_, err := jwt.ValidateJWT(tokenString)

	if err != nil {
		fmt.Printf("\n%s\n", err)
		t.Fail()
	}

}
func TestFetchJWKS(t *testing.T) {

	t.SkipNow() // Test intended to be used manually

	jwt.FetchJWKS("https://9spokestest.b2clogin.com/9spokestest.onmicrosoft.com/discovery/v2.0/keys?p=b2c_1a_signup_signin")

}

func TestFetchOIDCConfiguration(t *testing.T) {

	t.SkipNow() // Test intended to be used manually

	jwt.FetchOIDCConfiguration("https://9SpokesTest.b2clogin.com/9SpokesTest.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_signup_signin")

}
