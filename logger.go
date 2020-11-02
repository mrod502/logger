// Package logger logs all messages through a single, clean interface
package logger

import (
	"fmt"

	prefixed "github.com/chappjc/logrus-prefix"

	"github.com/sirupsen/logrus"

	"os"
	"strings"
)

var (
	log = logrus.New()
)

const (
	logLength = 52
)

func init() {

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.DisableSorting = true
	formatter.TimestampFormat = "01-02 15:04:05"

	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "cyan+b",
		TimestampStyle: "white+h",
	})

	log.Formatter = formatter
	log.Level = logrus.DebugLevel

	log.SetOutput(os.Stdout)
}

// LogLevel -- change LogLevel on-the-fly
func LogLevel(lvl string) logrus.Level {
	switch lvl {
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}

// Error -- log an error
func errorLog(info string, x ...string) {

	suffix := strings.Join(x, " ")
	suffix = fmt.Sprintf("%-56s", suffix) // replce with length at some point

	//_, caller := utils.CallInfo()	// shows who called the Error func

	if len(info) > 9 {
		info = info[0:9]
	}
	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf("%9s", info),
		//"Caller": caller,
		//"Time": timeStamp,
	}).Error(suffix)

}

// Warn -- log an error
func warn(info string, x ...string) {

	suffix := strings.Join(x, " ")
	suffix = fmt.Sprintf("%-55s", suffix) // replce with length at some point

	//_, caller := utils.CallInfo()	// shows who called the Error func

	if len(info) > 9 {
		info = info[0:9]
	}
	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf("%9s", info),
		//"Caller": caller,
		//"Time": timeStamp,
	}).Warn(suffix)
}

// Info -- log information
func info(info string, x ...string) {

	suffix := strings.Join(x, " ")
	suffix = fmt.Sprintf("%-55s", suffix) // replce with length at some point

	if len(info) > 9 {
		info = info[0:9]
	}
	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf("%9s", info),
	}).Info(suffix)
}

// Send -- Send a log to the "logger" system
func send(prefix string, name string, key string, x ...string) {
	logInfo := strings.Join(x, " ")

	if len(prefix) > 9 {
		prefix = prefix[0:9]
	}
	if len(name) > 6 {
		name = name[0:6]
	}
	if len(key) > 7 {
		key = key[0:7]
	}
	if len(logInfo) > logLength {
		logInfo = logInfo[0:logLength]
	}

	logInfo = fmt.Sprintf("%-56s", logInfo) // replce with length at some point

	log.WithFields(logrus.Fields{
		"prefix": fmt.Sprintf("%9s", strings.ToUpper(prefix)),
		name:     fmt.Sprintf("%-6s", key),
	}).Info(logInfo)
}
