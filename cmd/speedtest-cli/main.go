package main

import (
	"fmt"
	"os"

	"github.com/fopina/speedtest-cli/cmd/speedtest-cli/internal/fastdotcom"
	"github.com/fopina/speedtest-cli/cmd/speedtest-cli/internal/speedtestdotnet"
	"github.com/fopina/speedtest-cli/cmd/speedtest-cli/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "speedtest-cli",
		Short: "A command-line interface for internet speed testing",
		Long: `speedtest-cli provides a command-line interface for testing internet connection speeds
using multiple backends including speedtest.net and fast.com.

If no subcommand is provided, it defaults to speedtest.net (st).`,
		Run: func(cmd *cobra.Command, args []string) {
			// Default to speedtestdotnet when no subcommand is provided
			speedtestdotnet.Main(args)
		},
	}

	// Add speedtestdotnet command
	speedtestCmd := &cobra.Command{
		Use:     "st [OPTIONS]",
		Aliases: []string{"speedtest.net"},
		Short:   "Run speed test using speedtest.net",
		Long:    "Run a speed test using the speedtest.net backend",
		Run:     func(cmd *cobra.Command, args []string) { speedtestdotnet.Main(args) },
	}
	speedtestdotnet.InitFlags(speedtestCmd)
	cmd.AddCommand(speedtestCmd)

	// Add fastdotcom command
	fastCmd := &cobra.Command{
		Use:     "f [OPTIONS]",
		Aliases: []string{"fast.com"},
		Short:   "Run speed test using fast.com",
		Long:    "Run a speed test using the fast.com backend",
		Run:     func(cmd *cobra.Command, args []string) { fastdotcom.Main(args) },
	}
	fastdotcom.InitFlags(fastCmd)
	cmd.AddCommand(fastCmd)

	// Add version command
	versionCmd := &cobra.Command{
		Use:     "v [OPTIONS]",
		Aliases: []string{"version"},
		Short:   "Show version information",
		Long:    "Display version information about speedtest-cli",
		Run:     func(cmd *cobra.Command, args []string) { version.Main(args) },
	}
	cmd.AddCommand(versionCmd)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
