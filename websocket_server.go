package logger

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	gocache "github.com/mrod502/go-cache"
	"github.com/vmihailenco/msgpack/v5"
)

type WebsocketServer struct {
	//done                  *atomic.Bool
	logger                *Log
	upgrader              websocket.Upgrader
	tls                   bool
	port                  uint16
	conns                 *gocache.InterfaceCache
	cert                  string
	key                   string
	apiKeys               *gocache.BoolCache
	maxFailedConnAttempts int
	failedConnAttempts    *gocache.IntCache
}

func NewWebsocketServer(cfg ServerConfig) (*WebsocketServer, error) {
	var w = new(WebsocketServer)
	var err error
	w.logger, err = NewLog(cfg.LogPath, make(chan []string, 1024))
	if err != nil {
		return nil, err
	}
	w.upgrader = websocket.Upgrader{
		HandshakeTimeout:  time.Second,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	w.apiKeys = cfg.KeySignatures
	return w, nil
}

func (s *WebsocketServer) Serve() error {
	router := mux.NewRouter()
	router.HandleFunc("/ws", s.upgrade)

	if s.tls {
		return http.ListenAndServeTLS(fmt.Sprintf(":%d", s.port), s.cert, s.key, router)
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), router)
}

func (s *WebsocketServer) SetSyncInterval(i uint32) {
	s.logger.SetSyncInterval(i)
}
func (s *WebsocketServer) Quit() {
	s.logger.Stop()
}

func (s *WebsocketServer) upgrade(w http.ResponseWriter, r *http.Request) {
	if s.failedConnAttempts.Get(r.RemoteAddr) >= s.maxFailedConnAttempts {
		//s.doLog("BLACKLIST", "upgrade attempt", r.RemoteAddr)
		return
	}
	apiKey := r.Header.Get("API-Key")
	if !s.apiKeys.Get(sha256Sum(apiKey)) {
		s.failedConnAttempts.Add(r.RemoteAddr, 1)
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.doLog("UPGRADE", r.RemoteAddr, err.Error())
	}
	s.conns.Set(conn.RemoteAddr().String(), conn)
	go s.readMessages(conn)
}

func (s *WebsocketServer) doLog(inp ...string) {
	if len(inp) > 0 {
		s.logger.Write(inp...)
	}
}
func (s *WebsocketServer) readMessages(conn *websocket.Conn) {
	if conn == nil {
		s.doLog("READ", "Nil conn")
		return
	}

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			s.doLog("READ", conn.RemoteAddr().String(), err.Error())
			s.conns.Delete(conn.RemoteAddr().String())
			return
		}
		var msg []string
		err = msgpack.Unmarshal(b, &msg)
		if err != nil {
			s.doLog("READ", conn.RemoteAddr().String(), err.Error())
		} else {
			s.doLog(append([]string{"LOG", conn.RemoteAddr().String()}, msg...)...)
		}

	}

}

func sha256Sum(inp string) string {
	var keySum = sha256.Sum256([]byte(inp))
	return string(keySum[:])
}
