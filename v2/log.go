package logger

import (
	"os"
	"strings"
	"time"

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

func (l *Log) Write(inp ...string) {
	if len(inp) > 0 {
		l.c <- inp
	}
}
