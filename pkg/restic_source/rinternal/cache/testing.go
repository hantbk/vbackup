package cache

import (
	"testing"

	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/test"
)

// TestNewCache returns a cache in a temporary directory which is removed when
// cleanup is called.
func TestNewCache(t testing.TB) *Cache {
	dir := test.TempDir(t)
	t.Logf("created new cache at %v", dir)
	cache, err := New(restic.NewRandomID().String(), dir)
	if err != nil {
		t.Fatal(err)
	}
	return cache
}
