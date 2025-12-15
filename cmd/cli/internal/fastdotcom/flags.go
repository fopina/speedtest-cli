package fastdotcom

import (
	"github.com/spf13/cobra"
)

// Global variables to store flag values
var (
	fmtBytes bool
	urlCount int
	cfgTime  = 1  // stored as int for seconds
	dlTime   = 10 // stored as int for seconds
	ulTime   = 10 // stored as int for seconds
	jsonOut  bool
)

// InitFlags initializes Cobra flags for the fastdotcom command
func InitFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&fmtBytes, "bytes", false, "Display speeds in SI bytes (default is bits)")
	cmd.Flags().IntVar(&urlCount, "urls", 5, "Number of URLs to use to probe")
	cmd.Flags().IntVar(&cfgTime, "time.config", 1, "Timeout for getting initial configuration (seconds)")
	cmd.Flags().IntVar(&dlTime, "time.download", 10, "Maximum time to spend in download probe phase (seconds)")
	cmd.Flags().IntVar(&ulTime, "time.upload", 10, "Maximum time to spend in upload probe phase (seconds)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output results in JSON format")
}
