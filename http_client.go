package logger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/atomic"
)

type protocol string

const (
	pWSS   protocol = "wss"
	pWS    protocol = "ws"
	pHTTPS protocol = "https"
	pHTTP  protocol = "http"
)
const (
	EndpointLog string = "/log"
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
	return c.WriteLog(c.pref, "PING")
}

func NewHttpClient(cfg ClientConfig) (c *HttpClient, err error) {
	c = new(HttpClient)
	c.port = cfg.Port
	c.remoteIP = cfg.RemoteIP
	c.apiKey = cfg.APIKey
	r := c.buildRequestBody([]string{c.pref, "client initialized"})
	c.logLocally = atomic.NewBool(cfg.LogLocally)

	if cfg.EnableTLS {
		c.protocol = pHTTPS
	} else {
		c.protocol = pHTTP
	}

	c.logURI = c.baseURI() + EndpointLog

	res, err := http.DefaultClient.Post(c.logURI, "application/octet-stream", r)
	if err != nil {
		return nil, errors.New("unable to connect to client: " + err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unable to connect to client: ")
	}
	return
}

func (c HttpClient) WriteLog(inp ...string) (err error) {
	_, err = http.DefaultClient.Post(c.logURI, "application/octet-stream", c.buildRequestBody(inp))

	if c.LogLocally() {
		Info(inp...)
	}

	return err
}

type LogBody struct {
	Key string   `msgpack:"k,omitempty"`
	Log []string `msgpack:"l,omitempty"`
}

func (c *HttpClient) buildRequestBody(v []string) io.Reader {

	var body LogBody = LogBody{
		Key: c.apiKey,
		Log: v,
	}
	b, _ := msgpack.Marshal(body)
	return bytes.NewReader(b)
}
