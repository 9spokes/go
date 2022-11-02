package jwt

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/9spokes/go/logging/v3"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/square/go-jose.v2"
)

const (
	KID = "kid"
	X5C = "x5c"
)

// Context holds the config required to parse and validate a token
type Context struct {
	JWKSURLs     string
	TrustedKeys  map[string]rsa.PublicKey
	TrustedCerts []x509.Certificate
	PrivateKey   *rsa.PrivateKey
}

// New creates a new JWT context
func New(jwksURLs, trustStorePath, privateKeyPath, privateKeyPwd string) (*Context, error) {
	ctx := Context{
		JWKSURLs:     jwksURLs,
		TrustedKeys:  make(map[string]rsa.PublicKey),
		TrustedCerts: make([]x509.Certificate, 0),
		PrivateKey:   nil,
	}

	if jwksURLs != "" {

		keys, err := fetchJWKS(jwksURLs)
		if err != nil {
			return nil, fmt.Errorf("while retrieving web keys: %s", err.Error())
		}

		ctx.TrustedKeys = keys
	}

	if trustStorePath != "" {
		certs, err := loadCerts(trustStorePath)
		if err != nil {
			return nil, fmt.Errorf("while loading trusted certificates: %s", err.Error())
		}
		ctx.TrustedCerts = certs
	}

	if privateKeyPath != "" {
		// load private key
		key, err := loadPrivateKey(privateKeyPath, privateKeyPwd)
		if err != nil {
			return nil, fmt.Errorf("while loading private key: %s", err.Error())
		}
		ctx.PrivateKey = key
	}

	return &ctx, nil
}

// Validate checks the signature, decrypts if necessary and verifies the
// standard claims to ensure that a token is valid. Retruns the token claims if
// everythign is ok.
func (ctx *Context) Validate(input string) (map[string]interface{}, error) {
	// JWE's are made up of 5 parts (see https://www.rfc-editor.org/info/rfc7516)
	// JWS's are made up of 3 parts (see https://www.rfc-editor.org/info/rfc7515)
	if len(strings.Split(input, ".")) == 5 {
		decrypted, err := ctx.decrypt(input)
		if err != nil {
			return nil, fmt.Errorf("while decrypting token: %s", err.Error())
		}

		input = decrypted
	}

	token, err := jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {
		return ctx.getSigningKey(token)
	})
	if err != nil {
		return nil, fmt.Errorf("while parsing token: %s", err.Error())
	}

	return token.Claims.(jwt.MapClaims), nil
}

// Decrypts the token returning the payload as a string
func (ctx *Context) decrypt(tokenString string) (string, error) {
	token, err := jose.ParseEncrypted(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %s", err.Error())
	}

	if ctx.PrivateKey == nil {
		return "", fmt.Errorf("missing private key")
	}

	payload, err := token.Decrypt(ctx.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt token: %s", err.Error())
	}

	return string(payload), nil
}

// Identifies and returns the token signing key
func (ctx *Context) getSigningKey(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	// Get public key
	keyType, err := getKeyType(token.Header)
	if err != nil {
		return nil, fmt.Errorf("unrecognised public key type '%s': %s", keyType, err.Error())
	}

	switch keyType {
	case KID:
		kid := token.Header[KID].(string)
		publicKey, ok := ctx.TrustedKeys[kid]
		// If key is missing, refresh the keyMap just in case the keys were
		// rotated since we last read them
		if !ok {
			keys, err := fetchJWKS(ctx.JWKSURLs)
			if err != nil {
				return nil, fmt.Errorf("while retrieving web keys: %s", err.Error())
			}

			ctx.TrustedKeys = keys
			publicKey, ok = ctx.TrustedKeys[kid]
		}

		if !ok {
			return nil, fmt.Errorf("no key found for id [%s]", kid)
		}

		return &publicKey, nil

	case X5C:
		var certificate string

		if x5c, ok := token.Header[X5C].([]interface{}); ok {
			certificate = x5c[0].(string)
		} else if x5c, ok := token.Header[X5C].(string); ok {
			certificate = x5c
		} else {
			return nil, fmt.Errorf("invalid x5c header, not string and not array of strings")
		}

		// Our decoded DER certificate
		var der []byte

		// Assume a certificate chain is provided with proper BEGIN and END lines
		block, _ := pem.Decode([]byte(certificate))
		if block == nil {
			// If the Decoding fails, assume it is a single Base64-encoded DER certificate
			der, err = base64.StdEncoding.DecodeString(certificate)
			// If the decoding fails, we're unable to handle the data
			if err != nil {
				return nil, fmt.Errorf("could not decode x509 certificate: %s", certificate)
			}
		} else {
			der = block.Bytes
		}

		var cert *x509.Certificate
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			return nil, fmt.Errorf("could not load x5c certificate")
		}

		if !isJWSAuthorized(cert, ctx.TrustedCerts) {
			return nil, fmt.Errorf("certificate is not trusted")
		}

		publicKey := cert.PublicKey.(*rsa.PublicKey)

		return publicKey, nil
	}

	return nil, nil
}

