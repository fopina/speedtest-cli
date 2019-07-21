package speedtest

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newLatencyTestServer(latency time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(latency)
		if _, err := w.Write([]byte("test=test")); err != nil {
			panic(err)
		}
	}))
}

func TestStableSortServersByAverageLatency(t *testing.T) {
	const (
		sliceSize = 5
		runs      = 5
		timeScale = 10 * time.Millisecond
	)

	expected := make([]Server, sliceSize)
	for i := 0; i < sliceSize; i++ {
		ts := newLatencyTestServer(time.Duration(i+1) * timeScale)
		defer ts.Close()
		expected[i] = Server{
			ID:  ServerID(i),
			URL: ts.URL,
		}
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < runs; i++ {
		// Effectively a Fisher-Yates shuffle.
		//
		shuffled := make([]Server, sliceSize)
		perm := r.Perm(sliceSize)
		for i, j := range perm {
			shuffled[i] = expected[j]
		}

		allSame := true
		for i, s := range shuffled {
			if s != expected[i] {
				allSame = false
				break
			}
		}
		if allSame {
			t.Logf("Already in order on run %d", i)
		}

		_, err := StableSortServersByAverageLatency(shuffled, context.Background(), &Client{}, DefaultLatencySamples)
		if err != nil {
			t.Logf("Unexpected error: %v", err)
			t.Fail()
		} else {
			for j, s := range shuffled {
				if s != expected[j] {
					t.Logf("Failure on run %d at index %d", i, j)
					t.Fail()
				}
			}
		}
	}
}

func TestServer_AverageLatency(t *testing.T) {
	const expectedLatency = 10 * time.Millisecond
	ts := newLatencyTestServer(expectedLatency)
	defer ts.Close()
	s := Server{URL: ts.URL}

	latency, err := s.AverageLatency(context.Background(), &Client{}, DefaultLatencySamples)
	if err != nil {
		t.Logf("Unexpected error: %v", err)
	}
	if latency < expectedLatency {
		t.Fail()
	}
}

func TestServer_Latency(t *testing.T) {
	const expectedLatency = 10 * time.Millisecond
	ts := newLatencyTestServer(expectedLatency)
	defer ts.Close()
	s := Server{URL: ts.URL}

	latency, err := s.Latency(context.Background(), &Client{})
	if err != nil {
		t.Logf("Unexpected error: %v", err)
	}
	if latency < expectedLatency {
		t.Fail()
	}
}
