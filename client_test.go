package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	rasPiIp := "192.168.1.133:1738"
	c, err := NewClient(rasPiIp, "TESTCLIENT")

	if err != nil {
		t.Fatal(err)
	}
	var tnow time.Time
	var tt time.Duration
	tnow = time.Now()
	for i := 0; i < 1000; i++ {
		c.WriteLog("Hello", "this is a test at", time.Now().Format("2006-01-02 03:04:05.99"))
	}
	tt = time.Since(tnow)
	fmt.Println(tt)
}
