package logger

import (
	"github.com/rs/zerolog"
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
	logChan                       chan logMessage
	titleInf, titleWarn, titleErr string
	msgInf                        string
)

type logLevel byte

type logMessage struct {
	level logLevel
	msg   []string
}

func init() {

	zerolog.TimeFieldFormat = "01-02 15:04:05"
	logChan = make(chan logMessage, 1024)

	//Info Messages
	go func() {
		for {
			x := <-logChan
			if len(x.msg) == 0 && x.level != levelLine {
				continue
			}
			msgInf = ""
			if len(x.msg) > 1 {
				for _, v := range x.msg[1:] {
					msgInf += v + " "
				}
			}
			switch x.level {
			case levelInfo:
				info(x.msg[0], msgInf)
			case levelWarn:
				warn(x.msg[0], msgInf)
			case levelErr:
				errorLog(x.msg[0], msgInf)
			case levelSend:
				if len(x.msg) > 3 {
					send(x.msg[0], x.msg[1], x.msg[2], x.msg[3:]...)
				} else if len(x.msg) == 3 {
					send(x.msg[0], x.msg[1], x.msg[2])
				}
			case levelLine:
				line()
			}
		}
	}()

}

// Info -- log information
func Info(x ...string) {

	logChan <- logMessage{level: levelInfo, msg: x}

}

// Error -- log an error
func Error(x ...string) {
	logChan <- logMessage{level: levelErr, msg: x}
}

// Warn -- log an error
func Warn(x ...string) {
	logChan <- logMessage{level: levelWarn, msg: x}
}

//Line - a dashed line
func Line() {
	logChan <- logMessage{level: levelLine}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
