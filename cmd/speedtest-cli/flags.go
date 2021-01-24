package main

import (
	"flag"
	"strconv"
	"strings"
	"time"

	"go.jonnrb.io/speedtest/speedtestdotnet"
)

var (
	fmtBytes = flag.Bool("bytes", false, "Display speeds in SI bytes (default is bits)")
	list     = flag.Bool("list", false, "List the available servers and exit")
	srvID    = flag.Uint64("server", 0, "Override automatic server selection")
	cfgTime  = flag.Duration("time.config", 1*time.Second, "Timeout for getting initial configuration")
	pngTime  = flag.Duration("time.latency", 1*time.Second, "Timeout for latency detection phase")
	dlTime   = flag.Duration("time.download", 10*time.Second, "Maximum time to spend in download probe phase")
	ulTime   = flag.Duration("time.upload", 10*time.Second, "Maximum time to spend in upload probe phase")
)

var srvBlk serverIDList

func init() {
	flag.Var(&srvBlk, "server_blocklist", "CSV of server IDs to ignore")
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
