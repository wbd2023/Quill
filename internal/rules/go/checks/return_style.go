package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

/* ---------------------------------------- Return Rules ---------------------------------------- */

// CheckNamedReturns ensures all functions, methods, and interface methods use named,
// descriptive return values.
func CheckNamedReturns(fileSet *token.FileSet, file *ast.File) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			violations = append(
				violations,
				checkFuncReturns(fileSet, declaration.Name.Name, declaration.Type)...,
			)

		case *ast.InterfaceType:
			if declaration.Methods == nil {
				return true
			}
			for _, method := range declaration.Methods.List {
				funcType, ok := method.Type.(*ast.FuncType)
				if !ok {
					continue // embedded interface, skip
				}
				methodName := "(anonymous)"
				if len(method.Names) > 0 {
					methodName = method.Names[0].Name
				}
				violations = append(violations, checkFuncReturns(fileSet, methodName, funcType)...)
			}
		}
		return true
	})
	return violations
}

// checkFuncReturns reports a violation if any return value is unnamed.
func checkFuncReturns(
	fileSet *token.FileSet,
	funcName string,
	funcType *ast.FuncType,
) (violations []Violation) {
	if funcType.Results == nil || len(funcType.Results.List) == 0 {
		return nil
	}

	for _, field := range funcType.Results.List {
		if len(field.Names) == 0 {
			violations = append(violations, Violation{
				Position: fileSet.Position(funcType.Results.Pos()),
				Rule:     DiagnosticNamedReturnValues,
				Message:  fmt.Sprintf("function %q has unnamed return values", funcName),
			})
			return violations
		}

		for _, name := range field.Names {
			if isPlaceholderReturnName(name.Name) {
				violations = append(violations, Violation{
					Position: fileSet.Position(name.Pos()),
					Rule:     DiagnosticNoPlaceholderReturnNames,
					Message: fmt.Sprintf(
						"function %q uses placeholder return name %q",
						funcName,
						name.Name,
					),
				})
			}
		}
	}

	return violations
}

// CheckTypeElision ensures each parameter has its own type declaration.
func CheckTypeElision(fileSet *token.FileSet, file *ast.File) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		funcType, ok := node.(*ast.FuncType)
		if !ok {
			return true
		}

		if funcType.Params == nil {
			return true
		}

		for _, field := range funcType.Params.List {
			if len(field.Names) > 1 {
				names := make([]string, len(field.Names))
				for index, name := range field.Names {
					names[index] = name.Name
				}
				violations = append(violations, Violation{
					Position: fileSet.Position(field.Pos()),
					Rule:     DiagnosticNoTypeElision,
					Message: fmt.Sprintf(
						"type elision: parameters %s share a type",
						strings.Join(names, ", "),
					),
				})
			}
		}

		return true
	})

	return violations
}

// CheckNakedReturns reports naked returns in functions that declare named return values.
func CheckNakedReturns(fileSet *token.FileSet, file *ast.File) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		functionDecl, ok := node.(*ast.FuncDecl)
		if !ok ||
			functionDecl.Type == nil ||
			functionDecl.Type.Results == nil ||
			functionDecl.Body == nil {
			return true
		}

		if !funcHasNamedReturns(functionDecl.Type) {
			return true
		}

		ast.Inspect(functionDecl.Body, func(bodyNode ast.Node) bool {
			switch typed := bodyNode.(type) {
			case *ast.FuncLit:
				return false

			case *ast.ReturnStmt:
				if len(typed.Results) == 0 {
					violations = append(violations, Violation{
						Position: fileSet.Position(typed.Pos()),
						Rule:     DiagnosticNoNakedReturns,
						Message: fmt.Sprintf(
							"function %q uses a naked return",
							functionDecl.Name.Name,
						),
					})
				}
			}

			return true
		})

		return true
	})

	return violations
}

func funcHasNamedReturns(functionType *ast.FuncType) (found bool) {
	for _, resultField := range functionType.Results.List {
		if len(resultField.Names) > 0 {
			return true
		}
	}

	return false
}

func isPlaceholderReturnName(name string) (placeholder bool) {
	if !strings.HasPrefix(name, "result") {
		return false
	}

	suffix := strings.TrimPrefix(name, "result")
	if suffix == "" {
		return false
	}

	for _, runeValue := range suffix {
		if runeValue < '0' || runeValue > '9' {
			return false
		}
	}

	return true
}
