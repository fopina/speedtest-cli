package fastdotcom

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fopina/speedtest-cli/fastdotcom"
	"github.com/fopina/speedtest-cli/units"
	"github.com/spf13/cobra"
)

// Result holds the fast.com speed test results for JSON output
type Result struct {
	DownloadSpeed float64 `json:"download_speed_mbps"`
	UploadSpeed   float64 `json:"upload_speed_mbps"`
	DownloadBits  uint64  `json:"download_speed_bits_per_second"`
	UploadBits    uint64  `json:"upload_speed_bits_per_second"`
}

func Main(cmd *cobra.Command, args []string) {
	var client fastdotcom.Client

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfgTime)*time.Second)
	defer cancel()

	m, err := fastdotcom.GetManifest(ctx, urlCount)
	if err != nil {
		log.Fatalf("Error loading fast.com configuration: %v", err)
	}

	// Run download and upload tests and collect results
	var result Result
	downloadSpeed, downloadBits := runDownloadTest(m, &client)
	uploadSpeed, uploadBits := runUploadTest(m, &client)

	result.DownloadSpeed = downloadSpeed
	result.UploadSpeed = uploadSpeed
	result.DownloadBits = downloadBits
	result.UploadBits = uploadBits

	// Output results
	if jsonOut {
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Printf("Download: %.2f Mb/s\n", downloadSpeed)
		fmt.Printf("Upload: %.2f Mb/s\n", uploadSpeed)
	}
}

// Helper function to run download test and return results
func runDownloadTest(m *fastdotcom.Manifest, client *fastdotcom.Client) (float64, uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(dlTime)*time.Second)
	defer cancel()

	stream := make(chan units.BytesPerSecond, 1)

	go func() {
		speed, err := m.ProbeDownloadSpeed(ctx, client, stream)
		if err != nil {
			log.Fatalf("Error probing download speed: %v", err)
		}
		stream <- speed
	}()

	var finalSpeed units.BytesPerSecond
	select {
	case finalSpeed = <-stream:
	case <-ctx.Done():
		log.Fatalf("Download test timed out")
	}

	// Convert to Mbps and bits per second
	speedMbps := float64(finalSpeed) / 125000 // Convert bytes/s to Mbps
	bitsPerSecond := uint64(finalSpeed) * 8   // Convert bytes/s to bits/s

	return speedMbps, bitsPerSecond
}

// Helper function to run upload test and return results
func runUploadTest(m *fastdotcom.Manifest, client *fastdotcom.Client) (float64, uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ulTime)*time.Second)
	defer cancel()

	stream := make(chan units.BytesPerSecond, 1)

	go func() {
		speed, err := m.ProbeUploadSpeed(ctx, client, stream)
		if err != nil {
			log.Fatalf("Error probing upload speed: %v", err)
		}
		stream <- speed
	}()

	var finalSpeed units.BytesPerSecond
	select {
	case finalSpeed = <-stream:
	case <-ctx.Done():
		log.Fatalf("Upload test timed out")
	}

	// Convert to Mbps and bits per second
	speedMbps := float64(finalSpeed) / 125000 // Convert bytes/s to Mbps
	bitsPerSecond := uint64(finalSpeed) * 8   // Convert bytes/s to bits/s

	return speedMbps, bitsPerSecond
}
