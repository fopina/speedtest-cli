package speedtest

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

const DefaultLatencySamples = 4

// Probes each and every server and stable sorts them based on their average
// latencies. Will fail fast. Returns the average latencies in a map keyed by
// each server's ID and an error if any.
//
// Linear time constrained (len(servers)*samples) because server's are probed
// one at a time.
//
// TODO(jonnrb): Is failing fast the Right Thing To Do here? A very probable
// reason a probe failed was because the context expired, and if so, returning
// is necessary. A reasonable timeout per probe would be required.
//
func StableSortServersByAverageLatency(servers []Server, ctx context.Context, client *Client, samples int) (map[ServerID]time.Duration, error) {
	if samples <= 0 {
		return nil, fmt.Errorf("taking %v latency samples makes no sense", samples)
	}

	m := make(map[ServerID]time.Duration)
	for _, s := range servers {
		if _, ok := m[s.ID]; ok {
			continue
		}
		if latency, err := s.AverageLatency(ctx, client, samples); err != nil {
			return nil, err
		} else {
			m[s.ID] = latency
		}
	}

	sort.SliceStable(servers, func(i, j int) bool {
		return m[servers[i].ID] < m[servers[j].ID]
	})

	return m, nil
}

// Takes samples of a server's latency and returns the average.
//
// Serialized and fails fast. It is assumed that if there is an error doing a
// single latency probe to a server, that server is not a good candidate for a
// speed test.
//
func (server Server) AverageLatency(ctx context.Context, client *Client, samples int) (time.Duration, error) {
	if samples <= 0 {
		return time.Duration(0), fmt.Errorf("samples must be positive; was %v", samples)
	}

	var total time.Duration
	for i := 0; i < samples; i++ {
		if d, err := server.Latency(ctx, client); err != nil {
			return time.Duration(0), err
		} else {
			total += d
		}
	}

	return total / time.Duration(samples), nil
}

func (server *Server) Latency(ctx context.Context, client *Client) (time.Duration, error) {
	start := time.Now()

	url, err := server.RelativeURL("latency.txt")
	if err != nil {
		return time.Duration(0), fmt.Errorf("could not parse realtive path to latency.txt: %v", err)
	}

	resp, err := client.get(ctx, url)

	if resp != nil {
		url = resp.Request.URL.String()
	}

	if err != nil {
		return time.Duration(0), fmt.Errorf("[%s] Failed to detect latency: %v\n", url, err)
	}
	if resp.StatusCode != 200 {
		return time.Duration(0), fmt.Errorf("[%s] Invalid latency detection HTTP status: %d\n", url, resp.StatusCode)
	}

	content, err := resp.ReadContent()
	d := time.Since(start)

	if err != nil {
		return time.Duration(0), fmt.Errorf("[%s] Failed to read latency response: %v\n", url, err)
	}
	if !strings.HasPrefix(string(content), "test=test") {
		return time.Duration(0), fmt.Errorf("[%s] Invalid latency response: %s\n", url, content)
	}

	return d, nil
}
