package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

/* ------------------------------------- Interface Matching ------------------------------------- */

func TestStylecheckMatchesMockPrefixNaming(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(tempDir, "internal", "client", "application", "port", "account")
	mocksDirectory := filepath.Join(tempDir, "internal", "testkit", "mocks")

	portsPath := filepath.Join(portsDirectory, "user_repository.go")
	portsSource := `package ports

type UserRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`

	writeSourceFile(t, portsPath, portsSource)

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

	writeSourceFile(t, mockPath, mockSource)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
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
	portsDirectory := filepath.Join(tempDir, "internal", "client", "application", "port", "account")
	implDirectory := filepath.Join(tempDir, "internal", "adapters", "storage")

	portsPath := filepath.Join(portsDirectory, "user_repository.go")
	portsSource := `package ports

type UserRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`
	writeSourceFile(t, portsPath, portsSource)

	implPath := filepath.Join(implDirectory, "user_repository.go")
	implSource := `package storage

import "project/internal/client/application/port/account"

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
	writeSourceFile(t, implPath, implSource)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
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
	portsDirectory := filepath.Join(tempDir, "internal", "client", "application", "port", "account")
	implDirectory := filepath.Join(tempDir, "internal", "adapters", "storage")

	portsPath := filepath.Join(portsDirectory, "user_repository.go")
	portsSource := `package ports

type UserRepository interface {
	Save(value string) (err error)
	Load(value string) (err error)
}
`
	writeSourceFile(t, portsPath, portsSource)

	implPath := filepath.Join(implDirectory, "user_repository.go")
	implSource := `package storage

import "project/internal/client/application/port/account"

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
	writeSourceFile(t, implPath, implSource)

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, output:\n%s", output)
	}
}

func TestStylecheckReportsAmbiguousMockNaming(t *testing.T) {
	tempDir := t.TempDir()
	portsDirectory := filepath.Join(
		tempDir,
		"internal",
		"client",
		"application",
		"port",
		"identity",
	)
	mocksDirectory := filepath.Join(tempDir, "internal", "testkit", "mocks")

	portsPath := filepath.Join(portsDirectory, "identity_repository.go")
	portsSource := `package ports

type IdentityRepository interface {
	Load(value string) (err error)
}
`

	writeSourceFile(t, portsPath, portsSource)

	mockPrefixPath := filepath.Join(mocksDirectory, "identity_repository_prefix_mock.go")
	mockPrefixSource := `package mocks

type MockIdentityRepository struct{}

func (m *MockIdentityRepository) Load(value string) (err error) {
	return nil
}
`

	writeSourceFile(t, mockPrefixPath, mockPrefixSource)

	mockSuffixPath := filepath.Join(mocksDirectory, "identity_repository_suffix_mock.go")
	mockSuffixSource := `package mocks

type IdentityRepositoryMock struct{}

func (m *IdentityRepositoryMock) Load(value string) (err error) {
	return nil
}
`

	writeSourceFile(t, mockSuffixPath, mockSuffixSource)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail on ambiguous mock naming, output:\n%s", output)
	}

	if !strings.Contains(output, `multiple mock types match interface "IdentityRepository"`) {
		t.Fatalf("expected ambiguous mock naming violation, got:\n%s", output)
	}
}
