# speedtest

Provides a golang interface to the speedtest.net API. Comes with a very basic
CLI to test your internet speed (*now with live results*).

## Example usage

Assuming you have your `$GOPATH` set up correctly,

```
$ go get -u github.com/jonnrb/speedtest/cmd/speedtest-cli
$ speedtest-cli
Testing from LogicWeb (173.239.220.140)...
Using server hosted by Verizon (Branchburg, NJ) [20.22 km]: 72.5 ms
Download speed: 56.47 Mb/s
Upload speed: 56.60 Mb/s
```
