package main

import (
	"flag"
	"time"
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
