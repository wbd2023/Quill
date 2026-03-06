package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsFileStructureOrderViolation(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Service struct{}

const maxRetries = 3

func NewService() (service *Service) {
	return &Service{}
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `[2.9] declaration group "constants" appears after "types"`) {
		t.Fatalf("expected file-structure violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidFileStructureOrder(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "errors"

const maxRetries = 3

var ErrInvalidValue = errors.New("invalid value")

type ServicePort interface {
	Save(value string) (err error)
}

type Service struct{}

func NewService() (service *Service) {
	return &Service{}
}

func (s *Service) Save(value string) (err error) {
	_ = value
	return nil
}

func helper(value string) (err error) {
	_ = value
	return nil
}

var _ ServicePort = (*Service)(nil)
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}
