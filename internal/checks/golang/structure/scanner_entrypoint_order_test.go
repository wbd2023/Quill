package structure

import (
	"testing"

	"ciphera/tools/internal/checks/golang/analysis"
)

func TestCheckScannerEntrypointOrderRejectsHelperBeforeCheck(t *testing.T) {
	fileSet, file := parseGoSource(t, `package text

func helper() {}

func CheckThing() {}
`)

	violations := CheckScannerEntrypointOrder(
		fileSet,
		file,
		"/repo/tools/internal/checks/text/example.go",
	)
	if len(violations) != 1 || violations[0].Rule != analysis.DiagnosticScannerEntrypointOrder {
		t.Fatalf("expected scanner entrypoint violation, got %#v", violations)
	}
}

func TestCheckScannerEntrypointOrderAcceptsCheckBeforeHelpers(t *testing.T) {
	fileSet, file := parseGoSource(t, `package text

func CheckThing() {}

func helper() {}
`)

	violations := CheckScannerEntrypointOrder(
		fileSet,
		file,
		"/repo/tools/internal/checks/text/example.go",
	)
	if len(violations) != 0 {
		t.Fatalf("expected scanner order to pass, got %#v", violations)
	}
}
