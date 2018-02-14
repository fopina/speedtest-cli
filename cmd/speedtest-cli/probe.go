package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jonnrb/speedtest"
)

func download(client *speedtest.Client, server speedtest.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), *dlTime)
	defer cancel()

	if speed, err := server.ProbeDownloadSpeed(ctx, client); err != nil {
		log.Fatalf("Error probing download speed: %v", err)
	} else {
		// Default return speed is in bytes.
		if *fmtBytes {
			fmt.Printf("Download speed: %v\n", speed)
		} else {
			fmt.Printf("Download speed: %v\n", speed.BitsPerSecond())
		}
	}
}

func upload(client *speedtest.Client, server speedtest.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), *ulTime)
	defer cancel()

	if speed, err := server.ProbeUploadSpeed(ctx, client); err != nil {
		log.Fatalf("Error probing upload speed: %v", err)
	} else {
		// Default return speed is in bytes.
		if *fmtBytes {
			fmt.Printf("Upload speed: %v\n", speed)
		} else {
			fmt.Printf("Upload speed: %v\n", speed.BitsPerSecond())
		}
	}
}
