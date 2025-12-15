package speedtestdotnet

import (
	"context"
	"fmt"
	"io"

	"github.com/fopina/speedtest-cli/prober"
	"github.com/fopina/speedtest-cli/prober/proberutil"
	"github.com/fopina/speedtest-cli/units"
)

const (
	concurrentDownloadLimit = 6
	downloadBufferSize      = 4096
	downloadRepeats         = 5
)

var downloadImageSizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}

// Will probe download speed until enough samples are taken or ctx expires.
func (s Server) ProbeDownloadSpeed(
	ctx context.Context,
	client *Client,
	stream chan<- units.BytesPerSecond,
) (units.BytesPerSecond, error) {
	grp := prober.NewGroup(concurrentDownloadLimit)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, size := range downloadImageSizes {
		for i := 0; i < downloadRepeats; i++ {
			url, err := s.RelativeURL(
				fmt.Sprintf("random%dx%d.jpg", size, size))
			if err != nil {
				return 0, fmt.Errorf("error parsing url for %v: %v", s, err)
			}

			grp.Add(func() (prober.BytesTransferred, error) {
				return client.downloadFile(ctx, url)
			})
		}
	}

	return proberutil.SpeedCollect(grp, stream)
}

func (c *Client) downloadFile(
	ctx context.Context,
	url string,
) (t prober.BytesTransferred, err error) {
	// Check early failure where context is already canceled.
	if err = ctx.Err(); err != nil {
		return
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
