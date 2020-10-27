package logger

import (
	"fmt"
)

var logChan chan string

func init() {
	logChan = make(chan string, 128)
	go func() {
		for {
			fmt.Println(<-logChan)
		}
	}()
}

//Log logs stuff
func Log(in ...string) {
	msg := ""
	for _, v := range in {
		msg += v + " "
	}
	logChan <- msg
}
