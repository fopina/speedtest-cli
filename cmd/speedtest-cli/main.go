package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/fopina/speedtest-cli/cmd/speedtest-cli/internal/fastdotcom"
	"github.com/fopina/speedtest-cli/cmd/speedtest-cli/internal/speedtestdotnet"
	"github.com/fopina/speedtest-cli/cmd/speedtest-cli/internal/version"
)

type subcmd struct {
	mainFunc func(args []string)
	aliases  []string
}

var subcmds = []subcmd{
	{
		mainFunc: speedtestdotnet.Main,
		aliases:  []string{"st", "speedtest.net"},
	},
	{
		mainFunc: fastdotcom.Main,
		aliases:  []string{"f", "fast.com"},
	},
	{
		mainFunc: version.Main,
		aliases:  []string{"v", "version"},
	},
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) < 1 {
		// Default to first subcommand ("st" - speedtest.net) when no subcommand is provided
		subcmds[0].mainFunc([]string{})
		return
	}
	s := getSubcmd()
	if s == nil {
		flag.Usage()
		os.Exit(2)
	}
	s.mainFunc(flag.Args()[1:])
}

func getSubcmd() *subcmd {
	args := flag.Args()
	for _, s := range subcmds {
		if slices.Contains(s.aliases, args[0]) {
			return &s
		}
	}
	return nil
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "USAGE\n")
	for _, s := range subcmds {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"  %s %s [OPTIONS]\n",
			os.Args[0], strings.Join(s.aliases, "|"))
	}
	fmt.Fprintf(
		flag.CommandLine.Output(),
		"`%s` is used if none is specified\n",
		subcmds[0].aliases[0],
	)
	flag.PrintDefaults()
}
