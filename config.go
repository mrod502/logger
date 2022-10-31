package logger

import gocache "github.com/mrod502/go-cache"

const version = "1.1.20"

type BaseConfig struct {
	Port            uint16 `yaml:"port"`
	EnableWebsocket bool   `yaml:"enable_websocket"`
	EnableTLS       bool   `yaml:"enable_tls"`
}

type ServerConfig struct {
	BaseConfig
	LogPath       string                       `yaml:"log_path"`
	KeySignatures *gocache.Cache[string, bool] `yaml:"key_signatures"`
	CertFilePath  string                       `yaml:"cert_file_path"`
	KeyFilePath   string                       `yaml:"key_file_path"`
}

type ClientConfig struct {
	BaseConfig
	APIKey     string `yaml:"api_key"`
	Prefix     string `yaml:"prefix"`
	RemoteIP   string `yaml:"remote_ip"`
	LogLocally bool   `yaml:"log_locally"`
}

func Version() string {
	return version
}
