package scan

import (
	"testing"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
)

func TestRepositoryScanDriverRejectsMissingScanner(t *testing.T) {
	driver := repositoryScanDriver(driverkit.NewRepositoryScanners())
	_, err := driver(
		execution.RunContext{},
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
		_ execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return style.ExecutionResult{}, nil
	}

	registry := driverkit.NewRepositoryScanners()
	registry.Add("duplicate", scanner)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate scanner ID to panic")
		}
	}()

	registry.Add("duplicate", scanner)
}
