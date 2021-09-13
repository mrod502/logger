package logger

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	gocache "github.com/mrod502/go-cache"
	"github.com/vmihailenco/msgpack/v5"
)

type WebsocketServer struct {
	//done                  *atomic.Bool
	apiKeys               *gocache.BoolCache
	cert                  string
	conns                 *gocache.InterfaceCache
	failedConnAttempts    *gocache.IntCache
	key                   string
	logger                *FileLog
	maxFailedConnAttempts uint32
	port                  uint16
	router                *mux.Router
	tls                   bool
	upgrader              *websocket.Upgrader
}

func (s *WebsocketServer) Serve() error {
	go s.logger.Start()
	if s.tls {
		return http.ListenAndServeTLS(fmt.Sprintf(":%d", s.port), s.cert, s.key, s.router)
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}

func (s *WebsocketServer) SetSyncInterval(i uint32) {
	s.logger.SetSyncInterval(i)
}
func (s *WebsocketServer) Quit() {
	s.logger.Stop()
}

func (s *WebsocketServer) upgrade(w http.ResponseWriter, r *http.Request) {

	if uint32(s.failedConnAttempts.Get(r.RemoteAddr)) >= s.maxFailedConnAttempts {
		s.doLog("BLACKLIST", "upgrade attempt", r.RemoteAddr)
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
	s.doLog("CLIENT", "new client connected:", r.RemoteAddr)
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
	defer conn.Close()

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
			s.doLog("READ", conn.RemoteAddr().String(), err.Error(), " - disconnecting...")
			conn.Close()
			return
		} else {

			s.doLog(append([]string{"LOG", conn.RemoteAddr().String()}, msg...)...)
		}
	}
}

func (s *WebsocketServer) buildRouter() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/ws", s.upgrade)
	s.router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		Info("PING", "PONG")
		w.Write([]byte("pong"))
	})
}
