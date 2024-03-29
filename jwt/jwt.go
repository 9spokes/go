// Deprecated: insecure, use v2
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
	"time"

	"github.com/9spokes/go/misc"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/square/go-jose.v2"
)

var keyMap = make(map[string]rsa.PublicKey)

// Options are a set of options & flags used to validate a JWT token
type Options struct {
	OIDCDiscoveryURI string
	TrustedSigners   []x509.Certificate
}

//Params is a struct defining the required parameters needed to generate a new signed JWT token
type Params struct {
	Subject            string
	Claims             map[string]string
	PrivateKeyPath     string
	PrivateKeyPassword string
	PublicKeyPath      string
	Expiry             time.Duration
}

// ValidateJWT checks that the supplied string constitutes a JWT token
func ValidateJWT(tokenString string, opt Options) (*jwt.Token, error) {

	if opt.OIDCDiscoveryURI == "" && opt.TrustedSigners == nil {
		return nil, fmt.Errorf("neither an OIDC Discovery URL nor an x509 certificate pool was supplied")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get public key
		keyType, err := GetKeyType(token.Header)
		if err != nil {
			return nil, fmt.Errorf("unrecognised public key type '%s': %s", keyType, err.Error())
		}

		switch keyType {

		case "kid":

			kid := token.Header["kid"].(string)

			publicKey, ok := keyMap[kid]

			// If key is missing, refresh the keyMap
			if !ok {
				oidcConfig, err := FetchOIDCConfiguration(opt.OIDCDiscoveryURI)
				if err != nil {
					return nil, err
				}

				err = FetchJWKS(oidcConfig["jwks_uri"].(string))
				if err != nil {
					return nil, fmt.Errorf(err.Error())
				}

				publicKey, ok := keyMap[kid]
				if !ok {
					return nil, fmt.Errorf("no key found for kid [%s]", kid)
				}
				return &publicKey, nil
			}
			return &publicKey, nil

		case "x5c":

			var certificate string

			if x5c, ok := token.Header["x5c"].([]interface{}); ok {
				certificate = x5c[0].(string)
			} else if x5c, ok := token.Header["x5c"].(string); ok {
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

			if !isJWSAuthorized(cert, opt.TrustedSigners) {
				return nil, fmt.Errorf("certificate is not trusted")
			}

			publicKey := cert.PublicKey.(*rsa.PublicKey)

			return publicKey, nil
		}

		return nil, err
	})
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

// GetKeyType determines what type of key a JWT is using
func GetKeyType(joseHeader map[string]interface{}) (string, error) {

	var keyType string

	// KID ?
	_, ok := joseHeader["kid"]
	if ok {
		// logger.Debug("JOSE Key type KID")
		keyType = "kid"
	}

	// X5c ?
	_, ok = joseHeader["x5c"]
	if ok {
		// logger.Debug("JOSE Key type X5C")
		keyType = "x5c"
	}

	// Unknown key type !
	if keyType == "" {
		// logger.Debug("JOSE Key type Unrecognised")
		return "", errors.New("JOSE key type is unrecognised")
	}

	return keyType, nil
}

// FetchJWKS retrieves a JSON Web Key Set
func FetchJWKS(jwksURL string) error {

	// Get open id connect configuration
	response, err := http.Get(jwksURL)
	if err != nil {
		return fmt.Errorf("while connecting to remote endpoint '%s': %s", jwksURL, err.Error())
	}
	data, _ := ioutil.ReadAll(response.Body)

	var set jose.JSONWebKeySet
	if err := json.Unmarshal(data, &set); err != nil {
		return err
	}

	// Copy key to key map
	for _, key := range set.Keys {

		keyMap[key.KeyID] = *(key.Key.(*rsa.PublicKey))

	}

	return nil
}

// FetchOIDCConfiguration retrieves the OIDC configuration from a remote OIDC discovery URL
func FetchOIDCConfiguration(oidcDiscoveryURI string) (map[string]interface{}, error) {

	// Get open id connect configuration
	response, err := http.Get(oidcDiscoveryURI)
	if err != nil {
		return nil, fmt.Errorf("while connecting to discovery endpoint '%s': %s", oidcDiscoveryURI, err.Error())
	}
	data, _ := ioutil.ReadAll(response.Body)
	//logger.Debugf("Open ID connection configuration found %s\n", string(data))

	// Unmarshell configuration
	var oidcConfig map[string]interface{}
	err = json.Unmarshal(data, &oidcConfig)
	if err != nil {
		return nil, fmt.Errorf("while unmarshaling OIDC configuration: %s", err.Error())
	}

	return oidcConfig, nil
}

// MakeJWT creates a new signed JWT token given the parameters provided
func MakeJWT(params Params) (string, error) {

	// Load private key from PEM file
	raw, err := ioutil.ReadFile(params.PrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("while reading private key file '%s': %s", params.PrivateKeyPath, err.Error())
	}

	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return "", fmt.Errorf("while decoding private key '%s': %s", params.PrivateKeyPath, err.Error())
	}

	der, err := x509.DecryptPEMBlock(block, []byte(params.PrivateKeyPassword))
	if err != nil {
		return "", fmt.Errorf("while decrypting private key '%s': %s", params.PrivateKeyPath, err.Error())
	}

	priv, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return "", fmt.Errorf("while parsing private key '%s': %s", params.PrivateKeyPath, err.Error())
	}

	// Load certificate from PEM file
	cert, err := ioutil.ReadFile(params.PublicKeyPath)
	if err != nil {
		return "", fmt.Errorf("while reading public key '%s': %s", params.PublicKeyPath, err.Error())
	}

	// Place certificate into certificate chain ( required for format reasons )
	certChain := []string{string(cert)}

	// Add addtional claims
	claims := jwt.MapClaims{

		"iss": "9Spokes",
		"sub": params.Subject,
		"aud": "9spokes",
		"exp": time.Now().UTC().Add(time.Second * params.Expiry).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"jti": misc.GenUUIDv4(),
	}

	for k, v := range params.Claims {
		claims[k] = v
	}

	// Sign JWT
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["x5c"] = certChain
	jws, err := token.SignedString(priv)
	if err != nil {
		return "", fmt.Errorf("while signing JWT: %s", err.Error())
	}

	return jws, nil
}
