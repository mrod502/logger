package logger

type LogServer interface {
	Serve() error
	SetSyncInterval(i uint32)
	Quit()
}

func NewLogServer(cfg ServerConfig) (l LogServer, err error) {
	if cfg.EnableWebsocket {
		return NewWebsocketServer(cfg)
	}
	return NewHttpServer(cfg)
}
