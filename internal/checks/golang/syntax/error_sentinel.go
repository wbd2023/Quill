package syntax

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func collectSentinelErrorLocationViolations(
	fileSet *token.FileSet,
	file *ast.File,
	classifier analysis.PathClassifier,
) (violations []analysis.Violation) {
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

				violations = append(violations, analysis.Violation{
					Position: fileSet.Position(name.Pos()),
					Rule:     analysis.DiagnosticDomainErrorsLocation,
					Message: fmt.Sprintf(
						"sentinel errors must be declared in %s",
						classifier.FirstPattern(analysis.PathRoleDomainErrors),
					),
				})
			}
		}
	}

	return violations
}
