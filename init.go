package logger

import (
	"os"

	prefixed "github.com/chappjc/logrus-prefix"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

func init() {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.DisableSorting = true
	formatter.TimestampFormat = "01-02 15:04:05.999"

	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "cyan+b",
		TimestampStyle: "white+h",
	})
	formatter.ForceColors = true
	formatter.ForceFormatting = true
	log.Formatter = formatter
	log.Level = logrus.DebugLevel
	log.SetOutput(os.Stdout)
	logChan = make(chan logMessage, 1024)
	DebugLogs = atomic.NewBool(false)

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
			case levelDebug:
				debug(x.msg[0], msgInf)
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
