package fastdotcom

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fopina/speedtest-cli/prober"
)

func TestManifest_ProbeDownloadSpeed(t *testing.T) {
	t.Skip("Skipping for now, need some mock server")
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
