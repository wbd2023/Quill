package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// checkServiceTypeNaming enforces exported type naming in internal/core/services (2.2).
func checkServiceTypeNaming(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !strings.Contains(path, "/internal/core/services/") {
		return nil
	}

	if strings.Contains(path, "/internal/core/services/accountref/") {
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
			if strings.HasSuffix(name, "Service") ||
				strings.HasSuffix(name, "UseCase") ||
				strings.HasSuffix(name, "Config") {
				continue
			}
			violations = append(violations, violation{
				position: fileSet.Position(typeSpec.Pos()),
				rule:     "2.2",
				message: fmt.Sprintf(
					"exported type %q should end with Service, UseCase, or Config",
					name,
				),
			})
		}
	}
	return violations
}
