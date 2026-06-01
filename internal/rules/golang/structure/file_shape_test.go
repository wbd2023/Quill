package structure

import (
	"strings"
	"testing"

	"ciphera/tools/internal/rules/golang/analysis"
)

func TestCheckShapeRejectsVagueMixedFileName(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

type value struct{}

func helper() {}
`)

	violations := CheckShape(fileSet, file, "helpers.go", false)
	if !hasViolation(violations, analysis.DiagnosticFileShapeVagueName) {
		t.Fatalf("expected vague filename violation, got %#v", violations)
	}
}

func TestCheckShapeAllowsPackageWideTypesFile(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

type first struct{}
type second struct{}
`)

	violations := CheckShape(fileSet, file, "types.go", false)
	if hasViolation(violations, analysis.DiagnosticFileShapeVagueName) {
		t.Fatalf("expected package-wide types file to pass, got %#v", violations)
	}
}

func TestCheckShapeRejectsTinyGlueFile(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

func helper() {}
`)

	violations := CheckShape(fileSet, file, "glue.go", false)
	if !hasViolation(violations, analysis.DiagnosticFileShapeTinyGlue) {
		t.Fatalf("expected tiny glue violation, got %#v", violations)
	}
}

func TestCheckShapeRejectsLongFunction(t *testing.T) {
	var body strings.Builder
	for range 82 {
		body.WriteString("\t_ = 1\n")
	}

	fileSet, file := parseGoSource(
		t,
		"package example\n\nfunc longFunction() {\n"+body.String()+"}\n",
	)
	violations := CheckShape(fileSet, file, "shape.go", false)
	if !hasViolation(violations, analysis.DiagnosticFileShapeLongFunction) {
		t.Fatalf("expected long function violation, got %#v", violations)
	}
}
