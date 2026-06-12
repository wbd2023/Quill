package syntax

import (
	"fmt"
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
	"ciphera/tools/internal/checks/gopolicy"
)

// CheckDirectDomainValueCasts enforces parser/constructor usage for key domain values.
func CheckDirectDomainValueCasts(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	classifier analysis.PathClassifier,
	constructors gopolicy.DomainValueConstructors,
) (violations []analysis.Violation) {
	if classifier.HasRole(path, analysis.PathRoleDomain) {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok || len(callExpression.Args) != 1 {
			return true
		}

		selector, ok := callExpression.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		packageIdentifier, ok := selector.X.(*ast.Ident)
		if !ok || packageIdentifier.Name != "domain" {
			return true
		}

		recommendedConstructor, found := recommendedDomainValueConstructor(
			constructors,
			selector.Sel.Name,
		)
		if !found {
			return true
		}

		violations = append(violations, analysis.Violation{
			Position: fileSet.Position(callExpression.Pos()),
			Rule:     analysis.DiagnosticNoDirectDomainCasts,
			Message: fmt.Sprintf(
				"direct cast to domain.%s is disallowed; use %s",
				selector.Sel.Name,
				recommendedConstructor,
			),
		})
		return true
	})

	return violations
}
