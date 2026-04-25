package checks

import (
	"go/ast"
	"go/token"
)

func CheckProcessExecutionSafety(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []Violation) {
	execAliases := importAliases(file, "os/exec")
	if len(execAliases) == 0 {
		return nil
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
		if !ok || !execAliases[packageIdentifier.Name] {
			return true
		}

		commandIndex := shellCommandIndex(selector.Sel.Name)
		if commandIndex < 0 || len(callExpression.Args) <= commandIndex+1 {
			return true
		}

		commandName, found := stringLiteralValue(callExpression.Args[commandIndex])
		if !found || !isShellCommand(commandName) {
			return true
		}

		shellFlag, found := stringLiteralValue(callExpression.Args[commandIndex+1])
		if !found || !isShellInterpolationFlag(shellFlag) {
			return true
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(callExpression.Pos()),
			Rule:     DiagnosticNoShellInterpolation,
			Message: "process execution must avoid shell interpolation; pass command " +
				"arguments directly",
		})

		return true
	})

	return violations
}

func shellCommandIndex(selectorName string) (index int) {
	switch selectorName {
	case "Command":
		return 0
	case "CommandContext":
		return 1
	default:
		return -1
	}
}

func isShellCommand(commandName string) (found bool) {
	switch commandName {
	case "sh", "bash", "zsh", "/bin/sh", "/bin/bash", "/bin/zsh":
		return true
	default:
		return false
	}
}

func isShellInterpolationFlag(flag string) (found bool) {
	switch flag {
	case "-c", "-lc", "-ic":
		return true
	default:
		return false
	}
}
