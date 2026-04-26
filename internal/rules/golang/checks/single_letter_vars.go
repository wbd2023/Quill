package checks

import (
	"fmt"
	"go/ast"
	"go/token"
)

/* ---------------------------------------- Naming Rules ---------------------------------------- */

// CheckSingleLetterVars flags single-letter variable names that are not loop indices (i, j, k)
// or method receivers.
func CheckSingleLetterVars(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			violations = append(
				violations,
				checkSingleLetterFuncParams(fileSet, declaration)...,
			)

		case *ast.AssignStmt:
			violations = append(
				violations,
				checkSingleLetterAssignStmt(fileSet, declaration)...,
			)

		case *ast.RangeStmt:
			violations = append(
				violations,
				checkSingleLetterRangeStmt(fileSet, declaration)...,
			)

		case *ast.ValueSpec:
			violations = append(
				violations,
				checkSingleLetterValueSpec(fileSet, declaration)...,
			)
		}
		return true
	})

	return violations
}

func checkSingleLetterFuncParams(
	fileSet *token.FileSet,
	declaration *ast.FuncDecl,
) (violations []Violation) {
	if declaration.Type.Params == nil {
		return nil
	}

	for _, field := range declaration.Type.Params.List {
		for _, name := range field.Names {
			violation := singleLetterNameViolation(
				fileSet,
				name,
				DiagnosticSingleLetterNames,
				fmt.Sprintf(
					"single-letter parameter %q in function %q",
					name.Name,
					declaration.Name.Name,
				),
			)
			if violation != nil {
				violations = append(violations, *violation)
			}
		}
	}

	return violations
}

func checkSingleLetterAssignStmt(
	fileSet *token.FileSet,
	declaration *ast.AssignStmt,
) (violations []Violation) {
	if declaration.Tok != token.DEFINE {
		return nil
	}

	for _, lhs := range declaration.Lhs {
		identifierNode, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}

		violation := singleLetterNameViolation(
			fileSet,
			identifierNode,
			DiagnosticSingleLetterNames,
			fmt.Sprintf("single-letter variable %q", identifierNode.Name),
		)
		if violation != nil {
			violations = append(violations, *violation)
		}
	}

	return violations
}

func checkSingleLetterRangeStmt(
	fileSet *token.FileSet,
	declaration *ast.RangeStmt,
) (violations []Violation) {
	if declaration.Tok != token.DEFINE {
		return nil
	}

	if key, ok := declaration.Key.(*ast.Ident); ok {
		violation := singleLetterNameViolation(
			fileSet,
			key,
			DiagnosticSingleLetterNames,
			fmt.Sprintf("single-letter range variable %q", key.Name),
		)
		if violation != nil {
			violations = append(violations, *violation)
		}
	}

	if declaration.Value != nil {
		if value, ok := declaration.Value.(*ast.Ident); ok {
			violation := singleLetterNameViolation(
				fileSet,
				value,
				DiagnosticSingleLetterNames,
				fmt.Sprintf("single-letter range variable %q", value.Name),
			)
			if violation != nil {
				violations = append(violations, *violation)
			}
		}
	}

	return violations
}

func checkSingleLetterValueSpec(
	fileSet *token.FileSet,
	declaration *ast.ValueSpec,
) (violations []Violation) {
	for _, name := range declaration.Names {
		violation := singleLetterNameViolation(
			fileSet,
			name,
			DiagnosticSingleLetterNames,
			fmt.Sprintf("single-letter variable %q", name.Name),
		)
		if violation != nil {
			violations = append(violations, *violation)
		}
	}

	return violations
}

func singleLetterNameViolation(
	fileSet *token.FileSet,
	name *ast.Ident,
	rule string,
	message string,
) (violation *Violation) {
	if len(name.Name) != 1 || isAllowedSingleLetterName(name.Name) {
		return nil
	}

	return &Violation{
		Position: fileSet.Position(name.Pos()),
		Rule:     rule,
		Message:  message,
	}
}

func isAllowedSingleLetterName(name string) (allowed bool) {
	switch name {
	case "i", "j", "k", "_":
		return true
	default:
		return false
	}
}
