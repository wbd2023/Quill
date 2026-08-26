package syntax

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/wbd2023/quill/internal/checks/golang/analysis"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
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
		call, ok := node.(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := selector.X.(*ast.Ident)
		if !ok || ident.Name != "domain" {
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
			Position: fileSet.Position(call.Pos()),
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
