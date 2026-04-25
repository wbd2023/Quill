package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStylecheckReportsContextAndResourceViolations(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected context/resource failure, output:\n%s", output)
	}

	for _, fragment := range []string{
		"contexts must not be stored on struct fields",
		"http.Client literals must set Timeout explicitly",
		"network requests must use an http.Client with an explicit timeout",
		"ignored close errors require an inline comment explaining why they are safe",
	} {
		if strings.Contains(output, fragment) {
			continue
		}

		t.Fatalf("expected %q violation, got:\n%s", fragment, output)
	}
}

func TestStylecheckAcceptsContextAndResourcePatterns(t *testing.T) {
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

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected valid context/resource fixture to pass, output:\n%s", output)
	}
}
