package behaviour

import (
	"path/filepath"
	"strings"
	"testing"
)

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
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err == nil {
		t.Fatalf("expected custom Go check to fail, output:\n%s", output)
	}

	if !strings.Contains(
		output,
		`[go/order/file] declaration group "constants" appears after "types"`,
	) {
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
	writeSourceFile(t, sourcePath, sourceCode)

	output, err := runGoStyleCheck(t, tempDir)
	if err != nil {
		t.Fatalf("expected custom Go check to pass, output:\n%s", output)
	}
}
