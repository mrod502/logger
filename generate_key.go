package logger

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func MarshalKeyPair(k *rsa.PrivateKey) (priv, pub []byte, err error) {
	priv = make([]byte, 0, 1<<9)
	pub = make([]byte, 0, 1<<9)
	privBytes := x509.MarshalPKCS1PrivateKey(k)
	pubBytes, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	buf := bytes.NewBuffer(priv)
	err = pem.Encode(buf, privBlock)
	if err != nil {
		return
	}

	buf1 := bytes.NewBuffer(pub)

	err = pem.Encode(buf1, pubBlock)

	return
}
