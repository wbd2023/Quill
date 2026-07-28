package external_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wbd2023/quill/internal/pack/external"
)

/* ------------------------------------ Executable Resolution ----------------------------------- */

func TestResolveExecutableAcceptsRelativeExecutable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(binDir, "pack-quill")
	writeExecutable(t, binary)

	resolved, err := external.ResolveExecutable(dir, "bin/pack-quill")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != binary {
		t.Fatalf("expected %q, got %q", binary, resolved)
	}
}

func TestResolveExecutableRejectsSymlinkEscape(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("symlink semantics differ on Windows")
	}

	packDir := t.TempDir()
	outside := filepath.Join(t.TempDir(), "escape")
	writeExecutable(t, outside)

	link := filepath.Join(packDir, "runner")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	if _, err := external.ResolveExecutable(packDir, "runner"); err == nil {
		t.Fatal("expected symlink escape rejection")
	}
}

func TestResolveExecutableRejectsAbsoluteCommand(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if _, err := external.ResolveExecutable(dir, "/usr/bin/evil"); err == nil {
		t.Fatal("expected absolute command rejection")
	}
}

func TestResolveExecutableRejectsParentTraversal(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	parent := filepath.Dir(dir)
	target := filepath.Join(parent, "escape")
	writeExecutable(t, target)
	t.Cleanup(func() { _ = os.Remove(target) })

	if _, err := external.ResolveExecutable(dir, "../escape"); err == nil {
		t.Fatal("expected parent traversal rejection")
	}
}

func TestResolveExecutableRejectsMissingExecutable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if _, err := external.ResolveExecutable(dir, "bin/missing"); err == nil {
		t.Fatal("expected missing executable rejection")
	}
}

func TestResolveExecutableRejectsDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := external.ResolveExecutable(dir, "bin"); err == nil {
		t.Fatal("expected directory rejection")
	}
}

func TestResolveExecutableRejectsNonExecutable(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("executable bit is not enforced on Windows")
	}

	dir := t.TempDir()
	plain := filepath.Join(dir, "notexec")
	if err := os.WriteFile(plain, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := external.ResolveExecutable(dir, "notexec"); err == nil {
		t.Fatal("expected non-executable rejection")
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func writeExecutable(tb testing.TB, path string) {
	tb.Helper()
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		tb.Fatal(err)
	}
}
