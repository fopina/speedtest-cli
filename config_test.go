package speedtest

import (
	"context"
	"testing"
	"time"
)

// Dumb test that just prints out the config downloaded from speedtest.net
//
func TestClient_Config(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var client Client
	cfg, err := client.Config(ctx)
	if err != nil {
		t.Fail()
	}

	t.Log(cfg)
}
