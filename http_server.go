package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	gocache "github.com/mrod502/go-cache"
	"github.com/vmihailenco/msgpack/v5"
)

type HttpServer struct {
	logger  *Log
	apiKeys *gocache.BoolCache
	notify  chan []string
}

func NewHttpServer(cfg ServerConfig) (*HttpServer, error) {

	c := make(chan []string, 1024)

	l, err := NewLog("", c)
	if err != nil {
		return nil, err
	}
	ls := &HttpServer{
		logger: l,
		notify: c,
	}
	return ls, nil
}

func (l *HttpServer) Run(path, port string) {
	Info("LOGGER", "Starting up")

	router := mux.NewRouter()

	router.HandleFunc(EndpointLog, l.doLog).Methods("POST")
	var err error
	l.logger, err = NewLog(path, l.notify)

	if err != nil {
		Error("LOG", err.Error(), "exiting")
		return
	}
	defer l.logger.Stop()
	go l.logger.run()

	go http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	closeHandler()

}

func (l *HttpServer) doLog(w http.ResponseWriter, r *http.Request) {
	var inp LogBody
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
	if !l.validKey(inp.Key) {
		return
	}
	l.logger.Write(inp.Log...)
	w.WriteHeader(http.StatusOK)
}

func openLogFile(path string) (f *os.File, err error) {
	f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		Debug("LOGFILE", "open", err.Error())
	}
	return
}

func closeHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	Info("LOGGER", "exiting")
	time.Sleep(time.Second / 2)
}

func (l *HttpServer) validKey(key string) bool {

	return l.apiKeys.Get(key)
}
