package node

import (
	"testing"

	"github.com/wbd2023/Quill/internal/toolchain"
	"github.com/wbd2023/Quill/internal/workspace"
)

func TestCommandBuildsNpmInstallRequest(t *testing.T) {
	layout := workspace.NewLayout("/repo")
	tool := toolchain.Tool{
		ID:            "markdownlint",
		Name:          "markdownlint",
		Command:       "markdownlint",
		PinnedVersion: "0.45.0",
		Install:       toolchain.NpmInstall{Source: "markdownlint-cli"},
	}

	cmd, err := command(layout, tool, "/tool/bin:/usr/bin")
	if err != nil {
		t.Fatalf("command: %v", err)
	}

	if cmd.Name != "npm" {
		t.Fatalf("Name = %q, want %q", cmd.Name, "npm")
	}

	if cmd.Directory != InstallDirectory(layout) {
		t.Fatalf("Directory = %q, want %q", cmd.Directory, InstallDirectory(layout))
	}

	if cmd.Environment["npm_config_cache"] != CacheDirectory(layout) {
		t.Fatalf("npm_config_cache = %q, want %q",
			cmd.Environment["npm_config_cache"], CacheDirectory(layout))
	}

	found := false
	for _, arg := range cmd.Arguments {
		if arg == "markdownlint-cli@0.45.0" {
			found = true
		}
	}
	if !found {
		t.Fatalf("Arguments %v missing markdownlint-cli@0.45.0", cmd.Arguments)
	}
}
