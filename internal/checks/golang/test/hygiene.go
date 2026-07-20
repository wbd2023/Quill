package test

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

/* ------------------------------------- Test Hygiene Rules ------------------------------------- */

// CheckHygiene check hygiene.
func CheckHygiene(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []analysis.Violation) {
	if !strings.HasSuffix(path, "_test.go") {
		return nil
	}

	violations = append(violations, checkTestHelperPlacement(fileSet, file)...)

	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Body == nil || isTestEntrypoint(function.Name.Name) {
			continue
		}

		handleName, hasTestingHandle := testingHandleParameter(function)
		if !hasTestingHandle || callsTestingHandleHelper(function.Body, handleName) {
			continue
		}

		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(function.Pos()),
			Rule:     analysis.DiagnosticHelpersCallHelper,
			Message:  "test helpers that accept testing handles must call Helper()",
		})
	}

	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := callExpression.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		packageIdentifier, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}

		switch {
		case packageIdentifier.Name == "os" && selector.Sel.Name == "Setenv":
			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(callExpression.Pos()),
				Rule:     analysis.DiagnosticTestSetenv,
				Message:  "tests must use t.Setenv() instead of os.Setenv()",
			})

		case packageIdentifier.Name == "os" && selector.Sel.Name == "MkdirTemp":
			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(callExpression.Pos()),
				Rule:     analysis.DiagnosticTestTempDir,
				Message:  "tests must use t.TempDir() instead of os.MkdirTemp()",
			})

		case packageIdentifier.Name == "ioutil" && selector.Sel.Name == "TempDir":
			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(callExpression.Pos()),
				Rule:     analysis.DiagnosticTestTempDir,
				Message:  "tests must use t.TempDir() instead of ioutil.TempDir()",
			})

		case packageIdentifier.Name == "time" && selector.Sel.Name == "Sleep":
			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(callExpression.Pos()),
				Rule:     analysis.DiagnosticTestAvoidArbitrarySleeps,
				Message:  "tests must avoid time.Sleep() when a deterministic signal is possible",
			})
		}

		return true
	})

	return violations
}

func checkTestHelperPlacement(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []analysis.Violation) {
	pendingHelpers := make([]*ast.FuncDecl, 0)

	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			continue
		}

		if isTestEntrypoint(function.Name.Name) {
			for _, helper := range pendingHelpers {
				violations = append(violations, analysis.Violation{
					Position: fileSet.Position(helper.Pos()),
					Rule:     analysis.DiagnosticTestHelperOrder,
					Message:  "test helpers should appear below test cases",
				})
			}
			pendingHelpers = nil
			continue
		}

		if _, hasTestingHandle := testingHandleParameter(function); hasTestingHandle {
			pendingHelpers = append(pendingHelpers, function)
		}
	}

	return violations
}

/* -------------------------------------- Test Entrypoints -------------------------------------- */

func isTestEntrypoint(name string) (found bool) {
	return strings.HasPrefix(name, "Test") ||
		strings.HasPrefix(name, "Benchmark") ||
		strings.HasPrefix(name, "Fuzz") ||
		strings.HasPrefix(name, "Example")
}

/* --------------------------------------- Testing Handles -------------------------------------- */

func testingHandleParameter(function *ast.FuncDecl) (name string, found bool) {
	if function.Type.Params == nil {
		return "", false
	}

	for _, field := range function.Type.Params.List {
		if len(field.Names) == 0 || !isTestingHandle(field.Type) {
			continue
		}

		return field.Names[0].Name, true
	}

	return "", false
}

func isTestingHandle(expression ast.Expr) (found bool) {
	switch typedExpression := expression.(type) {
	case *ast.StarExpr:
		selector, ok := typedExpression.X.(*ast.SelectorExpr)
		if !ok {
			return false
		}

		packageIdentifier, ok := selector.X.(*ast.Ident)
		if !ok || packageIdentifier.Name != "testing" {
			return false
		}

		switch selector.Sel.Name {
		case "T", "B", "F":
			return true
		default:
			return false
		}

	case *ast.SelectorExpr:
		packageIdentifier, ok := typedExpression.X.(*ast.Ident)
		if !ok || packageIdentifier.Name != "testing" {
			return false
		}

		return typedExpression.Sel.Name == "TB"

	default:
		return false
	}
}

/* ---------------------------------------- Helper Calls ---------------------------------------- */

func callsTestingHandleHelper(body *ast.BlockStmt, handleName string) (found bool) {
	ast.Inspect(body, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := callExpression.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "Helper" {
			return true
		}

		identifier, ok := selector.X.(*ast.Ident)
		if !ok || identifier.Name != handleName {
			return true
		}

		found = true
		return false
	})

	return found
}
