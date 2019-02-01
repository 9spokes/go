package crypto

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/pbkdf2"
)

//Decrypt uses a PBKDF2 key encryption method with a AES-CBC 256 algorithm
func Decrypt(ciphertext []byte, secret []byte) (string, error) {

	ciphertext = ciphertext[4:]

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	key := pbkdf2.Key(secret, iv, 2048, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	unpadded, err := func(b []byte, blocksize int) ([]byte, error) {

		if len(b)%blocksize != 0 {
			return nil, errors.New("Invalid PKCS7 padding size")
		}
		c := b[len(b)-1]
		n := int(c)
		if n == 0 || n > len(b) {
			return nil, errors.New("Invalid PKCS7 padding size")
		}
		for i := 0; i < n; i++ {
			if b[len(b)-n+i] != c {
				return nil, errors.New("Invalid PKCS7 padding size")
			}
		}
		return b[:len(b)-n], nil
	}(ciphertext, aes.BlockSize)

	return string(unpadded), err

}

//Encrypt uses a PBKDF2 key encryption method with a AES-CBC 256 algorithm
func Encrypt(str string, secret []byte) ([]byte, error) {

	header := make([]byte, 4+aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, header); err != nil {
		return nil, err
	}

	iv := header[4:]

	bytes := func(b []byte, blocksize int) []byte {
		n := blocksize - (len(b) % blocksize)
		pb := make([]byte, len(b)+n)
		copy(pb, b)
		copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
		return pb
	}([]byte(str), aes.BlockSize)

	key := pbkdf2.Key(secret, iv, 2048, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(header)+len(bytes))
	copy(ciphertext, header)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[len(header):], bytes)

	return ciphertext, nil
}

//SignRSA creates the signature for oauth1 with rsa-sha1
func SignRSA(message []byte, filepath string) string {
	keyInfo, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Panic(err)
	}
	keyBlock, _ := pem.Decode(keyInfo)
	privatekey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		log.Panic(err)
	}
	hash := crypto.SHA1.New()
	hash.Write(message)
	data := hash.Sum(nil)
	signed, err := rsa.SignPKCS1v15(rand.Reader, privatekey, crypto.SHA1, data)
	if err != nil {
		log.Panic(err)
	}
	return base64.StdEncoding.EncodeToString(signed)
}
