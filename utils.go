package logger

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/atomic"
)

type protocol string

const (
	pWSS   protocol = "wss"
	pWS    protocol = "ws"
	pHTTPS protocol = "https"
	pHTTP  protocol = "http"

	EndpointLog string = "/log"
)

const (
	levelInfo logLevel = iota
	levelWarn
	levelErr
	levelShow
	levelVerb
	levelDebug
	levelSend
	levelLine
)
const (
	blankLine = "------------------------------------------------------------------"
)

var (
	logChan   chan logMessage
	msgInf    string
	DebugLogs *atomic.Bool
)

type logLevel byte

func stringOfLen(inp uint32) string {
	return fmt.Sprintf("%%%ds", inp)
}

func Sha256Sum(inp string) string {
	var keySum = sha256.Sum256([]byte(inp))
	return fmt.Sprintf("%x", keySum[:])
}

func closeHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	Info("LOGGER", "exiting")
	time.Sleep(time.Second / 2)
}
