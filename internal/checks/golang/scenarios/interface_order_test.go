package scenarios

import (
	"path/filepath"
	"testing"
)

/* ------------------------------------- Interface Matching ------------------------------------- */

func TestGoStyleMatchesMockPrefixNaming(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(
		t,
		result,
		`mock "MockUserRepository" for interface "UserRepository" method order mismatch`,
	)
}

func TestGoStyleReportsImplementationOrderMismatch(t *testing.T) {
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

import "profile/internal/client/application/port/account"

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

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(
		t,
		result,
		`implementation "UserFileRepository" for interface "UserRepository" method order mismatch`,
	)
}

func TestGoStylePassesImplementationOrderMatch(t *testing.T) {
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

import "profile/internal/client/application/port/account"

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

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, diagnostics: %#v", result.Diagnostics)
	}
}

func TestGoStyleReportsAmbiguousMockNaming(t *testing.T) {
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

	result, err := runGoStyleResult(t, tempDir)
	if err == nil {
		t.Fatalf(
			"expected custom Go check to fail on ambiguous mock naming, diagnostics: %#v",
			result.Diagnostics,
		)
	}

	expectDiagnosticMessage(t, result, `multiple mock types match interface "IdentityRepository"`)
}
