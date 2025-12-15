package fastdotcom

import (
	"context"
	"io"

	"github.com/fopina/speedtest-cli/prober"
	"github.com/fopina/speedtest-cli/prober/proberutil"
	"github.com/fopina/speedtest-cli/units"
)

const (
	concurrentUploadLimit = 8
	uploadRepeats         = 3
)

var uploadSizes = []int{
	256, 1024, 4096,
	131_072, 1_048_576, 8_388_608, 16_777_216,
	33_554_432}

// Will probe upload speed until enough samples are taken or ctx expires.
func (m *Manifest) ProbeUploadSpeed(
	ctx context.Context,
	client *Client,
	stream chan<- units.BytesPerSecond,
) (units.BytesPerSecond, error) {
	grp := prober.NewGroup(concurrentUploadLimit)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := range uploadSizes {
		for j := 0; j < uploadRepeats; j++ {
			for _, t := range m.m.Targets {
				size := uploadSizes[i]
				url := putSizeIntoURL(t.URL, size)
				grp.Add(func() (prober.BytesTransferred, error) {
					return client.uploadFile(ctx, url, size)
				})
			}
		}
	}

	return proberutil.SpeedCollect(grp, stream)
}

func (c *Client) uploadFile(
	ctx context.Context,
	url string,
	size int,
) (t prober.BytesTransferred, err error) {
	// Check early failure where context is already canceled.
	if err = ctx.Err(); err != nil {
		return
	}

	res, err := c.post(ctx, url, size)
	if err != nil {
		return t, err
	}
	defer res.Body.Close()
	if _, err = io.ReadAll(res.Body); err != nil {
		return prober.BytesTransferred(0), err
	}
	return prober.BytesTransferred(size), nil
}
