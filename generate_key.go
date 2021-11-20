package logger

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
)

func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func MarshalKeyPair(k *rsa.PrivateKey) (priv, pub []byte, err error) {

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
	buf := newBuffer()
	err = pem.Encode(buf, privBlock)
	if err != nil {
		return
	}

	buf1 := newBuffer()

	err = pem.Encode(buf1, pubBlock)

	return buf.b, buf1.b, err
}

func GenerateApiKey() (string, string) {
	var b = make([]byte, 0, 64)
	for i := 0; i < 64; i++ {

		v, _ := rand.Int(rand.Reader, big.NewInt(255))
		b = append(b, byte(v.Int64()))
	}

	key := base64.StdEncoding.EncodeToString(b)

	sig := sha256Sum(key)

	return key, sig
}

type buffer struct {
	b []byte
}

func (b *buffer) Write(bt []byte) (int, error) {

	b.b = append(b.b, bt...)
	return len(bt), nil
}

func newBuffer() *buffer {
	return &buffer{b: make([]byte, 0, 128)}
}
