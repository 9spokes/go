package http

import (
	"crypto/tls"
	"fmt"
)

// Implementation of the KeyStore interface that loads the certificate and
// private key from the file system.
type FileKeyStore struct {
	CertFile string
	KeyFile  string
}

func (r *FileKeyStore) Get() (tls.Certificate, error) {
	if r.CertFile == "" {
		return tls.Certificate{}, fmt.Errorf("certificate file not specified")
	}

	if r.KeyFile == "" {
		return tls.Certificate{}, fmt.Errorf("key file not specified")
	}

	return tls.LoadX509KeyPair(r.CertFile, r.KeyFile)
}
