package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/square/go-jose.v2"
)

var keyMap = make(map[string]rsa.PublicKey)

//Params is a struct defining the required parameters needed to generate a new signed JWT token
type Params struct {
	Subject            string
	Claims             map[string]interface{}
	PrivateKeyPath     string
	PrivateKeyPassword string
	PublicKeyPath      string
	Expiry             time.Duration
}

// ValidateJWT checks that the supplied string constitutes a JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Get public key
		keyType, err := GetKeyType(token.Header)
		if err != nil {
			return nil, fmt.Errorf("Error unrecognised public key type: %s", err.Error())
		}

		switch keyType {

		case "kid":

			kid := token.Header["kid"].(string)
			publicKey, ok := keyMap[kid]
			if !ok {
				return nil, fmt.Errorf("No key found for kid [%s]", kid)
			}
			return &publicKey, nil

		case "x5c":

			certificate := token.Header["x5c"].([]interface{})[0].(string)

			block, _ := pem.Decode([]byte(certificate))

			var cert *x509.Certificate
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("Could not load x5c certificate")
			}

			publicKey := cert.PublicKey.(*rsa.PublicKey)

			return publicKey, nil
		}

		return nil, err
	})

	return token, err
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
		return "", errors.New("JOSE Key type Unrecognised")
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

	// Empty claims
	claims := jwt.MapClaims{}

	// Add addtional claims
	for k, v := range params.Claims {
		claims[k] = v
	}

	// Add standard claims
	claims["iss"] = "9Spokes"
	claims["sub"] = subject
	claims["aud"] = "9spokes"
	claims["exp"] = time.Now().UTC().Add(time.Second * time.Duration(opt.TokenExpiryPeriod)).Unix()
	claims["nbf"] = time.Now().Unix()
	claims["iat"] = time.Now().Unix()
	claims["jti"] = logging.GenUUIDv4()

	// Sign JWT
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["x5c"] = certChain
	jws, err := token.SignedString(priv)
	if err != nil {
		return "", fmt.Errorf("while signing JWT: %s", err.Error())
	}

	return jws, nil
}
