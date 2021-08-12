package logger

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/atomic"
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

var logLocally *atomic.Bool

func SetLogLocally(l bool) {
	logLocally.Store(l)
}

func LogLocally() bool {
	return logLocally.Load()
}

func NewClient(addr string, logPrefix string) (c *Client, err error) {
	c = new(Client)
	c.addr = addr
	c.pref = logPrefix
	r := buildRequestBody([]string{c.pref, "client initialized"})

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
	_, err = http.DefaultClient.Post(c.logURI(), "application/octet-stream", buildRequestBody(inp))

	if LogLocally() {
		Info(inp...)
	}

	return err
}

type LogBody struct {
	Key string   `msgpack:"k"`
	Log []string `msgpack:"l"`
}

func buildRequestBody(v []string) io.Reader {

	var body LogBody = LogBody{
		Key: "",
		Log: v,
	}
	b, _ := msgpack.Marshal(body)
	return bytes.NewReader(b)
}
