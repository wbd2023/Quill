package golang

import (
	"testing"

	"ciphera/tools/internal/toolchain"
	"ciphera/tools/internal/workspace"
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

	cmd, err := command(layout, tool, "/tool/bin:/usr/bin")
	if err != nil {
		t.Fatalf("command: %v", err)
	}

	if cmd.Name != "go" {
		t.Fatalf("Name = %q, want %q", cmd.Name, "go")
	}

	if cmd.Directory != layout.StateDirectory {
		t.Fatalf("Directory = %q, want %q", cmd.Directory, layout.StateDirectory)
	}

	if cmd.Environment["GOBIN"] != layout.BinaryDirectory() {
		t.Fatalf("GOBIN = %q, want %q",
			cmd.Environment["GOBIN"], layout.BinaryDirectory())
	}

	if cmd.Arguments[1] != "golang.org/x/tools/cmd/goimports@v0.30.0" {
		t.Fatalf("Arguments = %v, want install source@version", cmd.Arguments)
	}
}
