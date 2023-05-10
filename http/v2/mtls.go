package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// Storage repository for private keys and certificates. In it's current
// implementation a KeyStore can only hold one private key and certificate.
type KeyStore interface {
	// Retrieve a certificate and private key from the store.
	Get() (tls.Certificate, error)
}

// mTLS client options
type Options struct {
	// Location of the client trust store
	TrustStore string
	// An implementation of the KeyStore interface
	KeyStore KeyStore
}

// Creates a new HTTP client that uses mTLS to authenticate itself and verify
// the identity of the server it connects to.
func NewClient(opt Options) (*http.Client, error) {
	var config tls.Config

	if opt.TrustStore != "" {
		caCert, err := os.ReadFile(opt.TrustStore)
		if err != nil {
			return nil, fmt.Errorf("while loading trusted certificates: %w", err)
		}
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			caCertPool = x509.NewCertPool()
		}
		caCertPool.AppendCertsFromPEM(caCert)
		config.RootCAs = caCertPool
	}

	if opt.KeyStore == nil {
		return nil, fmt.Errorf("key store not specified")
	}

	cert, err := opt.KeyStore.Get()
	if err != nil {
		return nil, fmt.Errorf("while loading the mTLS client certificate: %w", err)
	}

	config.Certificates = []tls.Certificate{cert}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &config,
		},
	}, nil
}
