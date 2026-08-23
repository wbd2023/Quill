//go:build !windows

package external_test

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/wbd2023/quill/internal/pack/external"
)

func TestResolveExecutableRejectsNamedPipe(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pipe := filepath.Join(dir, "runner")
	if err := syscall.Mkfifo(pipe, 0o755); err != nil {
		t.Fatal(err)
	}

	if _, err := external.ResolveExecutable(dir, "runner"); err == nil {
		t.Fatal("expected named pipe rejection")
	}
}
