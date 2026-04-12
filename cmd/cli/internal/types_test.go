package internal

import (
	"encoding/json"
	"testing"

	"github.com/fopina/speedtest-cli/units"
)

func TestResultSetSpeedsBytes(t *testing.T) {
	result := Result{}

	result.SetSpeeds(units.BytesPerSecond(1500), units.BytesPerSecond(2500), true)

	if result.DownloadSpeed != 1500 {
		t.Fatalf("download speed = %d, want 1500", result.DownloadSpeed)
	}
	if result.UploadSpeed != 2500 {
		t.Fatalf("upload speed = %d, want 2500", result.UploadSpeed)
	}
	if result.DownloadPretty != "1.50 KB/s" {
		t.Fatalf("download pretty = %q, want %q", result.DownloadPretty, "1.50 KB/s")
	}
	if result.UploadPretty != "2.50 KB/s" {
		t.Fatalf("upload pretty = %q, want %q", result.UploadPretty, "2.50 KB/s")
	}
}

func TestResultSetSpeedsBits(t *testing.T) {
	result := Result{}

	result.SetSpeeds(units.BytesPerSecond(1500), units.BytesPerSecond(2500), false)

	if result.DownloadSpeed != 12000 {
		t.Fatalf("download speed = %d, want 12000", result.DownloadSpeed)
	}
	if result.UploadSpeed != 20000 {
		t.Fatalf("upload speed = %d, want 20000", result.UploadSpeed)
	}
	if result.DownloadPretty != "12.00 Kb/s" {
		t.Fatalf("download pretty = %q, want %q", result.DownloadPretty, "12.00 Kb/s")
	}
	if result.UploadPretty != "20.00 Kb/s" {
		t.Fatalf("upload pretty = %q, want %q", result.UploadPretty, "20.00 Kb/s")
	}
}

func TestResultJSON(t *testing.T) {
	result := Result{
		ServerID:       123,
		ServerName:     "Test Server",
		DownloadSpeed:  1000,
		UploadSpeed:    2000,
		DownloadPretty: "1.00 Kb/s",
		UploadPretty:   "2.00 Kb/s",
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(result.JSON()), &decoded); err != nil {
		t.Fatalf("result JSON was invalid: %v", err)
	}

	if decoded["server_id"] != float64(123) {
		t.Fatalf("server_id = %v, want 123", decoded["server_id"])
	}
	if decoded["server_name"] != "Test Server" {
		t.Fatalf("server_name = %v, want Test Server", decoded["server_name"])
	}
	if decoded["download_speed_pretty"] != "1.00 Kb/s" {
		t.Fatalf("download_speed_pretty = %v, want 1.00 Kb/s", decoded["download_speed_pretty"])
	}
}
