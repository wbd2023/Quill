package drivers

import (
	"context"
	"testing"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
)

func TestRepositoryScanDriverRejectsMissingScanner(t *testing.T) {
	driver := repositoryScanDriver(NewRepositoryScanners())
	_, err := driver(
		context.Background(),
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
		_ context.Context,
		_ execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return style.ExecutionResult{}, nil
	}

	registry := NewRepositoryScanners()
	registry.Add("pack", "duplicate", scanner)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate scanner ID to panic")
		}
	}()

	registry.Add("pack", "duplicate", scanner)
}
