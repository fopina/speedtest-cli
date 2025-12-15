package speedtestdotnet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fopina/speedtest-cli/speedtestdotnet"
	"github.com/fopina/speedtest-cli/units"
	"github.com/spf13/cobra"
)

// Result holds the speed test results for JSON output
type Result struct {
	ServerID      uint64  `json:"server_id"`
	ServerName    string  `json:"server_name"`
	ServerSponsor string  `json:"server_sponsor"`
	Latency       float64 `json:"latency_ms"`
	DownloadSpeed float64 `json:"download_speed_mbps"`
	UploadSpeed   float64 `json:"upload_speed_mbps"`
	DownloadBits  uint64  `json:"download_speed_bits_per_second"`
	UploadBits    uint64  `json:"upload_speed_bits_per_second"`
	ISP           string  `json:"isp"`
	IP            string  `json:"ip"`
}

func Main(cmd *cobra.Command, args []string) {
	var client speedtestdotnet.Client

	if list {
		printServers(&client)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfgTime)*time.Second)
	defer cancel()

	cfg, err := client.Config(ctx)
	if err != nil {
		log.Fatalf("Error loading speedtest.net configuration: %v", err)
	}

	if !jsonOut {
		fmt.Printf("Testing from %s (%s)...\n", cfg.ISP, cfg.IP)
	}

	servers := listServers(ctx, &client)
	server := selectServer(&client, cfg, servers)

	// Collect results for JSON output
	var result Result
	result.ServerID = uint64(server.ID)
	result.ServerName = server.Name
	result.ServerSponsor = server.Sponsor
	result.ISP = cfg.ISP
	result.IP = cfg.IP

	// Get latency (from selectServer output)
	latency := getServerLatency(&client, server)
	result.Latency = float64(latency) / float64(time.Millisecond)

	// Run download and upload tests
	downloadSpeed, downloadBits := runDownloadTest(&client, server)
	uploadSpeed, uploadBits := runUploadTest(&client, server)

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

// Helper function to get server latency
func getServerLatency(client *speedtestdotnet.Client, server speedtestdotnet.Server) time.Duration {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pngTime)*time.Second)
	defer cancel()

	latency, err := server.AverageLatency(ctx, client, speedtestdotnet.DefaultLatencySamples)
	if err != nil {
		log.Printf("Warning: Could not get latency for server %d: %v", server.ID, err)
		return 0
	}
	return latency
}

// Helper function to run download test and return results
func runDownloadTest(client *speedtestdotnet.Client, server speedtestdotnet.Server) (float64, uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(dlTime)*time.Second)
	defer cancel()

	stream := make(chan units.BytesPerSecond, 1)

	go func() {
		speed, err := server.ProbeDownloadSpeed(ctx, client, stream)
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
func runUploadTest(client *speedtestdotnet.Client, server speedtestdotnet.Server) (float64, uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ulTime)*time.Second)
	defer cancel()

	stream := make(chan units.BytesPerSecond, 1)

	go func() {
		speed, err := server.ProbeUploadSpeed(ctx, client, stream)
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
