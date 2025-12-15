package speedtestdotnet

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_get(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	client := &Client{}
	resp, err := client.get(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer resp.Body.Close()

	content, err := resp.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent failed: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("expected 'test content', got '%s'", string(content))
	}
}

func TestClient_post(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("expected content-type 'text/plain', got '%s'", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "test body" {
			t.Errorf("expected body 'test body', got '%s'", string(body))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	defer server.Close()

	client := &Client{}
	body := bytes.NewReader([]byte("test body"))
	resp, err := client.post(context.Background(), server.URL, "text/plain", body)
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	content, err := resp.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent failed: %v", err)
	}
	if string(content) != "response" {
		t.Errorf("expected 'response', got '%s'", string(content))
	}
}

func TestResponse_ReadContent(t *testing.T) {
	content := "test XML content"
	resp := &response{Body: io.NopCloser(bytes.NewReader([]byte(content)))}
	defer resp.Body.Close()

	result, err := resp.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent failed: %v", err)
	}
	if string(result) != content {
		t.Errorf("expected '%s', got '%s'", content, string(result))
	}
}

func TestResponse_ReadXML(t *testing.T) {
	xmlContent := `<server><id>1</id><name>test</name></server>`
	resp := &response{Body: io.NopCloser(bytes.NewReader([]byte(xmlContent)))}
	defer resp.Body.Close()

	var server struct {
		ID   int    `xml:"id"`
		Name string `xml:"name"`
	}
	err := resp.ReadXML(&server)
	if err != nil {
		t.Fatalf("ReadXML failed: %v", err)
	}
	if server.ID != 1 {
		t.Errorf("expected ID 1, got %d", server.ID)
	}
	if server.Name != "test" {
		t.Errorf("expected Name 'test', got '%s'", server.Name)
	}
}
