package logger

type ServerConfig struct {
	LogPath       string   `yaml:"log_path"`
	KeySignatures []string `yaml:"key_signatures"`
}

type ClientConfig struct {
}
