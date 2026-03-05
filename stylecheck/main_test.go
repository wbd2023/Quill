package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestStylecheckReportsUnnamedReturns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Config struct{}

func Bad(value string) error {
	return nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `[2.2] function "Bad" has unnamed return values`) {
		t.Fatalf("expected unnamed return violation, got:\n%s", output)
	}
}

func TestStylecheckPassesValidFile(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Config struct{}

func Good(value string) (err error) {
	if value == "" {
		return nil
	}

	return nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}

func TestStylecheckReportsPlaceholderReturnNames(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	// Intentionally uses placeholder return naming to assert rule enforcement.
	sourceCode := `package sample

func Bad() (result0 string) {
	return "bad"
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `[2.2] function "Bad" uses placeholder return name "result0"`) {
		t.Fatalf("expected placeholder return-name violation, got:\n%s", output)
	}
}

func TestStylecheckReportsNakedReturns(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

func Bad(value string) (err error) {
	if value == "" {
		return nil
	}

	return
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `[2.2] function "Bad" uses a naked return`) {
		t.Fatalf("expected naked-return violation, got:\n%s", output)
	}
}

func TestStylecheckReportsDirectDomainIdentifierCast(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "project/internal/core/domain"

func Bad(raw string) (id domain.IdentityID, err error) {
	id = domain.IdentityID(raw)
	return id, nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(output, `direct cast to domain.IdentityID is disallowed`) {
		t.Fatalf("expected direct-cast violation, got:\n%s", output)
	}
}

func TestStylecheckPassesDomainIdentifierParserUsage(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

import "project/internal/core/domain"

func Good(raw string) (id domain.IdentityID, err error) {
	id, err = domain.ParseIdentityID(raw)
	if err != nil {
		return "", err
	}
	return id, nil
}
`
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}

func TestStylecheckMatchesMockPrefixNaming(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "core", "ports")
	mocksDirectory := filepath.Join(tempDir, "internal", "mocks")

	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}

	if err := os.MkdirAll(mocksDirectory, 0o700); err != nil {
		t.Fatalf("mkdir mocks: %v", err)
	}

	portsPath := filepath.Join(portsDirectory, "user_repository.go")
	portsSource := `package ports

type UserRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`

	if err := os.WriteFile(portsPath, []byte(portsSource), 0o600); err != nil {
		t.Fatalf("write ports: %v", err)
	}

	mockPath := filepath.Join(mocksDirectory, "user_repository_mock.go")
	mockSource := `package mocks

type MockUserRepository struct{}

func (m *MockUserRepository) Load(value string) (err error) {
	return nil
}

func (m *MockUserRepository) Save(value string) (err error) {
	return nil
}
`

	if err := os.WriteFile(mockPath, []byte(mockSource), 0o600); err != nil {
		t.Fatalf("write mock: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`mock "MockUserRepository" for interface "UserRepository" method order mismatch`,
	) {
		t.Fatalf("expected prefixed mock-order violation, got:\n%s", output)
	}
}

func TestStylecheckReportsImplementationOrderMismatch(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "core", "ports")
	implDirectory := filepath.Join(tempDir, "internal", "adapters", "storage")

	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}
	if err := os.MkdirAll(implDirectory, 0o700); err != nil {
		t.Fatalf("mkdir impl: %v", err)
	}

	portsPath := filepath.Join(portsDirectory, "user_repository.go")
	portsSource := `package ports

type UserRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`
	if err := os.WriteFile(portsPath, []byte(portsSource), 0o600); err != nil {
		t.Fatalf("write ports: %v", err)
	}

	implPath := filepath.Join(implDirectory, "user_repository.go")
	implSource := `package storage

import "project/internal/core/ports"

type UserFileRepository struct{}

func (r *UserFileRepository) Load(value string) (err error) {
	_ = value
	return nil
}

func (r *UserFileRepository) Save(value string) (err error) {
	_ = value
	return nil
}

var _ ports.UserRepository = (*UserFileRepository)(nil)
`
	if err := os.WriteFile(implPath, []byte(implSource), 0o600); err != nil {
		t.Fatalf("write impl: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`implementation "UserFileRepository" for interface "UserRepository" method order mismatch`,
	) {
		t.Fatalf("expected implementation-order violation, got:\n%s", output)
	}
}

