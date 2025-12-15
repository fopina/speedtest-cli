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
	ServerID       uint64  `json:"server_id"`
	ServerName     string  `json:"server_name"`
	ServerSponsor  string  `json:"server_sponsor"`
	Latency        float64 `json:"latency_ms"`
	DownloadSpeed  uint64  `json:"download_speed"`
	UploadSpeed    uint64  `json:"upload_speed"`
	DownloadPretty string  `json:"download_speed_pretty"`
	UploadPretty   string  `json:"upload_speed_pretty"`
	ISP            string  `json:"isp"`
	IP             string  `json:"ip"`
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

	log.Printf("Testing from %s (%s)...\n", cfg.ISP, cfg.IP)

	servers := listServers(ctx, &client)
	server := selectServer(&client, cfg, servers)

	if !jsonOut {
		download(&client, server)
		upload(&client, server)
	} else {
		// Collect results for JSON output
		var result Result
		result.ServerID = uint64(server.ID)
		result.ServerName = server.Name
		result.ServerSponsor = server.Sponsor
		result.Latency = float64(server.LastLatency) / float64(time.Millisecond)
		result.ISP = cfg.ISP
		result.IP = cfg.IP

		// Run download and upload tests
		dctx, cancel := context.WithTimeout(context.Background(), time.Duration(dlTime)*time.Second)
		defer cancel()
		dspeed, err := server.ProbeDownloadSpeed(dctx, &client, nil)
		if err != nil {
			log.Fatalf("Error probing download speed: %v", err)
			return
		}

		uctx, cancel := context.WithTimeout(context.Background(), time.Duration(dlTime)*time.Second)
		defer cancel()
		uspeed, err := server.ProbeUploadSpeed(uctx, &client, nil)
		if err != nil {
			log.Fatalf("Error probing upload speed: %v", err)
			return
		}

		if fmtBytes {
			result.DownloadSpeed = uint64(dspeed)
			result.UploadSpeed = uint64(uspeed)
			result.DownloadPretty = fmt.Sprintf("%v", dspeed)
			result.UploadPretty = fmt.Sprintf("%v", uspeed)
		} else {
			result.DownloadSpeed = uint64(dspeed.BitsPerSecond())
			result.UploadSpeed = uint64(uspeed.BitsPerSecond())
			result.DownloadPretty = fmt.Sprintf("%v", dspeed.BitsPerSecond())
			result.UploadPretty = fmt.Sprintf("%v", uspeed.BitsPerSecond())
		}

		// Output results
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}
		fmt.Println(string(jsonData))
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
