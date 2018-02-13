package speedtest

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

const (
	maxUploadDuration     = maxDownloadDuration
	concurrentUploadLimit = concurrentDownloadLimit
	uploadRepeats         = downloadRepeats * 25

	safeChars = "0123456789abcdefghijklmnopqrstuv"
)

var uploadSizes = []int{1000 * 1000 / 4, 1000 * 1000 / 2}

type safeReader struct {
	in io.Reader
}

func (r safeReader) Read(p []byte) (n int, err error) {
	n, err = r.in.Read(p)

	for i := 0; i < n; i++ {
		p[i] = safeChars[p[i]&31]
	}

	return n, err
}

// Will probe upload speed until enough samples are taken or ctx expires.
func (server *Server) ProbeUploadSpeed(ctx context.Context, client *Client) (BytesPerSecond, error) {
	pg := newProberGroup(concurrentDownloadLimit)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, size := range uploadSizes {
		for i := 0; i < uploadRepeats; i++ {
			pg.Add(func(size int) func() (bytesTransferred, error) {
				return func() (bytesTransferred, error) {
					err := client.uploadFile(ctx, server.URL, size)
					if err != nil {
						return bytesTransferred(0), err
					} else {
						return bytesTransferred(size), nil
					}
				}
			}(size))
		}
	}

	return pg.Collect()
}

func (client *Client) uploadFile(ctx context.Context, url string, size int) error {
	// Check early failure where context is already canceled.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	res, err := client.post(ctx, url, "application/x-www-form-urlencoded",
		io.MultiReader(
			strings.NewReader("content1="),
			io.LimitReader(&safeReader{rand.Reader}, int64(size-9))))
	if err != nil {
		return fmt.Errorf("upload failed: %v", url, err)
	}
	defer res.Body.Close()

	return nil
}
