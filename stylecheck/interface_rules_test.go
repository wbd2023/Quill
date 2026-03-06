package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

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
