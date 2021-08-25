package logger

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/atomic"
)

var (
	ErrCacheFull = errors.New("log cache is full")
)

type WebsocketClient struct {
	prefix           string
	remoteIP         string
	apiKey           string
	protocol         protocol
	port             uint16
	logURI           string
	logChan          chan []string
	logLocally       *atomic.Bool
	conn             *websocket.Conn
	brokenConnNotify chan bool
	brokenConn       *atomic.Bool
}

func NewWebsocketClient(cfg ClientConfig) (cli *WebsocketClient, err error) {
	cli = new(WebsocketClient)
	cli.prefix = cfg.Prefix
	cli.apiKey = cfg.APIKey
	cli.logChan = make(chan []string, 512)
	cli.brokenConnNotify = make(chan bool, 128)
	if cfg.EnableTLS {
		cli.protocol = pWSS
	} else {
		cli.protocol = pWS
	}
	cli.remoteIP = cfg.RemoteIP
	cli.port = cfg.Port
	cli.logURI = cli.baseURI()
	cli.logLocally = atomic.NewBool(cfg.LogLocally)
	cli.brokenConn = atomic.NewBool(false)
	go cli.connectionHandler()
	go cli.logWriter()
	return
}

func (c WebsocketClient) baseURI() string {
	return fmt.Sprintf("%s://%s:%d", c.protocol, c.remoteIP, c.port)
}
func (c *WebsocketClient) SetLogLocally(v bool) {
	c.logLocally.Store(v)

}
func (c *WebsocketClient) LogLocally() bool {
	return c.logLocally.Load()
}

func (c *WebsocketClient) logWriter() {
	for {
		l := <-c.logChan
		b, err := msgpack.Marshal(l)
		if err != nil {
			errorLog("MARSHAL", err.Error())
			continue
		}
		err = c.conn.WriteMessage(websocket.BinaryMessage, b)
		if err != nil {
			errorLog("WRITE", err.Error())
			if err == websocket.ErrCloseSent {
				c.notifyBrokenConn()
			}
		}
	}
}

func (c *WebsocketClient) notifyBrokenConn() {
	if c.brokenConn.Load() {
		return
	}
	c.brokenConn.Store(true)
	c.brokenConnNotify <- true
}

func (c *WebsocketClient) WriteLog(l ...string) error {
	if cap(c.logChan) == len(c.logChan) {
		return ErrCacheFull
	}
	c.logChan <- l
	return nil
}
func (c *WebsocketClient) Connect() error {
	h := make(http.Header)
	h.Set("API-Key", c.apiKey)
	dialer := &websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}
	conn, _, err := dialer.Dial(c.logURI, h)
	c.conn = conn
	if err != nil {
		errorLog("DIAL", err.Error())
	}
	return err
}

func (c *WebsocketClient) connectionHandler() {
	for {
		<-c.brokenConnNotify
		func() {
			defer c.brokenConn.Store(false)
			for err := c.Connect(); err != nil; {
				errorLog("CONN", "unable to establish connection to ", c.remoteIP, err.Error())
				time.Sleep(time.Second * 5)
			}
		}()
	}
}
