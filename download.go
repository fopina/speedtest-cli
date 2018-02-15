package speedtest

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/jonnrb/speedtest/prober"
)

const (
	concurrentDownloadLimit = 6
	downloadBufferSize      = 4096
	downloadRepeats         = 5
)

var downloadImageSizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}

// Will probe download speed until enough samples are taken or ctx expires.
func (s Server) ProbeDownloadSpeed(ctx context.Context, client *Client, stream chan BytesPerSecond) (BytesPerSecond, error) {
	grp := prober.NewGroup(concurrentDownloadLimit)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, size := range downloadImageSizes {
		for i := 0; i < downloadRepeats; i++ {
			url, err := s.RelativeURL(fmt.Sprintf("random%dx%d.jpg", size, size))
			if err != nil {
				return 0, fmt.Errorf("error parsing url for %v: %v", s, err)
			}

			grp.Add(func(url string) func() (prober.BytesTransferred, error) {
				return func() (prober.BytesTransferred, error) {
					return client.downloadFile(ctx, url)
				}
			}(url))
		}
	}

	return speedCollect(grp, stream)
}

func speedCollect(grp *prober.Group, stream chan BytesPerSecond) (BytesPerSecond, error) {
	start := time.Now()

	if stream != nil {
		inc := grp.GetIncremental()
		go func() {
			for b := range inc {
				d := float64(time.Since(start)) / float64(time.Second)
				stream <- BytesPerSecond(float64(b) / d)
			}
			close(stream)
		}()
	}

	b, err := grp.Collect()
	if err != nil {
		return BytesPerSecond(0), err
	} else {
		d := float64(time.Since(start)) / float64(time.Second)
		return BytesPerSecond(float64(b) / d), nil
	}
}

func (c *Client) downloadFile(ctx context.Context, url string) (prober.BytesTransferred, error) {
	var t prober.BytesTransferred

	// Check early failure where context is already canceled.
	select {
	case <-ctx.Done():
		return t, ctx.Err()
	default:
	}

	res, err := c.get(ctx, url)
	if err != nil {
		return t, err
	}
	defer res.Body.Close()

	var buf [downloadBufferSize]byte
	for {
		read, err := res.Body.Read(buf[:])
		t += prober.BytesTransferred(read)
		if err != nil {
			if err != io.EOF {
				return t, err
			}
			break
		}
	}
	return t, nil
}
