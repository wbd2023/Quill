package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/profile"
)

/* ------------------------------------------ Init Run ------------------------------------------ */

func TestInitCreatesProfileAndStyleGuide(t *testing.T) {
	t.Parallel()

	tool, stdout, stderr := newTestCLI()
	root := t.TempDir()

	exitCode := tool.Run(context.Background(), []string{"init", "--repository-root", root})
	if exitCode != 0 {
		t.Fatalf("expected exit 0 for init, got %d (stderr %q)", exitCode, stderr.String())
	}

	if _, err := os.Stat(filepath.Join(root, "STYLE.md")); err != nil {
		t.Fatalf("STYLE.md not written: %v", err)
	}

	profilePath := filepath.Join(root, profile.DefaultFilename)
	contents, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("quill.toml not written: %v", err)
	}

	if _, err := profile.Parse(string(contents)); err != nil {
		t.Fatalf("generated profile must parse and validate: %v", err)
	}

	if !strings.Contains(stdout.String(), "Initialised Quill") {
		t.Fatalf("expected status output, got %q", stdout.String())
	}
}

func TestInitGeneratesImmediatelyValidProfile(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	root := t.TempDir()

	if exitCode := tool.Run(context.Background(), []string{
		"init", "--repository-root", root,
	}); exitCode != 0 {
		t.Fatalf("init failed: %d (stderr %q)", exitCode, stderr.String())
	}

	// The generated repository must load, compile, and prepare through the same engine pipeline
	// other commands use. Metadata is metadata-only and launches no tool process.
	engineInstance, err := engine.New(root)
	if err != nil {
		t.Fatalf("engine.New: %v", err)
	}

	if _, err := engineInstance.Metadata(context.Background()); err != nil {
		t.Fatalf("generated profile must be self-contained and valid: %v", err)
	}
}

func TestInitRefusesToOverwriteProfile(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	root := t.TempDir()

	if err := os.WriteFile(
		filepath.Join(root, "quill.toml"), []byte("# existing"), 0o644,
	); err != nil {
		t.Fatalf("seed quill.toml: %v", err)
	}

	exitCode := tool.Run(context.Background(), []string{"init", "--repository-root", root})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 when refusing overwrite, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "refusing to overwrite") {
		t.Fatalf("expected overwrite-refusal message, got %q", stderr.String())
	}

	// The pre-existing file must be left untouched.
	contents, err := os.ReadFile(filepath.Join(root, "quill.toml"))
	if err != nil {
		t.Fatalf("read quill.toml: %v", err)
	}

	if string(contents) != "# existing" {
		t.Fatalf("existing quill.toml must not be modified, got %q", string(contents))
	}
}

func TestInitRefusesToOverwriteStyleGuide(t *testing.T) {
	t.Parallel()

	tool, _, _ := newTestCLI()
	root := t.TempDir()

	if err := os.WriteFile(
		filepath.Join(root, "STYLE.md"), []byte("# existing"), 0o644,
	); err != nil {
		t.Fatalf("seed STYLE.md: %v", err)
	}

	if exitCode := tool.Run(context.Background(), []string{
		"init", "--repository-root", root,
	}); exitCode != 1 {
		t.Fatalf("expected exit 1 when refusing overwrite, got %d", exitCode)
	}

	if _, err := os.Stat(filepath.Join(root, profile.DefaultFilename)); err == nil {
		t.Fatalf("quill.toml must not be written when STYLE.md already exists")
	}
}

func TestInitRejectsUnknownPreset(t *testing.T) {
	t.Parallel()

	tool, _, stderr := newTestCLI()
	root := t.TempDir()

	exitCode := tool.Run(context.Background(), []string{
		"init", "--preset", "bogus", "--repository-root", root,
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d for unknown preset, got %d", usageExitCode, exitCode)
	}

	if !strings.Contains(stderr.String(), "preset") {
		t.Fatalf("expected preset validation error, got %q", stderr.String())
	}

	// init must not write anything when the preset is invalid.
	if _, err := os.Stat(filepath.Join(root, "STYLE.md")); err == nil {
		t.Fatalf("STYLE.md must not be written for an invalid preset")
	}
}

func TestInitRejectsUnexpectedArgumentWithoutWritingFiles(t *testing.T) {
	t.Parallel()

	tool, stdout, stderr := newTestCLI()
	root := t.TempDir()

	exitCode := tool.Run(context.Background(), []string{
		"init", "--repository-root", root, "unexpected",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit %d, got %d", usageExitCode, exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout for invalid init, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unexpected argument") {
		t.Fatalf("expected unexpected-argument error, got %q", stderr.String())
	}

	for _, name := range []string{"STYLE.md", profile.DefaultFilename} {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			t.Fatalf("%s must not be written for invalid init", name)
		}
	}
}

func TestResolveInitTargetDefaultsToWorkingDirectory(t *testing.T) {
	t.Parallel()

	working, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}

	target, err := resolveInitTarget("")
	if err != nil {
		t.Fatalf("resolveInitTarget: %v", err)
	}

	// Both resolve to the same absolute working directory; compare after cleaning for filesystem
	// symlink differences on some platforms.
	want, err := filepath.Abs(working)
	if err != nil {
		t.Fatalf("filepath.Abs: %v", err)
	}

	if target != want {
		t.Fatalf("resolveInitTarget(\"\") = %q, want %q", target, want)
	}

	explicit, err := resolveInitTarget("sub/dir")
	if err != nil {
		t.Fatalf("resolveInitTarget explicit: %v", err)
	}

	if !filepath.IsAbs(explicit) {
		t.Fatalf("explicit target must be absolute, got %q", explicit)
	}
}

func TestInitWritesNoOutputToStdoutOnError(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	tool := New(&stdout, &bytes.Buffer{}, "test-version")
	root := t.TempDir()

	_ = os.WriteFile(filepath.Join(root, "STYLE.md"), []byte("#"), 0o644)
	_ = tool.Run(context.Background(), []string{"init", "--repository-root", root})

	if stdout.Len() != 0 {
		t.Fatalf("init must not write stdout on error, got %q", stdout.String())
	}
}

/* ---------------------------------------- Init Security --------------------------------------- */

// TestInitRejectsDanglingSymlink verifies init does not write through an attacker-controlled
// dangling symlink: a policy path that is a symlink is treated as occupied and refused.
func TestInitRejectsDanglingSymlink(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	stylePath := filepath.Join(root, "STYLE.md")
	missingTarget := filepath.Join(root, "missing")
	if err := os.Symlink(missingTarget, stylePath); err != nil {
		t.Skipf("symlink creation unsupported on this platform: %v", err)
	}

	tool, _, stderr := newTestCLI()
	exitCode := tool.Run(context.Background(), []string{"init", "--repository-root", root})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for dangling symlink, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "refusing to overwrite") {
		t.Fatalf("expected refusal message, got %q", stderr.String())
	}

	// The symlink must remain a dangling link (not resolved), the link target must not be
	// created (no write-through), and quill.toml must not be written.
	if _, err := os.Lstat(stylePath); err != nil {
		t.Fatalf("symlink must remain untouched: %v", err)
	}
	if _, err := os.Lstat(missingTarget); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("init must not create the symlink target (write-through): %v", err)
	}
	if _, err := os.Lstat(filepath.Join(root, profile.DefaultFilename)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("quill.toml must not be written, got %v", err)
	}
}
