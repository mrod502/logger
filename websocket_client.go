package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
)

type WebsocketClient struct {
	prefix     string
	remoteIP   string
	apiKey     string
	protocol   protocol
	port       uint16
	logURI     string
	logLocally *atomic.Bool
	conn       *websocket.Conn
}

func NewWebsocketClient(cfg ClientConfig) (cli *WebsocketClient, err error) {
	cli = new(WebsocketClient)
	cli.prefix = cfg.Prefix
	cli.apiKey = cfg.APIKey
	if cfg.EnableTLS {
		cli.protocol = pWSS
	} else {
		cli.protocol = pWS
	}
	cli.remoteIP = cfg.RemoteIP
	cli.port = cfg.Port
	cli.logURI = cli.baseURI()
	cli.logLocally = atomic.NewBool(cfg.LogLocally)

	return
}

func (c WebsocketClient) baseURI() string {
	return fmt.Sprintf("%s://%s:%d", c.protocol, c.remoteIP, c.port)
}
func (c *WebsocketClient) SetLogLocally(bool) {

}
func (c *WebsocketClient) LogLocally() bool {
	return c.logLocally.Load()
}
func (c *WebsocketClient) WriteLog(...string) error {
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
