package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
	rtest "github.com/hantbk/vbackup/pkg/restic_source/rinternal/test"

	"github.com/cenkalti/backoff/v4"
)

func TestNoSpacePermanent(t *testing.T) {
	oldTempFile := tempFile
	defer func() {
		tempFile = oldTempFile
	}()

	tempFile = func(_, _ string) (*os.File, error) {
		return nil, fmt.Errorf("not creating tempfile, %w", syscall.ENOSPC)
	}

	dir := rtest.TempDir(t)

	be, err := Open(context.Background(), Config{Path: dir, Connections: 2})
	rtest.OK(t, err)
	defer func() {
		rtest.OK(t, be.Close())
	}()

	h := restic.Handle{Type: restic.ConfigFile}
	err = be.Save(context.Background(), h, nil)
	_, ok := err.(*backoff.PermanentError)
	rtest.Assert(t, ok,
		"error type should be backoff.PermanentError, got %T", err)
	rtest.Assert(t, errors.Is(err, syscall.ENOSPC),
		"could not recover original ENOSPC error")
}
