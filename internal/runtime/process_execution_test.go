package runtime

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/fixtures"
)

func TestRunCommandResolvesCommandsFromProvidedPath(t *testing.T) {
	tempDir := t.TempDir()
	commandPath := filepath.Join(tempDir, "test-tool")
	fixtures.WriteExecutable(
		t,
		commandPath,
		"#!/bin/sh\necho resolved\n",
	)

	output, err := RunCommand(
		tempDir,
		map[string]string{"PATH": tempDir},
		"test-tool",
	)
	if err != nil {
		t.Fatalf("RunCommand: %v", err)
	}

	if output != "resolved\n" {
		t.Fatalf("unexpected output %q", output)
	}
}
