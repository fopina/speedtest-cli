package fastdotcom

import (
	"flag"
	"time"

	"github.com/spf13/cobra"
)

var (
	flagSet = flag.NewFlagSet("fastdotcom", flag.ExitOnError)

	fmtBytes = flagSet.Bool("bytes", false, "Display speeds in SI bytes (default is bits)")
	urlCount = flagSet.Int("urls", 5, "Number of URLs to use to probe")
	cfgTime  = flagSet.Duration("time.config", 1*time.Second, "Timeout for getting initial configuration")
	dlTime   = flagSet.Duration("time.download", 10*time.Second, "Maximum time to spend in download probe phase")
	ulTime   = flagSet.Duration("time.upload", 10*time.Second, "Maximum time to spend in upload probe phase")
)

// InitFlags initializes Cobra flags for the fastdotcom command
func InitFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("bytes", false, "Display speeds in SI bytes (default is bits)")
	cmd.Flags().Int("urls", 5, "Number of URLs to use to probe")
	cmd.Flags().Duration("time.config", 1*time.Second, "Timeout for getting initial configuration")
	cmd.Flags().Duration("time.download", 10*time.Second, "Maximum time to spend in download probe phase")
	cmd.Flags().Duration("time.upload", 10*time.Second, "Maximum time to spend in upload probe phase")
}
