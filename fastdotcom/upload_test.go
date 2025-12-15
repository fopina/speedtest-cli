package fastdotcom

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fopina/speedtest-cli/fastdotcom/internal"
	"github.com/fopina/speedtest-cli/prober"
	"github.com/fopina/speedtest-cli/units"
)

func TestManifest_ProbeUploadSpeed(t *testing.T) {
	// Skip if too long for CI
	if testing.Short() {
		t.Skip("Skipping long test in short mode")
	}

	// Create test server that accepts uploads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Consume the body to simulate processing
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("uploaded %d bytes", len(body))))
	}))
	defer server.Close()

	// Create manifest with one target
	manifest := &Manifest{
		m: &internal.Manifest{
			Targets: []internal.ManifestTarget{
				{URL: server.URL}, // use server.URL directly, putSizeIntoURL will add path
			},
		},
	}

	client := &Client{}
	stream := make(chan units.BytesPerSecond, 10)

	ctx, cancel := context.WithTimeout(context.Background(), 2000000000) // 2 seconds
	defer cancel()

	speed, err := manifest.ProbeUploadSpeed(ctx, client, stream)
	if err != nil {
		t.Fatalf("ProbeUploadSpeed failed: %v", err)
	}
	if speed <= 0 {
		t.Errorf("expected positive speed, got %v", speed)
	}
	// Drain stream
	close(stream)
	for range stream {
	}
}

func TestClient_uploadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Check content type
		if r.Header.Get("Content-Type") != "application/octet-stream" {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}
		// Read some body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read body: %v", err), http.StatusInternalServerError)
			return
		}
		// Respond with bytes uploaded as message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("uploaded %d bytes", len(body))))
	}))
	defer server.Close()

	client := &Client{}
	size := 100
	url := putSizeIntoURL(server.URL, size)
	transferred, err := client.uploadFile(context.Background(), url, size)
	if err != nil {
		t.Fatalf("uploadFile failed: %v", err)
	}
	if transferred != prober.BytesTransferred(size) {
		t.Errorf("expected %d bytes transferred, got %d", size, transferred)
	}
}

func TestClient_uploadFile_InvalidSize(t *testing.T) {
	client := &Client{}
	_, err := client.uploadFile(context.Background(), "http://invalid-url", 0)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestClient_uploadFile_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := &Client{}
	_, err := client.uploadFile(ctx, "http://example.com", 100)
	if err == nil {
		t.Error("expected error for canceled context")
	}
}

func TestUploadSizes(t *testing.T) {
	if len(uploadSizes) == 0 {
		t.Error("uploadSizes should not be empty")
	}
	for i, size := range uploadSizes {
		if size <= 0 {
			t.Errorf("uploadSizes[%d] is %d, should be positive", i, size)
		}
	}
}

func TestConcurrentUploadLimit(t *testing.T) {
	if concurrentUploadLimit <= 0 {
		t.Errorf("concurrentUploadLimit is %d, should be positive", concurrentUploadLimit)
	}
}

func TestUploadRepeats(t *testing.T) {
	if uploadRepeats <= 0 {
		t.Errorf("uploadRepeats is %d, should be positive", uploadRepeats)
	}
}
