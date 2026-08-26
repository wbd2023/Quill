package syntax

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/wbd2023/quill/internal/checks/golang/analysis"
)

// package-stutter constants.
const (
	packageStutterMarker = "allow-package-stutter"
)

/* ---------------------------------- Package Stutter Detection --------------------------------- */

// CheckPackageStutter flags exported top-level type and function declarations whose names repeat
// the enclosing package name case-insensitively (`package catalog` declaring `type Catalog`).
// Methods are exempt. Declarations may opt out with an inline `style: allow-package-stutter`
// marker on the declaration line.
func CheckPackageStutter(
	fileSet *token.FileSet,
	file *ast.File,
	lines []string,
) (violations []analysis.Violation) {
	for _, declaration := range file.Decls {
		violations = append(
			violations,
			checkPackageStutterDeclaration(fileSet, file.Name.Name, lines, declaration)...,
		)
	}

	return violations
}

func checkPackageStutterDeclaration(
	fileSet *token.FileSet,
	packageName string,
	lines []string,
	declaration ast.Decl,
) (violations []analysis.Violation) {
	switch declaration := declaration.(type) {
	case *ast.GenDecl:
		if declaration.Tok != token.TYPE {
			return nil
		}

		for _, spec := range declaration.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			violation := packageStutterViolation(fileSet, packageName, lines, "type", typeSpec.Name)
			if violation != nil {
				violations = append(violations, *violation)
			}
		}

	case *ast.FuncDecl:
		if declaration.Recv != nil {
			return nil
		}

		violation := packageStutterViolation(
			fileSet,
			packageName,
			lines,
			"function",
			declaration.Name,
		)
		if violation != nil {
			violations = append(violations, *violation)
		}
	}

	return violations
}

func packageStutterViolation(
	fileSet *token.FileSet,
	packageName string,
	lines []string,
	kind string,
	name *ast.Ident,
) (violation *analysis.Violation) {
	if !ast.IsExported(name.Name) || !strings.EqualFold(name.Name, packageName) ||
		allowsInlineException(lines, fileSet.Position(name.Pos()), packageStutterMarker) {
		return nil
	}

	return &analysis.Violation{
		Position: fileSet.Position(name.Pos()),
		Rule:     analysis.DiagnosticPackageStutter,
		Message: fmt.Sprintf(
			"exported %s %q stutters with package %q",
			kind,
			name.Name,
			packageName,
		),
	}
}
