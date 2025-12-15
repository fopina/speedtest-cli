package fastdotcom

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fopina/speedtest-cli/cmd/cli/internal"
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

	if !jsonOut {
		download(m, &client)
		upload(m, &client)
	} else {
		var result internal.Result
		dspeed, err := runTest(&client, m.ProbeDownloadSpeed)
		if err != nil {
			log.Fatalf("Error probing download speed: %v", err)
			return
		}

		uspeed, err := runTest(&client, m.ProbeUploadSpeed)
		if err != nil {
			log.Fatalf("Error probing upload speed: %v", err)
			return
		}

		result.SetSpeeds(dspeed, uspeed, fmtBytes)
		fmt.Println(result.JSON())
	}
}

func runTest(client *fastdotcom.Client, testFunc func(ctx context.Context,
	client *fastdotcom.Client,
	stream chan<- units.BytesPerSecond) (units.BytesPerSecond, error)) (units.BytesPerSecond, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(dlTime)*time.Second)
	defer cancel()
	speed, err := testFunc(ctx, client, nil)
	if err != nil {
		return 0, err
	}
	return speed, nil
}
