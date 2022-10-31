package logger

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {

	k, err := GenerateKey(1 << 12)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", *k)

	priv, pub, err := MarshalKeyPair(k)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(priv), len(pub))

	fmt.Println(string(priv))
	fmt.Println(string(pub))
}

func TestBuffer(t *testing.T) {
	var msg = []byte("hello")
	b := make([]byte, 0)

	buf := &buffer{b: b}

	buf.Write(msg)

	if !bytes.Equal(buf.b, msg) {
		t.Fatalf("expected:%v, got: %v", msg, buf.b)
	}
}

func TestGenerateApiKey(t *testing.T) {

	key, sig := GenerateApiKey()

	fmt.Println("key:", key)
	fmt.Println("sig:", sig)
	if Sha256Sum(key) != sig {
		t.Fatal("key signature did not match")
	}
}
