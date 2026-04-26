package checks

import (
	"go/ast"
	"go/token"
)

func CheckDataUsage(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier PathClassifier,
) (violations []Violation) {
	if !classifier.HasClass(path, PathClassGoSource) {
		return nil
	}

	interfaceNames := collectInterfaceTypeNames(file)
	sliceNames := collectSliceNames(file)

	if !isTestFile {
		violations = append(
			violations,
			checkNamedStructLiterals(fileSet, file)...,
		)
	}

	violations = append(
		violations,
		checkPointerToInterfaces(fileSet, file, interfaceNames)...,
	)
	violations = append(
		violations,
		checkSliceEmptinessStyle(fileSet, file, sliceNames)...,
	)

	return violations
}

/* --------------------------------------- Type Discovery --------------------------------------- */

func collectInterfaceTypeNames(file *ast.File) (names map[string]bool) {
	names = make(map[string]bool)

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := typeSpec.Type.(*ast.InterfaceType); !ok {
				continue
			}

			names[typeSpec.Name.Name] = true
		}
	}

	return names
}

func collectSliceNames(file *ast.File) (names map[string]bool) {
	names = make(map[string]bool)

	for _, declaration := range file.Decls {
		switch typedDecl := declaration.(type) {
		case *ast.FuncDecl:
			collectSliceNamesFromFieldList(names, typedDecl.Type.Params)
			collectSliceNamesFromFieldList(names, typedDecl.Type.Results)

		case *ast.GenDecl:
			for _, spec := range typedDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok || !isSliceType(valueSpec.Type) {
					continue
				}

				for _, name := range valueSpec.Names {
					names[name.Name] = true
				}
			}
		}
	}

	return names
}

func collectSliceNamesFromFieldList(names map[string]bool, fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, field := range fields.List {
		if !isSliceType(field.Type) {
			continue
		}

		for _, name := range field.Names {
			names[name.Name] = true
		}
	}
}

func isSliceType(expression ast.Expr) (found bool) {
	_, ok := expression.(*ast.ArrayType)
	return ok
}

/* --------------------------------------- Struct Literals -------------------------------------- */

func checkNamedStructLiterals(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		literal, ok := node.(*ast.CompositeLit)
		if !ok || len(literal.Elts) == 0 {
			return true
		}

		if !isNamedStructLiteralType(literal.Type) {
			return true
		}

		for _, element := range literal.Elts {
			if _, ok := element.(*ast.KeyValueExpr); ok {
				return true
			}
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(literal.Pos()),
			Rule:     DiagnosticNamedStructLiterals,
			Message:  "struct literals must use named fields by default",
		})

		return true
	})

	return violations
}

func isNamedStructLiteralType(expression ast.Expr) (found bool) {
	switch expression.(type) {
	case *ast.Ident, *ast.SelectorExpr:
		return true
	default:
		return false
	}
}

/* --------------------------------------- Interface Usage -------------------------------------- */

func checkPointerToInterfaces(
	fileSet *token.FileSet,
	file *ast.File,
	interfaceNames map[string]bool,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch typed := node.(type) {
		case *ast.Field:
			ruleID, message, found := pointerToInterfaceViolation(
				typed.Type,
				interfaceNames,
				true,
			)
			if !found {
				return true
			}

			violations = append(violations, Violation{
				Position: fileSet.Position(typed.Type.Pos()),
				Rule:     ruleID,
				Message:  message,
			})

		case *ast.ValueSpec:
			ruleID, message, found := pointerToInterfaceViolation(
				typed.Type,
				interfaceNames,
				false,
			)
			if !found {
				return true
			}

			violations = append(violations, Violation{
				Position: fileSet.Position(typed.Type.Pos()),
				Rule:     ruleID,
				Message:  message,
			})
		}

		return true
	})

	return violations
}

/* ----------------------------------------- Slice Usage ---------------------------------------- */

