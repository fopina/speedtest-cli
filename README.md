# speedtest [![Build Status](https://drone.jonnrb.com/api/badges/jon/speedtest/status.svg?branch=master)](https://drone.jonnrb.com/jon/speedtest) [![codecov](https://codecov.io/gh/jonnrb/speedtest/branch/master/graph/badge.svg)](https://codecov.io/gh/jonnrb/speedtest) [![GoDoc](https://godoc.org/github.com/jonnrb/speedtest?status.svg)](https://godoc.org/github.com/jonnrb/speedtest)

Provides a golang interface to the speedtest.net API. Comes with a very basic
CLI to test your internet speed (*now with live results*).

## Example usage

Assuming you have your `$GOPATH` set up correctly,

```
$ go get -u go.jonnrb.io/speedtest/cmd/speedtest-cli
$ speedtest-cli
Testing from LogicWeb (173.239.220.140)...
Using server hosted by Verizon (Branchburg, NJ) [20.22 km]: 72.5 ms
Download speed: 56.47 Mb/s
Upload speed: 56.60 Mb/s
```
