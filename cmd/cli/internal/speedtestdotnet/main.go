package speedtestdotnet

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fopina/speedtest-cli/cmd/cli/internal"
	"github.com/fopina/speedtest-cli/speedtestdotnet"
	"github.com/fopina/speedtest-cli/units"
	"github.com/spf13/cobra"
)

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
		var result internal.Result
		result.ServerID = uint64(server.ID)
		result.ServerName = server.Name
		result.ServerSponsor = server.Sponsor
		result.Latency = float64(server.LastLatency) / float64(time.Millisecond)
		result.Distance = float64(server.LastDistance)
		result.ISP = cfg.ISP
		result.IP = cfg.IP

		dspeed, err := runTest(&client, server.ProbeDownloadSpeed)
		if err != nil {
			log.Fatalf("Error probing download speed: %v", err)
			return
		}

		uspeed, err := runTest(&client, server.ProbeUploadSpeed)
		if err != nil {
			log.Fatalf("Error probing upload speed: %v", err)
			return
		}

		result.SetSpeeds(dspeed, uspeed, fmtBytes)
		fmt.Println(result.JSON())
	}
}

func runTest(client *speedtestdotnet.Client, testFunc func(ctx context.Context,
	client *speedtestdotnet.Client,
	stream chan<- units.BytesPerSecond) (units.BytesPerSecond, error)) (units.BytesPerSecond, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(dlTime)*time.Second)
	defer cancel()
	speed, err := testFunc(ctx, client, nil)
	if err != nil {
		return 0, err
	}
	return speed, nil
}
