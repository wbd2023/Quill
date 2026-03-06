package checker

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

/* ------------------------------------------ Constants ----------------------------------------- */

var placeholderReturnNamePattern = regexp.MustCompile(`^result[0-9]+$`)

/* ---------------------------------------- Return Rules ---------------------------------------- */

// checkNamedReturns ensures all functions, methods, and interface methods
// use named return values (2.2).
func checkNamedReturns(fileSet *token.FileSet, file *ast.File) (violations []violation) {
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
) (violations []violation) {
	if funcType.Results == nil || len(funcType.Results.List) == 0 {
		return nil
	}

	for _, field := range funcType.Results.List {
		if len(field.Names) == 0 {
			violations = append(violations, violation{
				position: fileSet.Position(funcType.Results.Pos()),
				rule:     "2.2",
				message:  fmt.Sprintf("function %q has unnamed return values", funcName),
			})
			return violations
		}

		for _, name := range field.Names {
			if placeholderReturnNamePattern.MatchString(name.Name) {
				violations = append(violations, violation{
					position: fileSet.Position(name.Pos()),
					rule:     "2.2",
					message: fmt.Sprintf(
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

// checkTypeElision ensures each parameter has its own type declaration (2.2).
func checkTypeElision(fileSet *token.FileSet, file *ast.File) (violations []violation) {
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
				violations = append(violations, violation{
					position: fileSet.Position(field.Pos()),
					rule:     "2.2",
					message: fmt.Sprintf(
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

// checkNakedReturns reports naked returns in functions that declare named return values (2.2).
func checkNakedReturns(fileSet *token.FileSet, file *ast.File) (violations []violation) {
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
					violations = append(violations, violation{
						position: fileSet.Position(typed.Pos()),
						rule:     "2.2",
						message: fmt.Sprintf(
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
