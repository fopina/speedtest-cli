package speedtestdotnet

import (
	"flag"
	"strconv"
	"strings"
	"time"

	"github.com/fopina/speedtest-cli/speedtestdotnet"
	"github.com/spf13/cobra"
)

var (
	flagSet  = flag.NewFlagSet("speedtestdotnet", flag.ExitOnError)
	fmtBytes = flagSet.Bool("bytes", false, "Display speeds in SI bytes (default is bits)")
	list     = flagSet.Bool("list", false, "List the available servers and exit")
	srvID    = flagSet.Uint64("server", 0, "Override automatic server selection")
	cfgTime  = flagSet.Duration("time.config", 1*time.Second, "Timeout for getting initial configuration")
	pngTime  = flagSet.Duration("time.latency", 1*time.Second, "Timeout for latency detection phase")
	dlTime   = flagSet.Duration("time.download", 10*time.Second, "Maximum time to spend in download probe phase")
	ulTime   = flagSet.Duration("time.upload", 10*time.Second, "Maximum time to spend in upload probe phase")
)

var srvBlk serverIDList

func init() {
	flagSet.Var(&srvBlk, "server_blocklist", "CSV of server IDs to ignore")
}

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
	cmd.Flags().Bool("bytes", false, "Display speeds in SI bytes (default is bits)")
	cmd.Flags().Bool("list", false, "List the available servers and exit")
	cmd.Flags().Uint64("server", 0, "Override automatic server selection")
	cmd.Flags().Duration("time.config", 1*time.Second, "Timeout for getting initial configuration")
	cmd.Flags().Duration("time.latency", 1*time.Second, "Timeout for latency detection phase")
	cmd.Flags().Duration("time.download", 10*time.Second, "Maximum time to spend in download probe phase")
	cmd.Flags().Duration("time.upload", 10*time.Second, "Maximum time to spend in upload probe phase")
	cmd.Flags().Var(&srvBlk, "server_blocklist", "CSV of server IDs to ignore")
}