func checkSliceEmptinessStyle(
	fileSet *token.FileSet,
	file *ast.File,
	sliceNames map[string]bool,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		binaryExpression, ok := node.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		switch binaryExpression.Op {
		case token.LOR, token.LAND:
		default:
			return true
		}

		name, found := sliceNilGuardPair(binaryExpression.X, binaryExpression.Y, sliceNames)
		if !found {
			name, found = sliceNilGuardPair(binaryExpression.Y, binaryExpression.X, sliceNames)
		}
		if !found {
			return true
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(binaryExpression.Pos()),
			Rule:     DiagnosticLenForSliceEmptiness,
			Message:  "slice emptiness checks must use len(" + name + ") instead of nil guards",
		})

		return true
	})

	return violations
}

func sliceNilGuardPair(
	nilSide ast.Expr,
	lenSide ast.Expr,
	sliceNames map[string]bool,
) (name string, found bool) {
	name, nilCompared, found := nilComparedSliceName(nilSide, sliceNames)
	if !found || !nilCompared {
		return "", false
	}

	if !containsSliceLenCheck(lenSide, name) {
		return "", false
	}

	return name, true
}

func nilComparedSliceName(
	expression ast.Expr,
	sliceNames map[string]bool,
) (name string, nilCompared bool, found bool) {
	binaryExpression, ok := expression.(*ast.BinaryExpr)
	if !ok {
		return "", false, false
	}

	switch binaryExpression.Op {
	case token.EQL, token.NEQ:
	default:
		return "", false, false
	}

	if identifier, ok := binaryExpression.X.(*ast.Ident); ok &&
		isNilIdentifier(binaryExpression.Y) {
		if sliceNames[identifier.Name] {
			return identifier.Name, true, true
		}
	}
	if identifier, ok := binaryExpression.Y.(*ast.Ident); ok &&
		isNilIdentifier(binaryExpression.X) {
		if sliceNames[identifier.Name] {
			return identifier.Name, true, true
		}
	}

	return "", false, false
}

func containsSliceLenCheck(expression ast.Expr, name string) (found bool) {
	binaryExpression, ok := expression.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	callExpression, ok := binaryExpression.X.(*ast.CallExpr)
	if !ok || !isLenCallForName(callExpression, name) {
		callExpression, ok = binaryExpression.Y.(*ast.CallExpr)
		if !ok || !isLenCallForName(callExpression, name) {
			return false
		}
	}

	return true
}

func isLenCallForName(callExpression *ast.CallExpr, name string) (found bool) {
	functionName, ok := callExpression.Fun.(*ast.Ident)
	if !ok || functionName.Name != "len" || len(callExpression.Args) != 1 {
		return false
	}

	identifier, ok := callExpression.Args[0].(*ast.Ident)
	return ok && identifier.Name == name
}

func isNilIdentifier(expression ast.Expr) (found bool) {
	identifier, ok := expression.(*ast.Ident)
	return ok && identifier.Name == "nil"
}

/* ---------------------------------------- Pointer Rules --------------------------------------- */

func pointerToInterfaceViolation(
	expression ast.Expr,
	interfaceNames map[string]bool,
	inField bool,
) (ruleID string, message string, found bool) {
	starExpression, ok := expression.(*ast.StarExpr)
	if !ok {
		return "", "", false
	}

	if !isInterfaceTypeExpression(starExpression.X, interfaceNames) {
		return "", "", false
	}

	if inField {
		return DiagnosticInterfaceValuesDirect,
			"functions and structs must pass interface values directly, not *interface types",
			true
	}

	return DiagnosticNoPointersToInterfaces,
		"pointers to interfaces are forbidden",
		true
}

func isInterfaceTypeExpression(expression ast.Expr, interfaceNames map[string]bool) (found bool) {
	switch typedExpression := expression.(type) {
	case *ast.InterfaceType:
		return true
	case *ast.Ident:
		return interfaceNames[typedExpression.Name]
	default:
		return false
	}
}
