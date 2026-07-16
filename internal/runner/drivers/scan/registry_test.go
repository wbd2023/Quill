package scan

import (
	"testing"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

func TestRepositoryScanDriverRejectsMissingScanner(t *testing.T) {
	driver := repositoryScanDriver(runtimebinding.NewRepositoryScanners())
	_, err := driver(
		runner.Context{},
		style.RepositoryScanExecution{
			Scanner: "missing",
		},
		nil,
	)
	if err == nil {
		t.Fatal("expected missing scanner error")
	}
}

func TestRepositoryScannersRejectDuplicateScannerID(t *testing.T) {
	scanner := func(
		_ runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return style.ExecutionResult{}, nil
	}

	registry := runtimebinding.NewRepositoryScanners()
	registry.Add("duplicate", scanner)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate scanner ID to panic")
		}
	}()

	registry.Add("duplicate", scanner)
}