func TestStylecheckPassesImplementationOrderMatch(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "core", "ports")
	implDirectory := filepath.Join(tempDir, "internal", "adapters", "storage")

	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}
	if err := os.MkdirAll(implDirectory, 0o700); err != nil {
		t.Fatalf("mkdir impl: %v", err)
	}

	portsPath := filepath.Join(portsDirectory, "user_repository.go")
	portsSource := `package ports

type UserRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`
	if err := os.WriteFile(portsPath, []byte(portsSource), 0o600); err != nil {
		t.Fatalf("write ports: %v", err)
	}

	implPath := filepath.Join(implDirectory, "user_repository.go")
	implSource := `package storage

import "project/internal/core/ports"

type UserFileRepository struct{}

func (r *UserFileRepository) Save(value string) (err error) {
	_ = value
	return nil
}

func (r *UserFileRepository) Load(value string) (err error) {
	_ = value
	return nil
}

var _ ports.UserRepository = (*UserFileRepository)(nil)
`
	if err := os.WriteFile(implPath, []byte(implSource), 0o600); err != nil {
		t.Fatalf("write impl: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err != nil {
		t.Fatalf("expected stylecheck to pass, output:\n%s", output)
	}
}

func TestStylecheckReportsAmbiguousMockNaming(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "core", "ports")
	mocksDirectory := filepath.Join(tempDir, "internal", "mocks")

	if err := os.MkdirAll(portsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir ports: %v", err)
	}

	if err := os.MkdirAll(mocksDirectory, 0o700); err != nil {
		t.Fatalf("mkdir mocks: %v", err)
	}

	portsPath := filepath.Join(portsDirectory, "identity_repository.go")
	portsSource := `package ports

type IdentityRepository interface {
	Load(value string) (err error)
}
`

	if err := os.WriteFile(portsPath, []byte(portsSource), 0o600); err != nil {
		t.Fatalf("write ports: %v", err)
	}

	mockPrefixPath := filepath.Join(mocksDirectory, "identity_repository_prefix_mock.go")
	mockPrefixSource := `package mocks

type MockIdentityRepository struct{}

func (m *MockIdentityRepository) Load(value string) (err error) {
	return nil
}
`

	if err := os.WriteFile(mockPrefixPath, []byte(mockPrefixSource), 0o600); err != nil {
		t.Fatalf("write prefix mock: %v", err)
	}

	mockSuffixPath := filepath.Join(mocksDirectory, "identity_repository_suffix_mock.go")
	mockSuffixSource := `package mocks

type IdentityRepositoryMock struct{}

func (m *IdentityRepositoryMock) Load(value string) (err error) {
	return nil
}
`

	if err := os.WriteFile(mockSuffixPath, []byte(mockSuffixSource), 0o600); err != nil {
		t.Fatalf("write suffix mock: %v", err)
	}

	output, err := runStylecheck(tempDir)
	if err == nil {
		t.Fatalf("expected stylecheck to fail on ambiguous mock naming, output:\n%s", output)
	}

	if !strings.Contains(output, `multiple mock types match interface "IdentityRepository"`) {
		t.Fatalf("expected ambiguous mock naming violation, got:\n%s", output)
	}
}

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

/* ------------------------------------------- Helpers ------------------------------------------ */

func runStylecheck(targetDirectory string) (output string, err error) {
	command := exec.Command("go", "run", ".", targetDirectory)
	command.Dir = stylecheckModuleDirectory()

	rawOutput, runErr := command.CombinedOutput()
	return string(rawOutput), runErr
}

func stylecheckModuleDirectory() (directory string) {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}
