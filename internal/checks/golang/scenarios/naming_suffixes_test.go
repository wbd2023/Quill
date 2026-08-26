package scenarios

import (
	"path/filepath"
	"testing"
)

/* -------------------------------- Weightless Suffix Diagnostics ------------------------------- */

func TestGoStyleReportsWeightlessSuffixes(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package sample

type EngineImpl struct {
	CacheWrapper string
}

type Runner interface {
	RunInstance()
}

func NewEngineInstance() {}

func (engine EngineImpl) LoadWrapper() {}

var serviceInstance = 1

const limitWrapper = 2
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/weightless-suffix] weightless suffix "Impl" in type "EngineImpl"`,
	)
	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/weightless-suffix] weightless suffix "Wrapper" in field "CacheWrapper"`,
	)
	expectDiagnostic(
		t,
		result,
		diagnosticMatch{
			code:    "go/naming/weightless-suffix",
			message: `suffix "Instance" in interface method "RunInstance"`,
		},
	)
	expectDiagnostic(
		t,
		result,
		diagnosticMatch{
			code:    "go/naming/weightless-suffix",
			message: `suffix "Instance" in function "NewEngineInstance"`,
		},
	)
	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/weightless-suffix] weightless suffix "Wrapper" in method "LoadWrapper"`,
	)
	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/weightless-suffix] weightless suffix "Instance" in variable "serviceInstance"`,
	)
	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/weightless-suffix] weightless suffix "Wrapper" in constant "limitWrapper"`,
	)
}

/* --------------------------------- Package Stutter Diagnostics -------------------------------- */

func TestGoStyleReportsPackageStutter(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package catalog

type Catalog struct{}

func Catalog() *Catalog { return nil }

func (item Catalog) Reset() {}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stutterDiagnostics := 0
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Code == "go/naming/package-stutter" {
			stutterDiagnostics++
		}
	}

	if stutterDiagnostics != 2 {
		t.Fatalf(
			"expected 2 package-stutter diagnostics, got %d: %#v",
			stutterDiagnostics,
			result.Diagnostics,
		)
	}

	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/package-stutter] exported type "Catalog" stutters with package "catalog"`,
	)
	expectDiagnosticMessage(
		t,
		result,
		`[go/naming/package-stutter] exported function "Catalog" stutters with package "catalog"`,
	)
}

/* -------------------------------- Marker Suppression Scenarios -------------------------------- */

func TestGoStyleAllowsMarkedNamingDeclarations(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")

	sourceCode := `package sample

type CacheImpl struct { // style: allow-weightless-suffix because: mirrors generated code
	Name string
}

type UserIDWrapper struct{}

func BuildCache() {}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rejectDiagnosticMessage(t, result, `weightless suffix "Impl" in type "CacheImpl"`)
	rejectDiagnosticMessage(
		t,
		result,
		`weightless suffix "Wrapper" in type "UserIDWrapper"`,
	)
}

func TestGoStyleAllowsMarkedPackageStutter(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package tree

type Tree struct{} // style: allow-package-stutter because: foundational container type

func NewTree() *Tree { return &Tree{} }
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rejectDiagnosticMessage(
		t,
		result,
		`exported type "Tree" stutters with package "tree"`,
	)
}

/* ------------------------------------ Clean Naming Baseline ----------------------------------- */

func TestGoStylePassesCleanNaming(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "sample.go")
	sourceCode := `package inventory

type Item struct {
	Label string
}

func List(items []Item) (labels []string, err error) {
	labels = make([]string, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.Label)
	}

	return labels, nil
}
`
	writeSourceFile(t, sourcePath, sourceCode)

	result, err := runGoStyleResult(t, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("expected clean naming fixture to pass, diagnostics: %#v", result.Diagnostics)
	}
}
