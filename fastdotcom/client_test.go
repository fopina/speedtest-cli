package fastdotcom

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response content"))
	}))
	defer server.Close()

	client := &Client{}
	resp, err := client.get(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if string(body) != "response content" {
		t.Errorf("expected 'response content', got '%s'", string(body))
	}
}

func TestClient_post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/octet-stream" {
			t.Errorf("expected Content-Type 'application/octet-stream', got '%s'", r.Header.Get("Content-Type"))
		}
		// Check Content-Length header
		contentLength := r.ContentLength
		if contentLength != 100 {
			t.Errorf("expected Content-Length 100, got %d", contentLength)
		}
		// Read body and check length
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		if len(body) != 100 {
			t.Errorf("expected body length 100, got %d", len(body))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("upload ok"))
	}))
	defer server.Close()

	client := &Client{}
	resp, err := client.post(context.Background(), server.URL, 100)
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if string(body) != "upload ok" {
		t.Errorf("expected 'upload ok', got '%s'", string(body))
	}
}

func TestPutSizeIntoURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		size     int
		expected string
	}{
		{
			name:     "simple URL",
			baseURL:  "http://example.com/test",
			size:     1024,
			expected: "http://example.com/test/range/0-1024",
		},
		{
			name:     "URL with existing path",
			baseURL:  "https://fast.com/path",
			size:     0,
			expected: "https://fast.com/path/range/0-0",
		},
		{
			name:     "URL with query params",
			baseURL:  "http://test.com/path?existing=param",
			size:     500,
			expected: "http://test.com/path/range/0-500?existing=param",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := putSizeIntoURL(tt.baseURL, tt.size)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestRandomBlob(t *testing.T) {
	size := 50
	reader := randomBlob(size)
	data := make([]byte, size)

	n, err := io.ReadFull(reader, data)
	if err != nil {
		t.Fatalf("failed to read full data: %v", err)
	}
	if n != size {
		t.Errorf("expected to read %d bytes, got %d", size, n)
	}
}

func TestRandomBlob_SizeLimit(t *testing.T) {
	size := 10
	reader := randomBlob(size)
	data := make([]byte, size+10)

	n, err := reader.Read(data)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		t.Errorf("unexpected error: %v", err)
	}
	if n != size {
		t.Errorf("expected to read %d bytes, got %d", size, n)
	}
}
