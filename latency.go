package speedtest

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const DefaultLatencySamples = 4

// Probes each and every server and stable sorts them based on their average
// latencies.
//
// Returns the average latencies in a map keyed by each server's ID and an error
// if all probes were unsuccessful. An error for a particular server gets
// treated as if the latency was higher than any successful probe.
//
func StableSortServersByAverageLatency(
	servers []Server,
	ctx context.Context,
	client *Client,
	samples int,
) (map[ServerID]time.Duration, error) {
	if samples <= 0 {
		return nil, fmt.Errorf("taking %v latency samples makes no sense", samples)
	}

	m, err := measureAllAverageLatencies(servers, ctx, client, samples)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(servers, func(i, j int) bool {
		return m[servers[i].ID] < m[servers[j].ID]
	})

	return m, nil
}

func measureAllAverageLatencies(
	servers []Server,
	ctx context.Context,
	client *Client,
	samples int,
) (m map[ServerID]time.Duration, err error) {
	dMax := maxLatencyFor(ctx)

	c := fanOutLatencyProbes(ctx, servers, client, samples)
	m = make(map[ServerID]time.Duration)
	var anyGood bool
	for r := range c {
		if r.err != nil {
			err = r.err
			r.d = dMax
		} else {
			anyGood = true
		}
		m[r.s] = r.d
	}
	if anyGood {
		err = nil
	}
	return
}

type latencyProbeRes struct {
	s   ServerID
	d   time.Duration
	err error
}

func fanOutLatencyProbes(
	ctx context.Context,
	servers []Server,
	client *Client,
	samples int,
) <-chan latencyProbeRes {
	c := make(chan latencyProbeRes)
	var g sync.WaitGroup
	g.Add(len(servers))
	go func() {
		g.Wait()
		close(c)
	}()

	for i := range servers {
		s := servers[i]
		go func() {
			d, err := s.AverageLatency(ctx, client, samples)
			c <- latencyProbeRes{s.ID, d, err}
			g.Done()
		}()
	}
	return c
}

func maxLatencyFor(ctx context.Context) time.Duration {
	t, ok := ctx.Deadline()
	if ok {
		return t.Sub(time.Now())
	}
	// Return some silly-high sentinel value. It just has to be larger than any
	// sensible result.
	return 24 * time.Hour
}

// Takes samples of a server's latency and returns the average.
//
// Serialized and fails fast. It is assumed that if there is an error doing a
// single latency probe to a server, that server is not a good candidate for a
// speed test.
//
func (s Server) AverageLatency(
	ctx context.Context,
	client *Client,
	samples int,
) (time.Duration, error) {
	if samples <= 0 {
		panic("must have samples > 0")
	}

	var total time.Duration
	for i := 0; i < samples; i++ {
		if d, err := s.Latency(ctx, client); err != nil {
			return time.Duration(0), err
		} else {
			total += d
		}
	}

	return total / time.Duration(samples), nil
}

func (s Server) Latency(
	ctx context.Context,
	client *Client,
) (time.Duration, error) {
	start := time.Now()
	err := s.download(ctx, client)
	return time.Since(start), err
}

func (s Server) download(ctx context.Context, client *Client) error {
	url, err := s.RelativeURL("latency.txt")
	if err != nil {
		return fmt.Errorf("could not parse realtive path to latency.txt: %v", err)
	}

	res, err := client.get(ctx, url)
	if res != nil {
		url = res.Request.URL.String()
	}
	if err != nil {
		return fmt.Errorf("[%s] Failed to detect latency: %v\n", url, err)
	}

	return checkDownloadResponse(res)
}

func checkDownloadResponse(res *response) error {
	url := res.Request.URL.String()
	if res.StatusCode != 200 {
		return fmt.Errorf(
			"[%s] Invalid latency detection HTTP status: %d\n",
			url, res.StatusCode)
	}

	content, err := res.ReadContent()
	if err != nil {
		return fmt.Errorf(
			"[%s] Failed to read latency response: %v\n",
			url, err)
	}

	if !strings.HasPrefix(string(content), "test=test") {
		return fmt.Errorf("[%s] Invalid latency response: %s\n", url, content)
	}
	return nil
}
