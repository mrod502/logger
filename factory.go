package logger

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	gocache "github.com/mrod502/go-cache"
	"go.uber.org/atomic"
)

//NewClient -
func NewClient(cfg ClientConfig) (Client, error) {
	if cfg.EnableWebsocket {
		return NewWebsocketClient(cfg)
	}
	return NewHttpClient(cfg)
}

func NewLogServer(cfg ServerConfig) (l LogServer, err error) {
	if cfg.EnableWebsocket {
		return NewWebsocketServer(cfg)
	}
	return NewHttpServer(cfg)
}

func NewFileLog(filePath string, c chan []string) (*FileLog, error) {
	l := new(FileLog)
	f, err := openLogFile(filePath)
	l.f = f
	l.c = c
	l.ctr = new(atomic.Uint32)
	l.syncInterval = atomic.NewUint32(20)
	l.queueEmpty = &sync.WaitGroup{}

	return l, err
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

func NewHttpServer(cfg ServerConfig) (*HttpServer, error) {
	c := make(chan []string, 1024)
	l, err := NewFileLog("", c)
	if err != nil {
		return nil, err
	}
	ls := &HttpServer{
		logger: l,
		notify: c,
	}
	return ls, nil
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
	return
}

func NewWebsocketServer(cfg ServerConfig) (*WebsocketServer, error) {
	var w = new(WebsocketServer)
	var err error

	w.apiKeys = cfg.KeySignatures
	w.cert = cfg.CertFilePath
	w.conns = gocache.NewInterfaceCache()
	w.failedConnAttempts = gocache.NewIntCache()
	w.key = cfg.KeyFilePath
	w.logger, err = NewFileLog(cfg.LogPath, make(chan []string, 1024))
	if err != nil {
		return nil, err
	}
	w.maxFailedConnAttempts = 20
	w.port = cfg.Port
	w.buildRouter()
	w.tls = cfg.EnableTLS

	w.upgrader = &websocket.Upgrader{
		HandshakeTimeout:  time.Second,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return w, nil
}
