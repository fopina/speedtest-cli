package prober

import (
	"fmt"
	"sync"
	"testing"

	"github.com/fortytw2/leaktest"
)

func TestGroup_AllSuccessful(t *testing.T) {
	const (
		concurrency = 5
		tasks       = 10
	)

	defer leaktest.Check(t)() // Check for goroutine leaks.

	grp := NewGroup(concurrency)

	for i := 0; i < tasks; i++ {
		grp.Add(func() (BytesTransferred, error) {
			return BytesTransferred(1), nil
		})
	}

	b, err := grp.Collect()
	if err != nil {
		t.Logf("Got an error: %v", err)
		t.Fail()
	}
	expected := BytesTransferred(tasks)
	if b != expected {
		t.Logf("Expected %v bytes transferred but got %v", expected, b)
		t.Fail()
	}
}

func TestGroup_AllFail(t *testing.T) {
	const (
		concurrency = 5
		tasks       = 10
	)

	defer leaktest.Check(t)() // Check for goroutine leaks.

	grp := NewGroup(concurrency)

	testErr := fmt.Errorf("test")
	for i := 0; i < tasks; i++ {
		grp.Add(func() (BytesTransferred, error) {
			return BytesTransferred(0), testErr
		})
	}

	b, err := grp.Collect()
	if err != testErr {
		t.Logf("Got unexpected error: %v", err)
		t.Fail()
	}
	expected := BytesTransferred(0)
	if b != expected {
		t.Logf("Expected %v bytes transferred but got %v", expected, b)
		t.Fail()
	}
}

func TestGroup_StreamingResults(t *testing.T) {
	const (
		concurrency = 5
		tasks       = 10
	)

	defer leaktest.Check(t)() // Check for goroutine leaks.

	grp := NewGroup(concurrency)

	for i := 0; i < tasks; i++ {
		grp.Add(func() (BytesTransferred, error) {
			return BytesTransferred(1), nil
		})
	}

	var (
		stream = grp.GetIncremental()
		mu     sync.Mutex
		done   bool
		gdone  = make(chan struct{})
	)
	go func() {
		var i BytesTransferred
		for b := range stream {
			func() {
				mu.Lock()
				defer mu.Unlock()

				if done {
					t.Fail()
					t.Log("Got streaming result after finish")
				}

				if b < i {
					t.Fail()
					t.Log("Streaming total went down")
				}
				i = b
			}()
		}
		close(gdone)
	}()

	b, err := grp.Collect()
	func() {
		mu.Lock()
		defer mu.Unlock()
		done = true
	}()
	if err != nil {
		t.Logf("Got an error: %v", err)
		t.Fail()
	}
	expected := BytesTransferred(tasks)
	if b != expected {
		t.Logf("Expected %v bytes transferred but got %v", expected, b)
		t.Fail()
	}
	<-gdone
}
