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
	fmt.Println(sha256Sum(string(b)), sha256Sum("abcde"))

}

func TestNewKey(t *testing.T) {
	fmt.Println(sha256Sum("abcde"))
}