// Checks whether this certificate (identified by thumbprint) is part of the trust keystore
func isJWSAuthorized(cert *x509.Certificate, pool []x509.Certificate) bool {
	thumbprint := sha256.Sum256(cert.Raw)

	for _, c := range pool {
		if sha256.Sum256(c.Raw) == thumbprint {
			return true
		}
	}
	return false
}

// getKeyType determines what type of key a JWT is using
func getKeyType(joseHeader map[string]interface{}) (string, error) {

	var keyType string

	// KID ?
	_, ok := joseHeader[KID]
	if ok {
		// logger.Debug("JOSE Key type KID")
		keyType = KID
	}

	// X5c ?
	_, ok = joseHeader[X5C]
	if ok {
		// logger.Debug("JOSE Key type X5C")
		keyType = X5C
	}

	// Unknown key type !
	if keyType == "" {
		// logger.Debug("JOSE Key type Unrecognised")
		return "", errors.New("JOSE key type is unrecognised")
	}

	return keyType, nil
}

// fetchJWKS retrieves a JSON Web Key Set
func fetchJWKS(jwksURLs string) (map[string]rsa.PublicKey, error) {
	keys := make(map[string]rsa.PublicKey)

	for _, jwksURL := range strings.Split(jwksURLs, " ") {

		// Get open id connect configuration
		response, err := http.Get(jwksURL)
		if err != nil {
			return keys, fmt.Errorf("while connecting to remote endpoint '%s': %s", jwksURL, err.Error())
		}
		data, _ := ioutil.ReadAll(response.Body)

		var set jose.JSONWebKeySet
		if err := json.Unmarshal(data, &set); err != nil {
			return keys, err
		}

		// Copy key to key map
		for _, key := range set.Keys {
			keys[key.KeyID] = *(key.Key.(*rsa.PublicKey))
		}
	}

	return keys, nil
}

// loadCerts loads certificates in the trust store
func loadCerts(path string) ([]x509.Certificate, error) {
	certs := make([]x509.Certificate, 0)

	raw, err := os.ReadFile(path)
	if err != nil {
		return certs, fmt.Errorf("while reading keystore %s: %s", path, err.Error())
	}

	for i := 0; ; i++ {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			logging.Warningf("failed to load certificate #%d: %s", i+1, err.Error())
		} else {
			certs = append(certs, *cert)
		}
		raw = rest
	}

	return certs, nil
}

// loadPrivateKey loads a RSA private key from the disk
func loadPrivateKey(path, secret string) (*rsa.PrivateKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("while reading private key file '%s': %s", path, err.Error())
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("while decoding private key '%s'", path)
	}

	var der []byte
	if secret != "" {
		if der, err = x509.DecryptPEMBlock(block, []byte(secret)); err != nil {
			return nil, fmt.Errorf("cannot decrypt key: %s", err.Error())
		}
	} else {
		der = block.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(der); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(der); err != nil {
			return nil, fmt.Errorf("while parsing private key '%s': %s", path, err.Error())
		}
	}

	key, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unable to parse RSA private key")
	}

	return key, nil
}
