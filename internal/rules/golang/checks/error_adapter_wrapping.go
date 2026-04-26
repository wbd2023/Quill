package checks

import (
	"go/ast"
	"go/token"
)

// CheckAdapterErrorWrapping rejects bare error propagation in adapters.
func CheckAdapterErrorWrapping(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
) (violations []Violation) {
	if isTestFile || !classifier.HasClass(path, PathClassConcreteInfra) {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		function, ok := node.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			return true
		}

		ast.Inspect(function.Body, func(bodyNode ast.Node) bool {
			switch statement := bodyNode.(type) {
			case *ast.FuncLit:
				return false

			case *ast.ReturnStmt:
				if !isBareErrReturn(statement) {
					return true
				}

				violations = append(violations, Violation{
					Position: fileSet.Position(statement.Return),
					Rule:     DiagnosticAdapterWrapsCause,
					Message:  "adapter error returns must wrap low-level errors with context (%w)",
				})
			}

			return true
		})

		return false
	})

	return violations
}
