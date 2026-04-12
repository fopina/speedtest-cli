package fastdotcom

import (
	"context"
	"errors"
	"testing"
	"time"

	provider "github.com/fopina/speedtest-cli/fastdotcom"
	"github.com/fopina/speedtest-cli/units"
	"github.com/spf13/cobra"
)

func TestInitFlagsDefaults(t *testing.T) {
	cmd := &cobra.Command{Use: "f"}
	InitFlags(cmd)

	for _, tc := range []struct {
		name string
		want string
	}{
		{name: "bytes", want: "false"},
		{name: "urls", want: "5"},
		{name: "time.config", want: "3"},
		{name: "time.download", want: "10"},
		{name: "time.upload", want: "10"},
		{name: "json", want: "false"},
	} {
		if got := cmd.Flags().Lookup(tc.name).DefValue; got != tc.want {
			t.Fatalf("%s default = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestFormatSpeed(t *testing.T) {
	oldFmtBytes := fmtBytes
	t.Cleanup(func() {
		fmtBytes = oldFmtBytes
	})

	fmtBytes = false
	if got := formatSpeed("Download speed", units.BytesPerSecond(1500)); got != "Download speed: 12.00 Kb/s" {
		t.Fatalf("formatSpeed(bits) = %q", got)
	}

	fmtBytes = true
	if got := formatSpeed("Download speed", units.BytesPerSecond(1500)); got != "Download speed: 1.50 KB/s" {
		t.Fatalf("formatSpeed(bytes) = %q", got)
	}
}

func TestRunTestUsesConfiguredTimeout(t *testing.T) {
	oldDLTime := dlTime
	dlTime = 1
	t.Cleanup(func() {
		dlTime = oldDLTime
	})

	var sawContext context.Context
	client := &provider.Client{}
	got, err := runTest(client, func(ctx context.Context, gotClient *provider.Client, stream chan<- units.BytesPerSecond) (units.BytesPerSecond, error) {
		sawContext = ctx
		if gotClient != client {
			t.Fatalf("client pointer changed")
		}
		if stream != nil {
			t.Fatalf("stream = %v, want nil", stream)
		}
		return units.BytesPerSecond(1234), nil
	})
	if err != nil {
		t.Fatalf("runTest returned error: %v", err)
	}
	if got != units.BytesPerSecond(1234) {
		t.Fatalf("speed = %v, want 1234 B/s", got)
	}

	deadline, ok := sawContext.Deadline()
	if !ok {
		t.Fatal("expected context deadline")
	}
	remaining := time.Until(deadline)
	if remaining <= 0 || remaining > time.Second {
		t.Fatalf("deadline remaining = %v, want within (0, 1s]", remaining)
	}
}

func TestRunTestReturnsProbeError(t *testing.T) {
	wantErr := errors.New("boom")

	_, err := runTest(&provider.Client{}, func(ctx context.Context, client *provider.Client, stream chan<- units.BytesPerSecond) (units.BytesPerSecond, error) {
		return 0, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
