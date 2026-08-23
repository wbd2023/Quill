package installer

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wbd2023/quill/internal/lockfile"
	"github.com/wbd2023/quill/internal/toolchain"
	"github.com/wbd2023/quill/internal/workspace"
)

func TestInstallStopsBeforeStartingSecondToolAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firstErr := errors.New("first failure")

	var calls []string
	err := installWith(
		ctx,
		workspace.Layout{},
		io.Discard,
		[]toolchain.Tool{{ID: "first"}, {ID: "later"}},
		lockfile.Lockfile{},
		func(
			_ context.Context,
			_ workspace.Layout,
			_ io.Writer,
			tool toolchain.Tool,
			_ lockfile.Lockfile,
		) error {
			calls = append(calls, tool.ID)
			cancel()
			return firstErr
		},
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Install error = %v, want context.Canceled", err)
	}
	if !errors.Is(err, firstErr) {
		t.Fatalf("Install error = %v, want first failure", err)
	}
	if len(calls) != 1 || calls[0] != "first" {
		t.Fatalf("Install calls = %v, want [first]", calls)
	}
}

// TestInstallToolSkipsBootstrapPathForManagedTool proves the outer installer performs the
// installed-tool probe before it requires a trusted host bootstrap PATH.
func TestInstallToolSkipsBootstrapPathForManagedTool(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fixture is a POSIX executable script")
	}

	layout := workspace.NewLayout(t.TempDir())
	if err := os.MkdirAll(layout.BinaryDirectory(), 0o755); err != nil {
		t.Fatalf("mkdir managed binary directory: %v", err)
	}
	binary := filepath.Join(layout.BinaryDirectory(), "goimports")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\necho 1.2.3\n"), 0o755); err != nil {
		t.Fatalf("write managed binary: %v", err)
	}
	t.Setenv("PATH", layout.StateDirectory)

	tool := toolchain.Tool{
		ID:            "goimports",
		Name:          "goimports",
		Command:       "goimports",
		PinnedVersion: "1.2.3",
		Install:       toolchain.GoInstall{Source: "example.com/goimports"},
		Version:       toolchain.DetectByCommand("--version", toolchain.ExtractFirstToken),
	}
	if err := installTool(context.Background(), layout, io.Discard, tool, lockfile.Lockfile{}); err != nil {
		t.Fatalf("installTool: %v", err)
	}
}
