package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MTLSRequest is a struct that contains the client MTLS credentials, used by New()
type MTLSRequest struct {
	KeyFile         string
	CertificateFile string
	CAFile          string
}

// New returns a pointer to an http.Client instance with the supplied paths to the client  key, certificate, and CA
func (req MTLSRequest) New() (*http.Client, error) {

	var config tls.Config

	if req.CAFile != "" {
		caCert, err := ioutil.ReadFile(req.CAFile)
		if err != nil {
			return nil, fmt.Errorf("while reading openbanking CA certificate %s: %s", req.CAFile, err.Error())
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		config.RootCAs = caCertPool
	}

	if req.CertificateFile != "" && req.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(req.CertificateFile, req.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("while loading the mTLS client certificate & key: %s", err.Error())
		}
		config.Certificates = []tls.Certificate{cert}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &config,
		},
	}
	return client, nil
}
