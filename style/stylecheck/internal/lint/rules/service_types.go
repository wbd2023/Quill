package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"stylecheck/internal/lint/paths"
	"stylecheck/internal/lint/report"
)

// CheckServiceTypeNaming enforces exported type naming in application service packages (2.2).
func CheckServiceTypeNaming(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []report.Violation) {
	if !paths.IsApplicationServicePath(path) {
		return nil
	}

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || !typeSpec.Name.IsExported() {
				continue
			}
			name := typeSpec.Name.Name
			if strings.HasSuffix(name, "Service") || strings.HasSuffix(name, "Config") {
				continue
			}
			violations = append(violations, report.Violation{
				Position: fileSet.Position(typeSpec.Pos()),
				Rule:     "2.2",
				Message: fmt.Sprintf(
					"exported type %q should end with Service or Config",
					name,
				),
			})
		}
	}
	return violations
}
