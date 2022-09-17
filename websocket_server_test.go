package logger

import (
	"testing"

	gocache "github.com/mrod502/go-cache"
)

func TestWebsocketServerTls(t *testing.T) {
	var baseConfig = BaseConfig{
		Port:            2121,
		EnableWebsocket: true,
		EnableTLS:       true,
	}
	var sigs = make(map[string]bool)
	var validKey = "some-valid-key"
	var validKeySignature = sha256Sum(validKey)

	sigs[validKeySignature] = true

	var cfg = ServerConfig{
		BaseConfig:    baseConfig,
		KeySignatures: gocache.New(sigs),
		LogPath:       "log-test.log",
		CertFilePath:  "test.crt",
		KeyFilePath:   "test.key",
	}

	server, err := NewWebsocketServer(cfg)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		err := server.Serve()
		if err != nil {
			panic(err)
		}
	}()

	cliCfg := ClientConfig{
		BaseConfig: baseConfig,
		RemoteIP:   "localhost",
		APIKey:     validKey,
	}

	cli, err := NewWebsocketClient(cliCfg)
	if err != nil {
		t.Fatal(err)
	}
	err = cli.Connect()
	if err != nil {
		t.Fatal(err)
	}

}
