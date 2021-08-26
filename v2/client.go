package logger

type Client interface {
	SetLogLocally(bool)
	LogLocally() bool
	WriteLog(...string) error
	Connect() error
}

func NewClient(cfg ClientConfig) (Client, error) {
	if cfg.EnableWebsocket {
		return NewWebsocketClient(cfg)
	}
	return NewHttpClient(cfg)

}
