package logger

import (
	"fmt"
	"os"
	"testing"
	"time"

	gocache "github.com/mrod502/go-cache"
)

func TestWebsocketServer(t *testing.T) {
	Info("TEST", "starting test")
	var baseConfig = BaseConfig{
		Port:            3838,
		EnableWebsocket: true,
		EnableTLS:       false,
	}

	var sigs = make(map[string]bool)
	var validKey = "some-valid-key"
	var invalidKey = "some-invalid-key"
	var validKeySignature = sha256Sum(validKey)
	var invalidKeySignature = sha256Sum(invalidKey)

	sigs[validKeySignature] = true
	sigs[invalidKeySignature] = false
	var cfg = ServerConfig{
		BaseConfig:    baseConfig,
		KeySignatures: gocache.NewBoolCache(sigs),
		LogPath:       "log-test.log",
	}
	var server LogServer
	var err error
	server, err = NewLogServer(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := server.(*WebsocketServer); !ok {
		t.Fatalf("expected *WebsocketServer, got %T", server)
	}
	go server.Serve()
	defer os.Remove("log-test.log")
	defer server.Quit()
	time.Sleep(time.Second)
	var cli Client

	cliCfg := ClientConfig{
		BaseConfig: baseConfig,
		APIKey:     validKey,
		RemoteIP:   "127.0.0.1",
		LogLocally: true,
	}

	cli, err = NewClient(cliCfg)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := cli.(*WebsocketClient); !ok {
		t.Fatalf("expected *WebsocketClient, got %T", cli)
	}

	err = cli.Connect()
	if err != nil {
		t.Fatal(err)
	}

	cli.SetLogLocally(true)
	for i := 0; i < 100; i++ {
		cli.Write("HELLO", "HI", fmt.Sprintf("%x", i))
	}
	time.Sleep(time.Second)

}
