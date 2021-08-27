package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/atomic"
)

type HttpClient struct {
	pref       string
	remoteIP   string
	apiKey     string
	protocol   protocol
	port       uint16
	logURI     string
	logLocally *atomic.Bool
}

func (c HttpClient) baseURI() string {
	return fmt.Sprintf("%s://%s:%d", c.protocol, c.remoteIP, c.port)
}

func (c *HttpClient) SetLogLocally(l bool) {
	c.logLocally.Store(l)
}

func (c *HttpClient) LogLocally() bool {
	return c.logLocally.Load()
}

func (c *HttpClient) Connect() error {
	return c.Write(c.pref, "PING")
}

func (c *HttpClient) Write(inp ...string) (err error) {
	if c.LogLocally() {
		Info(inp...)
	}

	var req *http.Request

	req, err = http.NewRequest(http.MethodPost, c.logURI, c.buildRequestBody(inp))
	if err != nil {
		return
	}
	var res *http.Response
	req.Header.Set("content-type", "application/octet-stream")
	req.Header.Set("API-Key", c.apiKey)

	res, err = http.DefaultClient.Do(req)
	if res.StatusCode > 299 {
		Error("WRITE", fmt.Sprintf("send request: %d", res.StatusCode))
	}
	return err
}

func (c *HttpClient) buildRequestBody(v []string) io.Reader {
	b, _ := msgpack.Marshal(v)
	return bytes.NewReader(b)
}
