package logger

type Client interface {
	SetLogLocally(bool)
	LogLocally() bool
	Write(...string) error
	Connect() error
}

type LogServer interface {
	Serve() error
	Quit()
}

type FileLogger interface {
	SetSyncInterval(uint32)
	Start() error
	Stop()
	Write(...string)
}
