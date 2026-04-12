package speedtestdotnet

import (
	"context"
	"errors"
	"testing"
	"time"

	provider "github.com/fopina/speedtest-cli/speedtestdotnet"
	"github.com/fopina/speedtest-cli/units"
	"github.com/spf13/cobra"
)

func TestServerIDListSetAndString(t *testing.T) {
	var ids serverIDList

	if err := ids.Set("1,20,300"); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	if got := ids.String(); got != "1,20,300" {
		t.Fatalf("String() = %q, want %q", got, "1,20,300")
	}
	if got := ids.Type(); got != "serverIDList" {
		t.Fatalf("Type() = %q, want %q", got, "serverIDList")
	}
}

func TestServerIDListSetError(t *testing.T) {
	var ids serverIDList

	if err := ids.Set("1,nope,3"); err == nil {
		t.Fatal("Set returned nil error for invalid input")
	}
}

func TestInitFlagsDefaults(t *testing.T) {
	cmd := &cobra.Command{Use: "st"}
	InitFlags(cmd)

	for _, tc := range []struct {
		name string
		want string
	}{
		{name: "bytes", want: "false"},
		{name: "list", want: "false"},
		{name: "server", want: "0"},
		{name: "time.config", want: "3"},
		{name: "time.latency", want: "3"},
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
	if got := formatSpeed("Upload speed", units.BytesPerSecond(2500)); got != "Upload speed: 20.00 Kb/s" {
		t.Fatalf("formatSpeed(bits) = %q", got)
	}

	fmtBytes = true
	if got := formatSpeed("Upload speed", units.BytesPerSecond(2500)); got != "Upload speed: 2.50 KB/s" {
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
		return units.BytesPerSecond(4321), nil
	})
	if err != nil {
		t.Fatalf("runTest returned error: %v", err)
	}
	if got != units.BytesPerSecond(4321) {
		t.Fatalf("speed = %v, want 4321 B/s", got)
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
