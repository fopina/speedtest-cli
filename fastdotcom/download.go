package fastdotcom

import (
	"context"
	"io"
	"net/url"
	"strconv"

	"go.jonnrb.io/speedtest/prober"
	"go.jonnrb.io/speedtest/prober/proberutil"
	"go.jonnrb.io/speedtest/units"
)

const (
	concurrentDownloadLimit = 12
	downloadBufferSize      = 4096
	downloadRepeats         = 5
)

var downloadSizes = []int{
	256, 512, 1024, 2048, 4096,
	131072, 1048576, 4194304, 8388608, 16777216}

// Will probe download speed until enough samples are taken or ctx expires.
func (m *Manifest) ProbeDownloadSpeed(
	ctx context.Context,
	client *Client,
	stream chan<- units.BytesPerSecond,
) (units.BytesPerSecond, error) {
	grp := prober.NewGroup(concurrentDownloadLimit)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, size := range downloadSizes {
		for i := 0; i < downloadRepeats; i++ {
			for _, t := range m.m.Targets {
				url := putSizeIntoURL(t.URL, size)
				grp.Add(func() (prober.BytesTransferred, error) {
					return client.downloadFile(ctx, url)
				})
			}
		}
	}

	return proberutil.SpeedCollect(grp, stream)
}

func putSizeIntoURL(base string, size int) string {
	u, err := url.Parse(base)
	if err != nil {
		// If we can't parse the URL, it's garbage and we'll fail to actually
		// download the URL.
		return base
	}
	if size < 0 {
		panic("fastdotcom: negative size")
	}
	u.Path += "/range/0-" + strconv.Itoa(size)
	return u.String()
}

func (c *Client) downloadFile(
	ctx context.Context,
	url string,
) (t prober.BytesTransferred, err error) {
	// Check early failure where context is already canceled.
	if err = ctx.Err(); err != nil {
		return
	}

	res, err := c.emptyPost(ctx, url)
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
