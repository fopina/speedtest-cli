package main

import (
	"context"
	"fmt"
	"log"

	"go.jonnrb.io/speedtest"
)

//
// Loads the list of servers and exits the program on failure.
//
func listServers(ctx context.Context, client *speedtest.Client) []speedtest.Server {
	servers, err := client.LoadAllServers(ctx)
	if err != nil {
		log.Fatalf("Failed to load server list: %v\n", err)
	}
	if len(servers) == 0 {
		log.Fatalf("No servers found somehow...")
	}
	return servers
}

//
// Iterates through the list of server and prints them out.
//
func printServers(client *speedtest.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *cfgTime)
	defer cancel()

	for _, s := range listServers(ctx, client) {
		fmt.Println(s)
	}
}
