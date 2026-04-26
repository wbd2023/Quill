package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsContextAndResourceViolations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "bad.go")
	sourceCode := `package service

import (
	"context"
	"net/http"
)

type Service struct {
	ctx context.Context
}

func New() (service *Service) {
	client := &http.Client{}
	_ = client
	_, _ = http.Get("https://example.com")
	return &Service{}
}

func Close(response *http.Response) {
	defer func() {
		_ = response.Body.Close()
	}()
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected context/resource failure, diagnostics: %#v", result.Diagnostics)
	}

	for _, fragment := range []string{
		"contexts must not be stored on struct fields",
		"http.Client literals must set Timeout explicitly",
		"network requests must use an http.Client with an explicit timeout",
		"ignored close errors require an inline comment explaining why they are safe",
	} {
		if hasDiagnostic(result, diagnosticMatch{message: fragment}) {
			continue
		}

		t.Fatalf("expected %q violation, got: %#v", fragment, result.Diagnostics)
	}
}

func TestGoStyleAcceptsContextAndResourcePatterns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "internal", "client", "application", "service", "good.go")
	sourceCode := `package service

import (
	"net/http"
	"time"
)

const relayTimeout = 5 * time.Second

func New() (client *http.Client) {
	client = &http.Client{Timeout: relayTimeout}
	return client
}

func Close(response *http.Response) {
	defer func() {
		_ = response.Body.Close() // safe: response bodies are fully consumed before close
	}()
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf(
			"expected valid context/resource fixture to pass, diagnostics: %#v",
			result.Diagnostics,
		)
	}
}
