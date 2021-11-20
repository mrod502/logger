package logger

import (
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
