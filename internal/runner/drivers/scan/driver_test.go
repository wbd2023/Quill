package scan

import (
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/text"
	"ciphera/tools/internal/pack/shipped/vocabulary"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

func TestRunRepositoryScanRuleAcceptsKnownScanner(t *testing.T) {
	context := testContext(t, testutil.RepositoryRoot(t), style.Scope("tools"))

	if _, err := testRepositoryScanDriver()(
		context,
		repositoryScanSpec(text.ScannerASCII),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanDriver(ascii): %v", err)
	}
}

func TestRunRepositoryScanRuleRejectsUnknownScanner(t *testing.T) {
	context := testContext(t, testutil.RepositoryRoot(t), style.Scope("all"))

	_, err := testRepositoryScanDriver()(
		context,
		repositoryScanSpec("unknown"),
		nil,
	)
	if err == nil {
		t.Fatal("expected unknown scanner to be rejected")
	}

	if !strings.Contains(err.Error(), `"unknown"`) {
		t.Fatalf("error = %q, want scanner ID", err)
	}
}

func TestRunRepositoryScanRuleSupportsAlternateProfile(t *testing.T) {
	fixtureRoot := t.TempDir()
	alternateProfile := buildScanDriverPolicyFixture(t)
	profiles.Write(t, fixtureRoot, alternateProfile)
	testutil.WriteFile(t, fixtureRoot, "ALTROOT", "")
	testutil.WriteFile(
		t,
		fixtureRoot,
		"go.mod",
		"module example.com/altchat\n\ngo 1.24.5\n",
	)
	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "errors.go"),
		"package domain\n\nvar ErrMissing = error(nil)\n",
	)
	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "app", "ports", "message_store.go"),
		"package ports\n\ntype Message"+"Store interface { ListMessages() }\n",
	)
	testutil.WriteFile(
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
	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "message.go"),
		"package domain\n\ntype Message struct{}\n",
	)

	context := testContext(t, fixtureRoot, style.Scope("all"))
	if result, err := testRepositoryScanDriver()(
		context,
		repositoryScanSpec(golang.ScannerArchitecture),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanDriver(architecture): %v\n%s", err, result.Output)
	}

	result, err := testRepositoryScanDriver()(
		context,
		repositoryScanSpec(vocabulary.ScannerVocabulary),
		nil,
	)
	if err == nil {
		t.Fatal("expected alternate vocabulary policy to reject Repository suffixes")
	}

	if !hasDiagnosticMatching(
		result,
		"vocabulary/project-terms/go-type-suffix",
		"internal/app/services/message_service.go",
		8,
		"use Store not Repository",
	) {
		t.Fatalf("expected alternate vocabulary diagnostic, got: %#v", result.Diagnostics)
	}
}

func repositoryScanSpec(scanner string) (spec style.ExecutionSpec) {
	return style.ExecutionSpec{
		Kind: style.ExecutionRepositoryScan,
		Detail: style.RepositoryScanExecution{
			Scanner: scanner,
		},
	}
}

func testRepositoryScanDriver() (driver runner.Driver) {
	scanners := runtimebinding.NewRepositoryScanners()
	scanners.Add(text.ScannerASCII, CheckASCII())
	scanners.Add(golang.ScannerArchitecture, CheckGoArchitecture(golang.PackID))
	scanners.Add(vocabulary.ScannerVocabulary, CheckVocabulary(vocabulary.PackID))
	return repositoryScanDriver(scanners)
}
