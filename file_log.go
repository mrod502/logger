package logger

import (
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"
)

type FileLog struct {
	f            *os.File
	c            chan []string
	ctr          *atomic.Uint32
	syncInterval *atomic.Uint32
	queueEmpty   *sync.WaitGroup
}

func (l *FileLog) Start() {
	for {
		v := <-l.c
		Info(v...)
		s := time.Now().Format("2006-01-02 15:04:05.999") + " " + strings.Replace(strings.Join(v, " "), "\n", "\\n", -1) + "\n"
		l.f.WriteString(s)
		l.ctr.Inc()
		if l.ctr.Load() > l.syncInterval.Load() {
			l.f.Sync()
			l.ctr.Store(0)
		}
		l.queueEmpty.Done()
	}
}

func (l *FileLog) SetSyncInterval(i uint32) {
	l.syncInterval.Store(i)
}

func (l *FileLog) Stop() {
	l.queueEmpty.Wait()
	l.f.Sync()
	l.f.Close()
}

func (l *FileLog) Write(inp ...string) {
	if len(inp) > 0 {
		l.queueEmpty.Add(1)
		l.c <- inp
	}
}
