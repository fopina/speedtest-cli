package units

import (
	"strings"
	"testing"
)

func TestBytesPerSecond_BitsPerSecond(t *testing.T) {
	tests := []struct {
		name     string
		input    BytesPerSecond
		expected BitsPerSecond
	}{
		{"zero", 0, 0},
		{"one byte", 1, 8},
		{"1000 bytes", 1000, 8000},
		{"negative", -10, -80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.BitsPerSecond()
			if result != tt.expected {
				t.Errorf("expected %v bits/s, got %v", tt.expected, result)
			}
		})
	}
}

func TestBytesPerSecond_String(t *testing.T) {
	tests := []struct {
		name     string
		input    BytesPerSecond
		contains string
	}{
		{"zero", 0, "0 B/s"},
		{"bytes", 500, "500 B/s"},
		{"KB/s", 1500, "KB/s"},
		{"MB/s", 1500000, "MB/s"},
		{"GB/s", 1500000000, "GB/s"},
		{"negative", -1000, "-1000 B/s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if !strings.Contains(result, tt.contains) {
				t.Errorf("expected result to contain '%s', got '%s'", tt.contains, result)
			}
		})
	}
}

func TestBitsPerSecond_BytesPerSecond(t *testing.T) {
	tests := []struct {
		name     string
		input    BitsPerSecond
		expected BytesPerSecond
	}{
		{"zero", 0, 0},
		{"eight bits", 8, 1},
		{"1000 bits", 1000, 125},
		{"negative", -80, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.BytesPerSecond()
			if result != tt.expected {
				t.Errorf("expected %v bytes/s, got %v", tt.expected, result)
			}
		})
	}
}

func TestBitsPerSecond_String(t *testing.T) {
	tests := []struct {
		name     string
		input    BitsPerSecond
		contains string
	}{
		{"zero", 0, "0 b/s"},
		{"bits", 500, "500 b/s"},
		{"Kb/s", 1500, "Kb/s"},
		{"Mb/s", 1500000, "Mb/s"},
		{"Gb/s", 1500000000, "Gb/s"},
		{"negative", -1000, "-1000 b/s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if !strings.Contains(result, tt.contains) {
				t.Errorf("expected result to contain '%s', got '%s'", tt.contains, result)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that constants are defined correctly
	if KBps != 1000 {
		t.Errorf("expected KBps to be 1000, got %v", KBps)
	}
	if MBps != 1000000 {
		t.Errorf("expected MBps to be 1000000, got %v", MBps)
	}
	if GBps != 1000000000 {
		t.Errorf("expected GBps to be 1000000000, got %v", GBps)
	}

	if Kbps != 1000 {
		t.Errorf("expected Kbps to be 1000, got %v", Kbps)
	}
	if Mbps != 1000000 {
		t.Errorf("expected Mbps to be 1000000, got %v", Mbps)
	}
	if Gbps != 1000000000 {
		t.Errorf("expected Gbps to be 1000000000, got %v", Gbps)
	}
}
