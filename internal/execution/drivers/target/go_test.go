package target

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/pack/shipped/golang"
	"github.com/wbd2023/Quill/internal/pack/shipped/tool"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func TestRunGolangciRulePassesCleanRepository(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))
	testutil.WriteFile(t, repoRoot, "cmd/quill/main.go", "package main\n\nfunc main() {}\n")
	testutil.WriteFile(t, repoRoot, "internal/example/example.go", "package example\n")
	writeExecutable(t, repoRoot, "goimports")
	writeExecutable(t, repoRoot, "golangci-lint")

	runCtx := testContext(t, repoRoot, style.Scope("all"))

	job := style.TargetCommandJob{
		ToolIDs:  []string{tool.Go, tool.Goimports, tool.GolangciLint},
		Action:   golang.TargetActionGolangci,
		Language: golang.Language,
		Targets:  []string{"go"},
	}

	result, err := testTargetCommandDriver()(context.Background(), runCtx, job, nil)
	if err != nil {
		t.Fatalf("golangciDriver(all): %v", err)
	}

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected repository lint diagnostics: %#v", result.Diagnostics)
	}
}

func writeExecutable(t *testing.T, repoRoot string, name string) {
	t.Helper()

	path := testutil.WriteFile(
		t,
		repoRoot,
		filepath.Join(".cache", "quill", "bin", name),
		"#!/bin/sh\nexit 0\n",
	)
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("make %s executable: %v", name, err)
	}
}

func testTargetCommandDriver() (driver execution.Driver) {
	commands := driverkit.NewTargetCommands()
	commands.Add(
		golang.TargetActionGolangci,
		RunGolangci(golang.PackID, tool.GolangciLint, tool.Goimports, golang.Language),
	)
	return targetCommandDriver(commands)
}
