package golang

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/wbd2023/quill/internal/toolchain"
	"github.com/wbd2023/quill/internal/workspace"
)

// fakeInspectionRunner satisfies toolchain.CommandRunner for the installed-tool probe: ResolvePath
// returns the candidate verbatim and Run reports a configured version, so IsInstalled can report a
// valid installed tool without a real goimports binary on the host.
type fakeInspectionRunner struct {
	version string
}

func (fakeInspectionRunner) ResolvePath(
	_ context.Context,
	_ map[string]string,
	command string,
) (string, error) {
	return command, nil
}

func (runner fakeInspectionRunner) Run(
	_ context.Context,
	_ map[string]string,
	_ string,
	_ []string,
) (string, error) {
	return runner.version, nil
}

// TestInstallSkipsBootstrapWhenToolAlreadyInstalled is the QUILL-TRUST-006 early-return
// regression: a tool that is already installed at its pinned version returns before the host go
// executable is resolved or run. The bootstrap PATH deliberately has no resolvable go, so any
// attempt to bootstrap would fail; a nil result with no progress output proves the bootstrap was
// skipped.
func TestInstallSkipsBootstrapWhenToolAlreadyInstalled(t *testing.T) {
	layout := workspace.NewLayout(t.TempDir())
	if err := os.MkdirAll(layout.BinaryDirectory(), 0o755); err != nil {
		t.Fatalf("mkdir binary directory: %v", err)
	}
	binary := BinaryPath(layout, "goimports")
	if err := os.WriteFile(binary, []byte("installed"), 0o755); err != nil {
		t.Fatalf("write installed binary: %v", err)
	}

	tool := toolchain.Tool{
		ID:            "goimports",
		Name:          "goimports",
		Command:       "goimports",
		PinnedVersion: "1.2.3",
		Install:       toolchain.GoInstall{Source: "example.com/goimports"},
		Version:       toolchain.DetectByCommand("--version", toolchain.ExtractFirstToken),
	}

	var output bytes.Buffer
	err := Install(
		context.Background(),
		layout,
		&output,
		tool,
		"",
		"/nonexistent-quill-bootstrap-path",
		fakeInspectionRunner{version: "1.2.3"},
	)
	if err != nil {
		t.Fatalf("Install = %v, want nil for installed tool", err)
	}
	if output.Len() != 0 {
		t.Fatalf("Install wrote %q, want no progress output for installed tool", output.String())
	}
}
