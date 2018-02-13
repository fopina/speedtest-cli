package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jonnrb/speedtest"
	"github.com/jonnrb/speedtest/geo"
)

func main() {
	opts := speedtest.ParseOpts()

	client, err := speedtest.NewClient(opts)
	if err != nil {
		log.Fatalf("Error parsing options: %v", err)
	}

	cfg, err := client.Config(context.TODO())
	if err != nil {
		log.Fatalf("Error loading speedtest.net configuration: %v", err)
	}

	servers, err := client.LoadAllServers(context.TODO(), cfg)
	if err != nil {
		log.Fatalf("Failed to load server list: %v\n", err)
	}
	if len(servers) == 0 {
		log.Fatalf("No servers found somehow...")
	}

	if opts.List {
		for _, s := range servers {
			fmt.Println(s)
		}
		return
	}

	fmt.Printf("Testing from %s (%s)...\n", cfg.ISP, cfg.IP)

	server, distance, latency := selectServer(opts, cfg, client, servers)
	fmt.Printf("Using server hosted by %s (%s) [%v]: %.1f ms\n",
		server.Sponsor, server.Name, distance, float64(latency)/float64(time.Millisecond))

	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()
	if speed, err := server.ProbeDownloadSpeed(ctx, client); err != nil {
		log.Fatal("Error probing download speed: %v", err)
	} else {
		fmt.Printf("Download speed: %v\n", speed)
	}

	ctx, c = context.WithTimeout(context.Background(), 5*time.Second)
	defer c()
	if speed, err := server.ProbeUploadSpeed(ctx, client); err != nil {
		log.Fatal("Error probing upload speed: %v", err)
	} else {
		fmt.Printf("Upload speed: %v\n", speed)
	}
}

func selectServer(opts *speedtest.Opts, cfg speedtest.Config, client *speedtest.Client, servers []speedtest.Server) (speedtest.Server, geo.Kilometers, time.Duration) {
	if opts.Server != 0 {
		// Meh, linear search.
		i := -1
		if i == -1 {
			log.Fatalf("Server not found: %d\n", opts.Server)
		}
		for j, s := range servers {
			if s.ID == opts.Server {
				i = j
				break
			}
		}
		selected := servers[i]

		latency, err := selected.AverageLatency(context.TODO(), client, speedtest.DefaultLatencySamples)
		if err != nil {
			log.Fatalf("Error getting latency for (%v): %v", selected, err)
		}

		distance := cfg.Coordinates.DistanceTo(selected.Coordinates)

		return selected, distance, latency
	} else {
		distanceMap := speedtest.SortServersByDistance(servers, cfg.Coordinates)
		const maxCloseServers = 5
		closestServers := func() []speedtest.Server {
			if len(servers) > maxCloseServers {
				return servers[:maxCloseServers]
			} else {
				return servers
			}
		}()

		latencyMap, err := speedtest.StableSortServersByAverageLatency(
			closestServers, context.TODO(), client, speedtest.DefaultLatencySamples)
		if err != nil {
			log.Fatalf("Error getting server latencies: %v", err)
		}

		selected := closestServers[0]

		return selected, distanceMap[selected.ID], latencyMap[selected.ID]
	}
}
