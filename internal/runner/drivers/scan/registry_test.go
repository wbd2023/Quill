package scan

import (
	"testing"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
)

func TestRepositoryScanDriverRejectsMissingScanner(t *testing.T) {
	driver := repositoryScanDriver(binding.NewRepositoryScanners())
	_, err := driver(
		runner.Context{},
		style.ExecutionSpec{
			Kind: style.ExecutionRepositoryScan,
			Detail: style.RepositoryScanExecution{
				Scanner: "missing",
			},
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

	registry := binding.NewRepositoryScanners()
	registry.Add("duplicate", scanner)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate scanner ID to panic")
		}
	}()

	registry.Add("duplicate", scanner)
}
