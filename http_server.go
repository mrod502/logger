package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	gocache "github.com/mrod502/go-cache"
	"github.com/vmihailenco/msgpack/v5"
)

type HttpServer struct {
	logger      *FileLog
	apiKeys     *gocache.Cache[bool, string]
	notify      chan []string
	enableTLS   bool
	router      *mux.Router
	certFilePah string
	keyFilePath string
	port        uint16
	path        string
}

func (l *HttpServer) Quit() {
	l.logger.Stop()
}

func (l *HttpServer) Serve() error {
	Info("LOGGER", "Starting up")
	l.router = mux.NewRouter()
	l.router.HandleFunc(EndpointLog, l.doLog).Methods("POST")
	var err error
	l.logger, err = NewFileLog(l.path, l.notify)
	if err != nil {
		Error("LOG", err.Error(), "exiting")
		return fmt.Errorf("unable to open log file: %s", err.Error())
	}
	go l.logger.Start()

	if l.enableTLS {
		return http.ListenAndServeTLS(fmt.Sprintf(":%d", l.port), l.certFilePah, l.keyFilePath, l.router)
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", l.port), l.router)
}

func (l *HttpServer) SetSyncInterval(v uint32) {
	l.logger.SetSyncInterval(v)

}

func (l *HttpServer) doLog(w http.ResponseWriter, r *http.Request) {
	if !l.authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		Warn("LOG", "unauthorized request", r.RemoteAddr)
		return
	}
	var inp []string
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Warn("READ", "unable to read body", err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = msgpack.Unmarshal(b, &inp)
	if err != nil {
		Warn("UNMARSHAL", "unable to unmarshal body", err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	l.logger.Write(inp...)
	w.WriteHeader(http.StatusOK)
}

func openLogFile(path string) (f *os.File, err error) {
	f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		Debug("LOGFILE", "open", err.Error())
	}
	return
}

func (l *HttpServer) authorized(req *http.Request) bool {
	v, _ := l.apiKeys.Get(sha256Sum(req.Header.Get("API-Key")))
	return v
}
