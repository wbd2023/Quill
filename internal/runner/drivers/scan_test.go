package drivers

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/pack/builtin"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

func TestRunRepositoryScanRuleAcceptsKnownScanner(t *testing.T) {
	context := testContext(t, fixtures.RepositoryRoot(t), contract.Scope("tools"))

	if _, err := repositoryScanExecutor(
		context,
		repositoryScanSpec(builtin.ScannerASCII),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanExecutor(ascii): %v", err)
	}
}

func TestRunRepositoryScanRuleRejectsUnknownScanner(t *testing.T) {
	context := testContext(t, fixtures.RepositoryRoot(t), contract.Scope("all"))

	if _, err := repositoryScanExecutor(
		context,
		repositoryScanSpec("unknown"),
		nil,
	); err == nil {
		t.Fatal("expected unknown scanner to be rejected")
	}
}

func TestRunRepositoryScanRuleSupportsAlternateProfile(t *testing.T) {
	fixtureRoot := t.TempDir()
	alternateProfile := alternatePolicyForTest(t)
	profiles.Write(t, fixtureRoot, alternateProfile)
	fixtures.WriteFile(t, fixtureRoot, "ALTROOT", "")
	fixtures.WriteFile(
		t,
		fixtureRoot,
		"go.mod",
		"module example.com/altchat\n\ngo 1.24.5\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "errors.go"),
		"package domain\n\nvar ErrMissing = error(nil)\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "app", "ports", "message_store.go"),
		"package ports\n\ntype Message"+"Store interface { ListMessages() }\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "app", "services", "message_service.go"),
		"package services\n\n"+
			"import (\n"+
			"\t\"example.com/altchat/internal/app/ports\"\n"+
			"\t\"example.com/altchat/internal/domain\"\n"+
			")\n\n"+
			"type Message"+"Repository interface {\n"+
			"\tListMessages() []domain.Message\n"+
			"}\n\n"+
			"type MessageService struct {\n"+
			"\tstore ports.Message"+"Store\n"+
			"}\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "message.go"),
		"package domain\n\ntype Message struct{}\n",
	)

	context := testContext(t, fixtureRoot, contract.Scope("all"))
	if result, err := repositoryScanExecutor(
		context,
		repositoryScanSpec(builtin.ScannerArchitecture),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanExecutor(architecture): %v\n%s", err, result.Output)
	}

	result, err := repositoryScanExecutor(
		context,
		repositoryScanSpec(builtin.ScannerVocabulary),
		nil,
	)
	if err == nil {
		t.Fatal("expected alternate vocabulary policy to reject Repository suffixes")
	}

	if !hasDiagnostic(
		result,
		"vocabulary/project-terms/go-type-suffix",
		"internal/app/services/message_service.go",
		8,
		"use Store not Repository",
	) {
		t.Fatalf("expected alternate vocabulary diagnostic, got: %#v", result.Diagnostics)
	}
}

func repositoryScanSpec(scanner string) (spec contract.ExecutionSpec) {
	return contract.ExecutionSpec{
		Kind: contract.ExecutorRepositoryScan,
		Detail: contract.RepositoryScanExecution{
			Scanner: scanner,
		},
	}
}
