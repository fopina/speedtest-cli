package fastdotcom

import (
	"context"
	"testing"
)

func TestGetManifest(t *testing.T) {
	// This test might make real HTTP calls, so skip if not wanted
	t.Skip("Skipping real HTTP test for GetManifest")

	ctx, cancel := context.WithTimeout(context.Background(), 5000000000) // 5 seconds
	defer cancel()

	manifest, err := GetManifest(ctx, 1)
	if err != nil {
		t.Fatalf("GetManifest failed: %v", err)
	}
	if manifest == nil {
		t.Fatal("expected manifest, got nil")
	}
	if len(manifest.m.Targets) == 0 {
		t.Error("expected at least one target in manifest")
	}
}

func TestGetManifest_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := GetManifest(ctx, 1)
	if err == nil {
		t.Error("expected error for canceled context")
	}
}
