package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsDataUsageViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "bad.go")
	sourceCode := `package service

type repository interface {
	Load() (err error)
}

type payload struct {
	Name string
	Age  int
}

func Build(repositoryPointer *repository, values []string) (payloadPointer *payload) {
	if values == nil || len(values) == 0 {
		return &payload{}
	}

	payloadPointer = &payload{"alice", 42}
	_ = repositoryPointer
	return payloadPointer
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected data-usage failure, diagnostics: %#v", result.Diagnostics)
	}

	for _, fragment := range []string{
		"functions and structs must pass interface values directly",
		"struct literals must use named fields by default",
		"slice emptiness checks must use len(values) instead of nil guards",
	} {
		if hasDiagnostic(result, diagnosticMatch{message: fragment}) {
			continue
		}

		t.Fatalf("expected %q violation, got: %#v", fragment, result.Diagnostics)
	}
}

func TestGoStyleAcceptsDataUsagePatterns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "good.go")
	sourceCode := `package service

type repository interface {
	Load() (err error)
}

type payload struct {
	Name string
	Age  int
}

func Build(store repository, values []string) (item payload) {
	if len(values) == 0 {
		return payload{}
	}

	_ = store
	item = payload{Name: "alice", Age: 42}
	return item
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected valid data-usage fixture to pass, diagnostics: %#v", result.Diagnostics)
	}
}
