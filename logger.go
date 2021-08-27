package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"

	"strings"
)

var (
	log = logrus.New()
)

var (
	logLength *atomic.Uint32 = atomic.NewUint32(20) // default to 20
)

func SetLogLength(v uint32) {
	logLength.Store(v)
}

func cropInfo(info string) string {
	if len(info) > int(logLength.Load()) {
		return info[0:logLength.Load()]
	}
	return info
}

// Error -- log an error
func errorLog(info string, x ...string) {

	suffix := strings.Join(x, " ")
	suffix = fmt.Sprintf(stringOfLen(logLength.Load()), suffix)

	//_, caller := utils.CallInfo()	// shows who called the Error func

	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf(stringOfLen(logLength.Load()), cropInfo(info)),
	}).Error(suffix)

}

// Warn -- log an error
func warn(info string, x ...string) {

	suffix := strings.Join(x, " ")
	suffix = fmt.Sprintf(stringOfLen(logLength.Load()), suffix)

	//_, caller := utils.CallInfo()	// shows who called the Error func

	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf(stringOfLen(logLength.Load()), cropInfo(info)),
		//"Caller": caller,
		//"Time": timeStamp,
	}).Warn(suffix)
}

// Info -- log information
func info(info string, x ...string) {

	suffix := strings.Join(x, " ")
	suffix = fmt.Sprintf(stringOfLen(logLength.Load()), suffix)

	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf(stringOfLen(logLength.Load()), cropInfo(info)),
	}).Info(suffix)
}

// Info -- log information
func debug(info string, x ...string) {

	suffix := strings.Join(x, " ")
	suffix = fmt.Sprintf(stringOfLen(logLength.Load()), suffix)

	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf(stringOfLen(logLength.Load()), cropInfo(info)),
	}).Debug(suffix)
}

// Send -- Send a log to the "logger" system
func send(prefix string, name string, key string, x ...string) {
	logInfo := strings.Join(x, " ")

	if len(name) > 6 {
		name = name[0:6]
	}
	if len(key) > 7 {
		key = key[0:7]
	}
	if len(logInfo) > int(logLength.Load()) {
		logInfo = logInfo[0:logLength.Load()]
	}

	logInfo = fmt.Sprintf(stringOfLen(logLength.Load()), logInfo)

	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf(stringOfLen(logLength.Load()), strings.ToUpper(cropInfo(prefix))),
		name:     fmt.Sprintf(stringOfLen(logLength.Load()), key),
	}).Info(logInfo)
}

// Line -- logs a line
func line() {

	log.WithFields(logrus.Fields{
		"prefix": "---------",
	}).Info(blankLine)
}

type logMessage struct {
	level logLevel
	msg   []string
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

func Debug(x ...string) {
	if !DebugLogs.Load() {
		return
	}
	logChan <- logMessage{level: levelDebug, msg: x}
}
