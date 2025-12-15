package fastdotcom

import (
	"testing"
)

// Tests for GetManifest are skipped to avoid real HTTP requests to fast.com
// The function depends on internal fast.com API calls that cannot be easily mocked
// without significant refactoring.

func TestGetManifest_Skipped(t *testing.T) {
	t.Skip("Skipping all GetManifest tests to prevent real HTTP requests")
}
