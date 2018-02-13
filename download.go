package speedtest

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

const concurrentDownloadLimit = 6
const maxDownloadDuration = 10 * time.Second
const downloadBufferSize = 4096
const downloadRepeats = 5

var downloadImageSizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}

type bytesTransferred int

type BytesPerSecond float64

const (
	KBps BytesPerSecond = 1000
	MBps                = 1000 * KBps
	GBps                = 1000 * MBps
)

func (s BytesPerSecond) String() string {
	if s < KBps {
		return fmt.Sprintf("%.0f B/s", s)
	} else if s < MBps {
		return fmt.Sprintf("%.02f KB/s", s/KBps)
	} else if s < GBps {
		return fmt.Sprintf("%.02f MB/s", s/MBps)
	} else {
		return fmt.Sprintf("%.02f GB/s", s/GBps)
	}
}

// Will probe download speed until enough samples are taken or ctx expires.
func (server Server) ProbeDownloadSpeed(ctx context.Context, client *Client) (BytesPerSecond, error) {
	sem := make(chan struct{}, concurrentDownloadLimit)
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

	for _, size := range downloadImageSizes {
		for i := 0; i < downloadRepeats; i++ {
			url, err := server.RelativeURL(fmt.Sprintf("random%dx%d.jpg", size, size))
			if err != nil {
				return 0, fmt.Errorf("error parsing url for %v: %v", server, err)
			}

			wg.Add(1)
			go func(url string) {
				sem <- struct{}{}

				bytes, err := client.downloadFile(ctx, url)
				if err != nil {
					errors <- err
				}
				results <- bytes

				<-sem
				wg.Done()
			}(url)
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

	return BytesPerSecond(float64(totalSize) * float64(time.Second) / float64(duration)), nil
}

func (client *Client) downloadFile(ctx context.Context, url string) (bytesTransferred, error) {
	var t bytesTransferred

	// Check early failure where context is already canceled.
	select {
	case <-ctx.Done():
		return t, ctx.Err()
	default:
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resp, err := client.Get(ctx, url)
	if err != nil {
		return t, err
	}

	go func() {
		<-ctx.Done()
		resp.Body.Close()
	}()

	var buf [downloadBufferSize]byte
	for {
		read, err := resp.Body.Read(buf[:])
		t += bytesTransferred(read)
		if err != nil {
			if err != io.EOF {
				return t, err
			}
			break
		}
	}
	return t, nil
}
