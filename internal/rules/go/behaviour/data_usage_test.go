package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsDataUsageViolations(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected data-usage failure, output:\n%s", output)
	}

	for _, fragment := range []string{
		"functions and structs must pass interface values directly",
		"struct literals must use named fields by default",
		"slice emptiness checks must use len(values) instead of nil guards",
	} {
		if strings.Contains(output, fragment) {
			continue
		}

		t.Fatalf("expected %q violation, got:\n%s", fragment, output)
	}
}

func TestStylecheckAcceptsDataUsagePatterns(t *testing.T) {
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

func Build(repositoryValue repository, values []string) (payloadValue payload) {
	if len(values) == 0 {
		return payload{}
	}

	_ = repositoryValue
	payloadValue = payload{Name: "alice", Age: 42}
	return payloadValue
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected valid data-usage fixture to pass, output:\n%s", output)
	}
}
