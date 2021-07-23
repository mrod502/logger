package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Client struct {
	addr string
	pref string
}

func (c Client) logURI() string {
	return "http://" + c.addr + EndpointLog
}

const (
	EndpointLog = "/log"
)

func NewClient(addr string, logPrefix string) (c *Client, err error) {
	c = new(Client)
	c.addr = addr
	c.pref = logPrefix
	r := stringSlice2Reader([]string{c.pref, "client initialized"})

	res, err := http.DefaultClient.Post("http://"+c.addr+EndpointLog, "text/json", r)
	if err != nil {
		return nil, errors.New("unable to connect to client: " + err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unable to connect to client: ")
	}
	return
}

func (c Client) WriteLog(inp ...string) (err error) {
	_, err = http.DefaultClient.Post(c.logURI(), "text/json", stringSlice2Reader(inp))

	return err
}

func stringSlice2Reader(v []string) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
