package fastdotcom

import (
	"context"
	"log"
	"time"

	"github.com/fopina/speedtest-cli/fastdotcom"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	// Parse Cobra flags (this is handled by Cobra automatically)
	// cmd.ParseFlags(args) // Cobra handles this automatically

	var client fastdotcom.Client

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfgTime)*time.Second)
	defer cancel()

	m, err := fastdotcom.GetManifest(ctx, urlCount)
	if err != nil {
		log.Fatalf("Error loading fast.com configuration: %v", err)
	}

	download(m, &client)
	upload(m, &client)
}
