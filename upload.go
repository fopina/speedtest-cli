package speedtest

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
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
	sem := make(chan struct{}, concurrentUploadLimit)
	results := make(chan bytesTransferred)
	var wg sync.WaitGroup

	var lastErr error // Keep the last transfer error in case nothing works.
	errors := make(chan error)
	go func() {
		for err := range errors {
			lastErr = err
		}
	}()

	start := time.Now()

	for _, size := range uploadSizes {
		for i := 0; i < uploadRepeats; i++ {
			wg.Add(1)
			go func(url string, size int) {
				sem <- struct{}{}

				if err := client.uploadFile(ctx, url, size); err != nil {
					errors <- err
				} else {
					results <- bytesTransferred(size)
				}

				<-sem
				wg.Done()
			}(server.URL, size)
		}
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	var totalSize bytesTransferred
	for b := range results {
		totalSize += b
	}

	if totalSize == bytesTransferred(0) {
		return BytesPerSecond(0), lastErr
	}

	duration := time.Since(start)

	return BytesPerSecond(int64(totalSize) * int64(time.Second) / int64(duration)), nil
}

func (client *Client) uploadFile(ctx context.Context, url string, size int) error {
	// Check early failure where context is already canceled.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	res, err := client.Post(ctx, url, "application/x-www-form-urlencoded",
		io.MultiReader(
			strings.NewReader("content1="),
			io.LimitReader(&safeReader{rand.Reader}, int64(size-9))))
	if err != nil {
		return fmt.Errorf("upload failed: %v", url, err)
	}
	defer res.Body.Close()

	return nil
}
