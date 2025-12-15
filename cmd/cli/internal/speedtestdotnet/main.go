package speedtestdotnet

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fopina/speedtest-cli/speedtestdotnet"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	// Parse Cobra flags (this is handled by Cobra automatically)
	// cmd.ParseFlags(args) // Cobra handles this automatically

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
	fmt.Printf("Testing from %s (%s)...\n", cfg.ISP, cfg.IP)
	servers := listServers(ctx, &client)

	server := selectServer(&client, cfg, servers)

	download(&client, server)
	upload(&client, server)
}
