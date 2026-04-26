package checks

import (
	"strings"
	"testing"
)

func TestCheckFileShapeRejectsVagueMixedFileName(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

type value struct{}

func helper() {}
`)

	violations := CheckFileShape(fileSet, file, "helpers.go", false)
	if !hasViolation(violations, DiagnosticFileShapeVagueName) {
		t.Fatalf("expected vague filename violation, got %#v", violations)
	}
}

func TestCheckFileShapeAllowsPackageWideTypesFile(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

type first struct{}
type second struct{}
`)

	violations := CheckFileShape(fileSet, file, "types.go", false)
	if hasViolation(violations, DiagnosticFileShapeVagueName) {
		t.Fatalf("expected package-wide types file to pass, got %#v", violations)
	}
}

func TestCheckFileShapeRejectsTinyGlueFile(t *testing.T) {
	fileSet, file := parseGoSource(t, `package example

func helper() {}
`)

	violations := CheckFileShape(fileSet, file, "glue.go", false)
	if !hasViolation(violations, DiagnosticFileShapeTinyGlue) {
		t.Fatalf("expected tiny glue violation, got %#v", violations)
	}
}

func TestCheckFileShapeRejectsLongFunction(t *testing.T) {
	var body strings.Builder
	for range 82 {
		body.WriteString("\t_ = 1\n")
	}

	fileSet, file := parseGoSource(
		t,
		"package example\n\nfunc longFunction() {\n"+body.String()+"}\n",
	)
	violations := CheckFileShape(fileSet, file, "shape.go", false)
	if !hasViolation(violations, DiagnosticFileShapeLongFunction) {
		t.Fatalf("expected long function violation, got %#v", violations)
	}
}
