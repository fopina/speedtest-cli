package internal

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
