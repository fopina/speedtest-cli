package fastdotcom

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fopina/speedtest-cli/fastdotcom/internal"
	"github.com/fopina/speedtest-cli/prober"
	"github.com/fopina/speedtest-cli/units"
)

func TestManifest_ProbeDownloadSpeed(t *testing.T) {
	// Skip if too long for CI
	if testing.Short() {
		t.Skip("Skipping long test in short mode")
	}

	// Create test server that serves downloads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Serve some data
		data := make([]byte, 1024)
		for i := range data {
			data[i] = byte(i % 256)
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	// Create manifest with one target
	manifest := &Manifest{
		m: &internal.Manifest{
			Targets: []internal.ManifestTarget{
				{URL: server.URL},
			},
		},
	}

	client := &Client{}
	stream := make(chan units.BytesPerSecond, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 3000000000) // 3 seconds
	defer cancel()

	speed, err := manifest.ProbeDownloadSpeed(ctx, client, stream)
	if err != nil {
		t.Fatalf("ProbeDownloadSpeed failed: %v", err)
	}
	if speed <= 0 {
		t.Errorf("expected positive speed, got %v", speed)
	}
	close(stream)
}

func TestClient_downloadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		data := []byte("test data for download")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	client := &Client{}
	url := putSizeIntoURL(server.URL, len("test data for download"))
	transferred, err := client.downloadFile(context.Background(), url)
	if err != nil {
		t.Fatalf("downloadFile failed: %v", err)
	}
	expected := prober.BytesTransferred(len("test data for download"))
	if transferred != expected {
		t.Errorf("expected %d bytes transferred, got %d", expected, transferred)
	}
}

func TestClient_downloadFile_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := &Client{}
	_, err := client.downloadFile(ctx, "http://example.com/largefile")
	if err == nil {
		t.Error("expected error for canceled context")
	}
}

func TestDownloadSizes(t *testing.T) {
	if len(downloadSizes) == 0 {
		t.Error("downloadSizes should not be empty")
	}
	for i, size := range downloadSizes {
		if size <= 0 {
			t.Errorf("downloadSizes[%d] is %d, should be positive", i, size)
		}
	}
}

func TestConcurrentDownloadLimit(t *testing.T) {
	if concurrentDownloadLimit <= 0 {
		t.Errorf("concurrentDownloadLimit is %d, should be positive", concurrentDownloadLimit)
	}
}

func TestDownloadRepeats(t *testing.T) {
	if downloadRepeats <= 0 {
		t.Errorf("downloadRepeats is %d, should be positive", downloadRepeats)
	}
}
