package speedtestdotnet

import (
	"strconv"
	"strings"

	"github.com/fopina/speedtest-cli/speedtestdotnet"
	"github.com/spf13/cobra"
)

// Global variables to store flag values
var (
	fmtBytes bool
	list     bool
	srvID    uint64
	cfgTime  = 1  // stored as int for seconds
	pngTime  = 1  // stored as int for seconds
	dlTime   = 10 // stored as int for seconds
	ulTime   = 10 // stored as int for seconds
	srvBlk   serverIDList
)

type serverIDList []speedtestdotnet.ServerID

func (l *serverIDList) Set(s string) (err error) {
	sl := strings.Split(s, ",")
	*l = make(serverIDList, len(sl))
	for i, s := range sl {
		var n int
		n, err = strconv.Atoi(s)
		(*l)[i] = speedtestdotnet.ServerID(n)
		if err != nil {
			return
		}
	}
	return
}

func (l *serverIDList) String() string {
	sl := make([]string, len(*l))[:0]
	for i, j := range *l {
		sl[i] = strconv.Itoa(int(j))
	}
	return strings.Join(sl, ",")
}

func (l *serverIDList) Type() string {
	return "serverIDList"
}

// InitFlags initializes Cobra flags for the speedtestdotnet command
func InitFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&fmtBytes, "bytes", false, "Display speeds in SI bytes (default is bits)")
	cmd.Flags().BoolVar(&list, "list", false, "List the available servers and exit")
	cmd.Flags().Uint64Var(&srvID, "server", 0, "Override automatic server selection")
	cmd.Flags().IntVar(&cfgTime, "time.config", 1, "Timeout for getting initial configuration (seconds)")
	cmd.Flags().IntVar(&pngTime, "time.latency", 1, "Timeout for latency detection phase (seconds)")
	cmd.Flags().IntVar(&dlTime, "time.download", 10, "Maximum time to spend in download probe phase (seconds)")
	cmd.Flags().IntVar(&ulTime, "time.upload", 10, "Maximum time to spend in upload probe phase (seconds)")
	cmd.Flags().Var(&srvBlk, "server_blocklist", "CSV of server IDs to ignore")
}
