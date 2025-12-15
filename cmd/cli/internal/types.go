package internal

import (
	"encoding/json"
	"log"

	"github.com/fopina/speedtest-cli/units"
)

// Result holds the speed test results for JSON output
type Result struct {
	ServerID       uint64  `json:"server_id,omitempty"`
	ServerName     string  `json:"server_name,omitempty"`
	ServerSponsor  string  `json:"server_sponsor,omitempty"`
	Latency        float64 `json:"latency_ms,omitempty"`
	Distance       float64 `json:"distance_km,omitempty"`
	DownloadSpeed  uint64  `json:"download_speed"`
	UploadSpeed    uint64  `json:"upload_speed"`
	DownloadPretty string  `json:"download_speed_pretty"`
	UploadPretty   string  `json:"upload_speed_pretty"`
	ISP            string  `json:"isp,omitempty"`
	IP             string  `json:"ip,omitempty"`
}

func (result *Result) SetSpeeds(downSpeed, upSpeed units.BytesPerSecond, formatBytes bool) {
	if formatBytes {
		result.DownloadSpeed = uint64(downSpeed)
		result.UploadSpeed = uint64(upSpeed)
		result.DownloadPretty = downSpeed.String()
		result.UploadPretty = upSpeed.String()
	} else {
		result.DownloadSpeed = uint64(downSpeed.BitsPerSecond())
		result.UploadSpeed = uint64(upSpeed.BitsPerSecond())
		result.DownloadPretty = downSpeed.BitsPerSecond().String()
		result.UploadPretty = upSpeed.BitsPerSecond().String()
	}
}

func (result *Result) JSON() string {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
		return "<error>"
	}
	return string(jsonData)
}
