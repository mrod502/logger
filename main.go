package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mrod502/logger/logger"
	"go.uber.org/atomic"
)

type log struct {
	f   *os.File
	c   chan []string
	ctr *atomic.Uint32
}

func (l *log) run() {
	for {
		v := <-l.c
		logger.Info(v...)
		s := time.Now().Format("2006-01-02 15:04:05.99") + " " + strings.Replace(strings.Join(v, " "), "\n", "\\n", -1) + "\n"

		l.f.WriteString(s)
		l.ctr.Inc()
		if l.ctr.Load() > 20 {
			l.f.Sync()
			l.ctr.Store(0)
		}
	}
}
func init() {
	notify = make(chan []string, 2048)
}

var notify chan []string

func newLog(filePath string, c chan []string) (*log, error) {
	l := new(log)
	f, err := openLogFile(filePath)
	l.f = f
	l.c = c
	l.ctr = new(atomic.Uint32)
	return l, err
}

func main() {
	logger.Info("LOGGER", "Starting up")
	if len(os.Args) < 3 {
		logger.Error("LOG", "must supply file path and serve port")
		time.Sleep(time.Millisecond * 100)
		return
	}

	var path string = os.Args[1]
	var port string = os.Args[2]
	router := mux.NewRouter()

	router.HandleFunc(logger.EndpointLog, doLog)

	l, err := newLog(path, notify)

	if err != nil {
		logger.Error("LOG", err.Error(), "exiting")
		return
	}
	defer l.f.Close()
	defer l.f.Sync()
	go l.run()

	go http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	closeHandler()

}

func doLog(w http.ResponseWriter, r *http.Request) {
	var inp []string
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("READ", "unable to read body", err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = json.Unmarshal(b, &inp)

	if err != nil {
		logger.Warn("UNMARSHAL", "unable to unmarshal body", err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	notify <- inp
	w.WriteHeader(http.StatusOK)
}

func openLogFile(path string) (f *os.File, err error) {
	f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Debug("LOGFILE", "open", err.Error())
	}
	return
}

func closeHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info("LOGGER", "exiting")
	time.Sleep(time.Second / 2)
}
