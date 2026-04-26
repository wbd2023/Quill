package checks

import (
	"fmt"
	"go/ast"
	"go/token"
)

func collectSentinelErrorLocationViolations(
	fileSet *token.FileSet,
	file *ast.File,
	classifier PathClassifier,
) (violations []Violation) {
	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range valueSpec.Names {
				if !isSentinelErrorName(name.Name) {
					continue
				}

				violations = append(violations, Violation{
					Position: fileSet.Position(name.Pos()),
					Rule:     DiagnosticDomainErrorsLocation,
					Message: fmt.Sprintf(
						"sentinel errors must be declared in %s",
						classifier.FirstPattern(PathClassDomainErrors),
					),
				})
			}
		}
	}

	return violations
}
