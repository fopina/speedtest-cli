# speedtest.net CLI

[![goreference](https://pkg.go.dev/badge/github.com/fopina/speedtest-cli.svg)](https://pkg.go.dev/github.com/fopina/speedtest-cli)
[![release](https://img.shields.io/github/v/release/fopina/speedtest-cli)](https://github.com/fopina/speedtest-cli/releases)
[![downloads](https://img.shields.io/github/downloads/fopina/speedtest-cli/total.svg)](https://github.com/fopina/speedtest-cli/releases)
[![ci](https://github.com/fopina/speedtest-cli/actions/workflows/publish-main.yml/badge.svg)](https://github.com/fopina/speedtest-cli/actions/workflows/publish-main.yml)
[![test](https://github.com/fopina/speedtest-cli/actions/workflows/test.yml/badge.svg)](https://github.com/fopina/speedtest-cli/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/fopina/speedtest-cli/graph/badge.svg)](https://codecov.io/github/fopina/speedtest-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/fopina/speedtest-cli)](https://goreportcard.com/report/github.com/fopina/speedtest-cli)

> Forked from [surol/speedtest-cli](https://github.com/surol/speedtest-cli). Huge changes also merged from [jonnrb/speedtest](https://github.com/jonnrb/speedtest).  
> Detached as it has been archived and no longer makes sense to have it as default PR target  
> *Github, please fix this to allow fork network to be kept without that bad UX*

This is a simple command line client to speedtest.net written in Go.

## Installation

Dowload a pre-built binary from [releases](https://github.com/fopina/speedtest-cli/releases) such as:

```
curl -L https://github.com/fopina/speedtest-cli/releases/download/v2.0.0/speedtest-cli_2.0.0_linux_amd64 -o /usr/local/bin/speedtest-cli
chmod a+x /usr/local/bin/speedtest-cli
```

Or build from latest source

```
go install github.com/fopina/speedtest-cli@latest
```

## Usage

```
$ speedtest-cli --help

speedtest-cli provides a command-line interface for testing internet connection speeds
using multiple backends including speedtest.net and fast.com.

If no subcommand is provided, it defaults to speedtest.net (st).

Usage:
  speedtest-cli [flags]
  speedtest-cli [command]

Available Commands:
  f           Run speed test using fast.com
  help        Help about any command
  st          Run speed test using speedtest.net
  v           Show version information
```

Without any arguments `speedtest-cli` tests the speed using [speedtest.net](speedtest.net) closest server with the lowest latency.

```
$ speedtest-cli          
2025/12/15 16:45:02 Testing from Vodafone Portugal (1.4.1.4)...
2025/12/15 16:45:02 Using server 60452 hosted by DIGI (Lisbon, Portugal) [0.45 km]: 8.7 ms
Download speed: 405.00 Mb/s
Upload speed: 199.25 Mb/s
```

[fast.com](fast.com) can also be used

```
$ speedtest-cli f
Download speed: 378.53 Mb/s
Upload speed: 128.30 Mb/s
```

`--json` works on either provider for an output easier to consume by scripts

```
$ speedtest-cli st --json
2025/12/15 16:45:02 Testing from Vodafone Portugal (1.4.1.4)...
2025/12/15 16:45:02 Using server 60452 hosted by DIGI (Lisbon, Portugal) [0.45 km]: 8.7 ms
{
  "server_id": 60452,
  "server_name": "Lisbon, Portugal",
  "server_sponsor": "DIGI",
  "latency_ms": 8.71926,
  "distance_km": 0.452,
  "download_speed": 402340398,
  "upload_speed": 74594946,
  "download_speed_pretty": "402.34 Mb/s",
  "upload_speed_pretty": "74.59 Mb/s",
  "isp": "DIGI",
  "ip": "1.4.1.4"
}
```

Check `--help` on each subcommand for more options such as `--bytes` (output speeds in bytes rather than bits)