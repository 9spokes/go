package http

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/9spokes/go/logging/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

// Implementation of the KeyStore interface that uses Azure Key Vault Managed
// HSM for key management. Only RSA keys are supported.
type KeyVault struct {
	// Tenant id
	Tenant string
	// HSM name
	HSMName string
	// Private key id
	Key string
	// Private key version
	KeyVersion string

	// Unfortunately, Azure Key Vault cannot store certificates so we'll need
	// use a local file instead.
	CertificateFile string
	client          *azkeys.Client
}

func (kv *KeyVault) Get() (tls.Certificate, error) {
	if err := kv.validate(); err != nil {
		return tls.Certificate{}, fmt.Errorf("HSM name not specified")
	}

	err := kv.connectToKeyVault()
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("while creating HSM signer")
	}

	crt, err := kv.loadCert()
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("while loading certificate")
	}

	return tls.Certificate{
		PrivateKey:  kv,
		Leaf:        crt,
		Certificate: [][]byte{crt.Raw},
	}, nil
}

func (kv *KeyVault) validate() error {
	if kv.HSMName == "" {
		return fmt.Errorf("HSM name not specified")
	}

	if kv.Key == "" {
		return fmt.Errorf("key name not specified")
	}

	if kv.KeyVersion == "" {
		return fmt.Errorf("key version not specified")
	}

	if kv.CertificateFile == "" {
		return fmt.Errorf("certificate file not specified")
	}

	return nil
}

func (kv *KeyVault) loadCert() (*x509.Certificate, error) {
	certificate, err := os.ReadFile(kv.CertificateFile)
	if err != nil {
		return nil, fmt.Errorf("while reading certificate: %w", err)
	}

	// Our decoded DER certificate
	var der []byte

	// Assume a certificate chain is provided with proper BEGIN and END lines
	block, _ := pem.Decode(certificate)
	if block == nil {
		// If the Decoding fails, assume it is a single Base64-encoded DER certificate
		der, err = base64.StdEncoding.DecodeString(string(certificate))
		// If the decoding fails, we're unable to handle the data
		if err != nil {
			return nil, fmt.Errorf("could not decode x509 certificate: %s", certificate)
		}
	} else {
		der = block.Bytes
	}

	return x509.ParseCertificate(der)
}

func (ctx *KeyVault) connectToKeyVault() error {
	creds, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{
			TenantID: ctx.Tenant,
		})
	if err != nil {
		return fmt.Errorf("while getting Azure credentials: %w", err)
	}

	url := fmt.Sprintf("https://%s.managedhsm.azure.net/", ctx.HSMName)
	az, err := azkeys.NewClient(url, creds, nil)
	if err != nil {
		return fmt.Errorf("while creating Azure Key Vault client: %w", err)
	}

	ctx.client = az
	return nil
}

func (kv *KeyVault) Public() crypto.PublicKey {

	keyBundle, err := kv.client.GetKey(context.Background(), kv.Key, kv.KeyVersion, nil)
	if err != nil {
		logging.Errorf("Failed to retreive public key from Key Vault")
		return nil
	}

	if keyBundle.Key.Kty == nil {
		logging.Errorf("Failed to determine public key type")
		return nil
	}

	if *keyBundle.Key.Kty != azkeys.JSONWebKeyTypeRSAHSM {
		logging.Errorf("Unsupported public key type: %s", *keyBundle.Key.Kty)
		return nil
	}

	return &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(keyBundle.Key.N),
		E: int(big.NewInt(0).SetBytes(keyBundle.Key.E).Uint64()),
	}
}

func (kv *KeyVault) Sign(_ io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	params := azkeys.SignParameters{
		Value: digest,
	}

	pk := kv.Public()
	if pk == nil {
		return nil, fmt.Errorf("unable to retrieve public key")
	}

	signAlg := ""
	switch pk.(type) {
	case *rsa.PublicKey:
		if _, ok := opts.(*rsa.PSSOptions); ok {
			signAlg = "PS"
		} else {
			signAlg = "RS"
		}
	default:
		return nil, fmt.Errorf("unsupported public key type: %T", pk)
	}

	var hf crypto.Hash = crypto.SHA256
	if opts != nil {
		hf = opts.HashFunc()
	}

	var algo azkeys.JSONWebKeySignatureAlgorithm
	switch hf {
	case crypto.SHA256:
		algo = azkeys.JSONWebKeySignatureAlgorithm(signAlg + "256")
	case crypto.SHA384:
		algo = azkeys.JSONWebKeySignatureAlgorithm(signAlg + "384")
	case crypto.SHA512:
		algo = azkeys.JSONWebKeySignatureAlgorithm(signAlg + "512")
	default:
		return nil, fmt.Errorf("unsupported hashing function: %s", hf.String())
	}

	params.Algorithm = &algo

	res, err := kv.client.Sign(context.Background(), kv.Key, kv.KeyVersion, params, nil)
	if err != nil {
		return nil, fmt.Errorf("while signing digest: %w", err)
	}

	return res.Result, nil
}
