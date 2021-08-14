package logger

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	gocache "github.com/mrod502/go-cache"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/atomic"
)

type WebsocketServer struct {
	done                  *atomic.Bool
	logger                *Log
	notify                chan []string
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
}
func (s *WebsocketServer) Quit() {

}
func (s *WebsocketServer) processQueue() {

}

func (s *WebsocketServer) upgrade(w http.ResponseWriter, r *http.Request) {
	if s.failedConnAttempts.Get(r.RemoteAddr) >= s.maxFailedConnAttempts {
		//s.doLog("BLACKLIST", "upgrade attempt", r.RemoteAddr)
		return
	}
	apiKey := r.Header.Get("API-Key")
	var keySum = sha256.Sum256([]byte(apiKey))
	if !s.apiKeys.Get(string(keySum[:])) {
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

	s.notify <- inp
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
