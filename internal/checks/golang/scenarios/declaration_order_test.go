package scenarios

import (
	"path/filepath"
	"testing"
)

func TestGoStyleReportsFileStructureOrderViolation(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type Service struct{}

const maxRetries = 3

func NewService() (service *Service) {
	return &Service{}
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatalf("expected custom Go check to fail, diagnostics: %#v", result.Diagnostics)
	}

	expectDiagnosticMessage(
		t,
		result,
		`[go/order/file] declaration group "constants" appears after "types"`,
	)
}

func TestGoStylePassesValidFileStructureOrder(t *testing.T) {
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
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
