package speedtestdotnet

import (
	"strconv"
	"strings"
	"testing"

	"github.com/fopina/speedtest-cli/geo"
)

func TestServerID(t *testing.T) {
	// ServerID is uint64
	var id ServerID = 42
	if id != 42 {
		t.Errorf("ServerID should be 42, got %v", id)
	}
}

func TestServer_String(t *testing.T) {
	server := Server{
		ID:      1,
		Name:    "Test Server",
		Country: "Test Country",
		Sponsor: "Test Sponsor",
		Host:    "test.com:8080",
		URL:     "http://test.com/",
		Coordinates: geo.Coordinates{
			Latitude:  geo.Degrees(10.0),
			Longitude: geo.Degrees(20.0),
		},
	}

	str := server.String()
	if str == "" {
		t.Error("String() should not return empty")
	}
	// Check contains key fields
	expectedParts := []string{
		strconv.Itoa(int(server.ID)),
		server.Name,
		server.Country,
		server.Sponsor,
		server.URL, // Host is not in String output
	}
	for _, part := range expectedParts {
		if len(part) > 0 && !strings.Contains(str, part) {
			t.Errorf("String() should contain %s", part)
		}
	}
}

func TestServer_RelativeURL(t *testing.T) {
	server := Server{
		URL: "http://test.com/speedtest/",
	}

	tests := []struct {
		name     string
		local    string
		expected string
		hasError bool
	}{
		{"basic local", "download.php?size=100", "http://test.com/speedtest/download.php?size=100", false},
		{"absolute local", "/download.php", "http://test.com/download.php", false},
		{"full url", "http://other.com/file", "http://other.com/file", false}, // Actually resolves to full URL
		{"empty", "", "http://test.com/speedtest/", false},                    // Resolves to base
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.RelativeURL(tt.local)
			if tt.hasError {
				if err == nil {
					t.Error("expected error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}
