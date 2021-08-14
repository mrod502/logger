package logger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	gocache "github.com/mrod502/go-cache"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/atomic"
)

type Log struct {
	f            *os.File
	c            chan []string
	ctr          *atomic.Uint32
	syncInterval *atomic.Uint32
}

func (l *Log) run() {
	for {
		v := <-l.c
		Info(v...)
		s := time.Now().Format("2006-01-02 15:04:05.99") + " " + strings.Replace(strings.Join(v, " "), "\n", "\\n", -1) + "\n"

		l.f.WriteString(s)
		l.ctr.Inc()
		if l.ctr.Load() > l.syncInterval.Load() {
			l.f.Sync()
			l.ctr.Store(0)
		}
	}
}

func NewLog(filePath string, c chan []string) (*Log, error) {
	l := new(Log)
	f, err := openLogFile(filePath)
	l.f = f
	l.c = c
	l.ctr = new(atomic.Uint32)
	l.syncInterval = atomic.NewUint32(20)

	return l, err
}

func (l *Log) SetSyncInterval(i uint32) {
	l.syncInterval.Store(i)
}

func (l *Log) Stop() {
	l.f.Close()
	l.f.Sync()
}

type HttpServer struct {
	logger  *Log
	apiKeys *gocache.StringCache
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
	l.notify <- inp.Log
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

	return l.apiKeys.Exists(key)
}
