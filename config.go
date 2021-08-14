package logger

type ServerConfig struct {
	LogPath       string   `yaml:"log_path"`
	KeySignatures []string `yaml:"key_signatures"`
}

type ClientConfig struct {
	EnableWebsocket bool
	EnableTLS       bool
	APIKey          string
	Prefix          string
	RemoteIP        string
	LogLocally      bool
	Port            uint16
}
