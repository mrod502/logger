package logger

import (
	"fmt"
	"os"
	"testing"
)

func TestSum(t *testing.T) {

	b, err := os.ReadFile("key.txt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
	fmt.Println(Sha256Sum(string(b)), Sha256Sum("abcde"))

}

func TestNewKey(t *testing.T) {
	fmt.Println(Sha256Sum("abcde"))
}
