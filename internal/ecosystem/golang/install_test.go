package golang

import (
	"testing"

	"github.com/wbd2023/quill/internal/toolchain"
	"github.com/wbd2023/quill/internal/workspace"
)

func TestCommandBuildsGoInstallRequest(t *testing.T) {
	layout := workspace.NewLayout("/repo")
	tool := toolchain.Tool{
		ID:            "goimports",
		Name:          "goimports",
		Command:       "goimports",
		PinnedVersion: "v0.30.0",
		Install:       toolchain.GoInstall{Source: "golang.org/x/tools/cmd/goimports"},
	}

	const bootstrapPath = "/tool/bin:/usr/bin"
	const resolvedGo = "/tool/bin/go"

	cmd, err := command(layout, tool, resolvedGo, bootstrapPath)
	if err != nil {
		t.Fatalf("command: %v", err)
	}

	if cmd.Name != resolvedGo {
		t.Fatalf("Name = %q, want absolute bootstrap-resolved go %q", cmd.Name, resolvedGo)
	}

	if cmd.Variables["PATH"] != bootstrapPath {
		t.Fatalf("PATH = %q, want bootstrap PATH %q", cmd.Variables["PATH"], bootstrapPath)
	}

	if cmd.Directory != layout.StateDirectory {
		t.Fatalf("Directory = %q, want %q", cmd.Directory, layout.StateDirectory)
	}

	if cmd.Variables["GOBIN"] != layout.BinaryDirectory() {
		t.Fatalf("GOBIN = %q, want %q",
			cmd.Variables["GOBIN"], layout.BinaryDirectory())
	}

	if cmd.Arguments[1] != "golang.org/x/tools/cmd/goimports@v0.30.0" {
		t.Fatalf("Arguments = %v, want install source@version", cmd.Arguments)
	}
}
