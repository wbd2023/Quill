package scan

import (
	"context"
	"testing"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
)

func TestRepositoryScanDriverRejectsMissingScanner(t *testing.T) {
	driver := CheckDriver(driverkit.NewRepositoryScanners())
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

	registry := driverkit.NewRepositoryScanners()
	registry.Add("duplicate", scanner)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate scanner ID to panic")
		}
	}()

	registry.Add("duplicate", scanner)
}
