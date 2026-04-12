package version

import (
	"bytes"
	"io"
	"os"
	"testing"

	rootversion "github.com/fopina/speedtest-cli/version"
)

func TestMainPrintsVersion(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	oldVersion := rootversion.Version
	rootversion.Version = "test-version"
	t.Cleanup(func() {
		rootversion.Version = oldVersion
	})

	Main(nil)

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read stdout: %v", err)
	}

	if got := buf.String(); got != "version: test-version\n" {
		t.Fatalf("stdout = %q, want %q", got, "version: test-version\n")
	}
}
