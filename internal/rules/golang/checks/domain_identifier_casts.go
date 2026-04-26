package checks

import (
	"fmt"
	"go/ast"
	"go/token"

	"ciphera/tools/internal/policy"
)

// CheckDirectDomainIdentifierCasts enforces parser/constructor usage for key domain IDs.
func CheckDirectDomainIdentifierCasts(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	classifier PathClassifier,
	identifiers policy.GoDomainIdentifierConfig,
) (violations []Violation) {
	if classifier.HasClass(path, PathClassDomain) {
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

		recommendedConstructor, found := recommendedDomainIdentifierConstructor(
			identifiers,
			selector.Sel.Name,
		)
		if !found {
			return true
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(callExpression.Pos()),
			Rule:     DiagnosticNoDirectDomainCasts,
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
