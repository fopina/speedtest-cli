package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestMainHelp tests the main help command
func TestMainHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name: "main_help",
			args: []string{"--help"},
			expected: []string{
				"speedtest-cli provides a command-line interface",
				"Available Commands:",
				"st", "f", "v",
				"completion", "help",
			},
		},
		{
			name: "speedtestdotnet_help",
			args: []string{"st", "--help"},
			expected: []string{
				"Run a speed test using the speedtest.net backend",
				"Aliases:",
				"st, speedtest.net",
				"--bytes",
				"--list",
				"--server",
				"--server_blocklist",
				"--time.config",
				"--time.download",
				"--time.latency",
				"--time.upload",
			},
		},
		{
			name: "fastdotcom_help",
			args: []string{"f", "--help"},
			expected: []string{
				"Run a speed test using the fast.com backend",
				"Aliases:",
				"f, fast.com",
				"--bytes",
				"--urls",
				"--time.config",
				"--time.download",
				"--time.upload",
			},
		},
		{
			name: "version_help",
			args: []string{"v", "--help"},
			expected: []string{
				"Display version information about speedtest-cli",
				"Aliases:",
				"v, version",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCommand()

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("execute failed: %v", err)
			}

			output := buf.String()

			// Check that all expected substrings are present in the output
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected substring %q not found in output.\nFull output:\n%s", expected, output)
				}
			}
		})
	}
}
