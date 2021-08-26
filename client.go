package logger

type Client interface {
	SetLogLocally(bool)
	LogLocally() bool
	WriteLog(...string) error
	Connect() error
}

//NewClient -
func NewClient(cfg ClientConfig) (Client, error) {
	if cfg.EnableWebsocket {
		return NewWebsocketClient(cfg)
	}
	return NewHttpClient(cfg)

}
