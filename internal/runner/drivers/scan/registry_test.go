package scan

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/runner"
)

func TestBuiltinPackScannersHaveDrivers(t *testing.T) {
	registry, err := builtin.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, rule := range registry.Rules() {
		execution, found := rule.Check.RepositoryScanExecution()
		if !found {
			continue
		}

		if _, found := scanners[execution.Scanner]; !found {
			t.Fatalf(
				"rule %q uses scanner %q without a driver",
				rule.ID,
				execution.Scanner,
			)
		}
	}
}

func TestAddRepositoryScannersRejectsDuplicateScannerID(t *testing.T) {
	scanner := func(
		_ runner.Context,
		_ contract.RepositoryScanExecution,
	) (contract.ExecutionResult, error) {
		return contract.ExecutionResult{}, nil
	}

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate scanner ID to panic")
		}
	}()

	addRepositoryScanners(
		map[string]repositoryScanner{"duplicate": scanner},
		map[string]repositoryScanner{"duplicate": scanner},
	)
}
