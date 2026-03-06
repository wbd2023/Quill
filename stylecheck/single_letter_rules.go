package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

var allowedSingleLetterNames = map[string]bool{
	"i": true,
	"j": true,
	"k": true,
	"_": true,
}

/* ---------------------------------------- Naming Rules ---------------------------------------- */

// checkSingleLetterVars flags single-letter variable names that are not
// loop indices (i, j, k) or method receivers (2.2).
func checkSingleLetterVars(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			violations = append(
				violations,
				checkSingleLetterFuncParams(fileSet, declaration, allowedSingleLetterNames)...,
			)

		case *ast.AssignStmt:
			violations = append(
				violations,
				checkSingleLetterAssignStmt(fileSet, declaration, allowedSingleLetterNames)...,
			)

		case *ast.RangeStmt:
			violations = append(
				violations,
				checkSingleLetterRangeStmt(fileSet, declaration, allowedSingleLetterNames)...,
			)

		case *ast.ValueSpec:
			violations = append(
				violations,
				checkSingleLetterValueSpec(fileSet, declaration, allowedSingleLetterNames)...,
			)
		}
		return true
	})

	return violations
}

func checkSingleLetterFuncParams(
	fileSet *token.FileSet,
	declaration *ast.FuncDecl,
	allowed map[string]bool,
) (violations []violation) {
	if declaration.Type.Params == nil {
		return nil
	}

	for _, field := range declaration.Type.Params.List {
		for _, name := range field.Names {
			violationValue := singleLetterNameViolation(
				fileSet,
				name,
				allowed,
				"2.2",
				fmt.Sprintf(
					"single-letter parameter %q in function %q",
					name.Name,
					declaration.Name.Name,
				),
			)
			if violationValue != nil {
				violations = append(violations, *violationValue)
			}
		}
	}

	return violations
}

func checkSingleLetterAssignStmt(
	fileSet *token.FileSet,
	declaration *ast.AssignStmt,
	allowed map[string]bool,
) (violations []violation) {
	if declaration.Tok != token.DEFINE {
		return nil
	}

	for _, lhs := range declaration.Lhs {
		identifierNode, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}

		violationValue := singleLetterNameViolation(
			fileSet,
			identifierNode,
			allowed,
			"2.2",
			fmt.Sprintf("single-letter variable %q", identifierNode.Name),
		)
		if violationValue != nil {
			violations = append(violations, *violationValue)
		}
	}

	return violations
}

func checkSingleLetterRangeStmt(
	fileSet *token.FileSet,
	declaration *ast.RangeStmt,
	allowed map[string]bool,
) (violations []violation) {
	if declaration.Tok != token.DEFINE {
		return nil
	}

	if key, ok := declaration.Key.(*ast.Ident); ok {
		violationValue := singleLetterNameViolation(
			fileSet,
			key,
			allowed,
			"2.2",
			fmt.Sprintf("single-letter range variable %q", key.Name),
		)
		if violationValue != nil {
			violations = append(violations, *violationValue)
		}
	}

	if declaration.Value != nil {
		if value, ok := declaration.Value.(*ast.Ident); ok {
			violationValue := singleLetterNameViolation(
				fileSet,
				value,
				allowed,
				"2.2",
				fmt.Sprintf("single-letter range variable %q", value.Name),
			)
			if violationValue != nil {
				violations = append(violations, *violationValue)
			}
		}
	}

	return violations
}

func checkSingleLetterValueSpec(
	fileSet *token.FileSet,
	declaration *ast.ValueSpec,
	allowed map[string]bool,
) (violations []violation) {
	for _, name := range declaration.Names {
		violationValue := singleLetterNameViolation(
			fileSet,
			name,
			allowed,
			"2.2",
			fmt.Sprintf("single-letter variable %q", name.Name),
		)
		if violationValue != nil {
			violations = append(violations, *violationValue)
		}
	}

	return violations
}

func singleLetterNameViolation(
	fileSet *token.FileSet,
	name *ast.Ident,
	allowed map[string]bool,
	rule string,
	message string,
) (violationValue *violation) {
	if len(name.Name) != 1 || allowed[name.Name] {
		return nil
	}

	return &violation{
		position: fileSet.Position(name.Pos()),
		rule:     rule,
		message:  message,
	}
}
