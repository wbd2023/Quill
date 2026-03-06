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

func TestStylecheckReportsServiceTypeNamingViolation(t *testing.T) {
	tempDir := t.TempDir()
	serviceDirectory := filepath.Join(tempDir, "internal", "core", "services", "message")
	if err := os.MkdirAll(serviceDirectory, 0o700); err != nil {
		t.Fatalf("mkdir services: %v", err)
	}

	sourcePath := filepath.Join(serviceDirectory, "service.go")
	sourceCode := `package message

type Manager struct{}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[2.2] exported type "Manager" should end with Service, UseCase, or Config`,
	) {
		t.Fatalf("expected service-type naming violation, got:\n%s", output)
	}
}

func TestStylecheckReportsCRUDLOrderViolation(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "core", "ports")
	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}

	sourcePath := filepath.Join(portsDirectory, "repository.go")
	sourceCode := `package ports

type MessageRepository interface {
	Load(value string) (err error)
	Save(value string) (err error)
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[2.5] method "Save" in interface "MessageRepository" is out of CRUD-L order`,
	) {
		t.Fatalf("expected CRUD-L order violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidServiceTypeNamingAndCRUDLOrder(t *testing.T) {
	tempDir := t.TempDir()
	serviceDirectory := filepath.Join(tempDir, "internal", "core", "services", "message")
	portsDirectory := filepath.Join(tempDir, "internal", "core", "ports")

	if err := os.MkdirAll(serviceDirectory, 0o700); err != nil {
		t.Fatalf("mkdir services: %v", err)
	}

	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}

	servicePath := filepath.Join(serviceDirectory, "service.go")
	serviceSource := `package message

type MessageService struct{}
`
	if err := os.WriteFile(servicePath, []byte(serviceSource), 0o600); err != nil {
		t.Fatalf("write service source: %v", err)
	}

	portsPath := filepath.Join(portsDirectory, "repository.go")
	portsSource := `package ports

type MessageRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`
	if err := os.WriteFile(portsPath, []byte(portsSource), 0o600); err != nil {
		t.Fatalf("write ports source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}
