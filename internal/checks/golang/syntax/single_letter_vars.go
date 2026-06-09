package syntax

import (
	"fmt"
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

/* ---------------------------------------- Naming Rules ---------------------------------------- */

// CheckSingleLetterVars flags single-letter variable names that are not loop indices (i, j, k)
// or method receivers.
func CheckSingleLetterVars(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []analysis.Violation) {
	loopIndexes := loopIndexPositions(file)

	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl:
			violations = append(
				violations,
				checkSingleLetterFuncParams(fileSet, declaration, loopIndexes)...,
			)

		case *ast.AssignStmt:
			violations = append(
				violations,
				checkSingleLetterAssignStmt(fileSet, declaration, loopIndexes)...,
			)

		case *ast.RangeStmt:
			violations = append(
				violations,
				checkSingleLetterRangeStmt(fileSet, declaration, loopIndexes)...,
			)

		case *ast.ValueSpec:
			violations = append(
				violations,
				checkSingleLetterValueSpec(fileSet, declaration, loopIndexes)...,
			)
		}
		return true
	})

	return violations
}

func checkSingleLetterFuncParams(
	fileSet *token.FileSet,
	declaration *ast.FuncDecl,
	loopIndexes map[token.Pos]bool,
) (violations []analysis.Violation) {
	if declaration.Type.Params == nil {
		return nil
	}

	for _, field := range declaration.Type.Params.List {
		for _, name := range field.Names {
			violation := singleLetterNameViolation(
				fileSet,
				name,
				loopIndexes,
				analysis.DiagnosticSingleLetterNames,
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
	loopIndexes map[token.Pos]bool,
) (violations []analysis.Violation) {
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
			loopIndexes,
			analysis.DiagnosticSingleLetterNames,
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
	loopIndexes map[token.Pos]bool,
) (violations []analysis.Violation) {
	if declaration.Tok != token.DEFINE {
		return nil
	}

	if key, ok := declaration.Key.(*ast.Ident); ok {
		violation := singleLetterNameViolation(
			fileSet,
			key,
			loopIndexes,
			analysis.DiagnosticSingleLetterNames,
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
				loopIndexes,
				analysis.DiagnosticSingleLetterNames,
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
	loopIndexes map[token.Pos]bool,
) (violations []analysis.Violation) {
	for _, name := range declaration.Names {
		violation := singleLetterNameViolation(
			fileSet,
			name,
			loopIndexes,
			analysis.DiagnosticSingleLetterNames,
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
	loopIndexes map[token.Pos]bool,
	rule string,
	message string,
) (violation *analysis.Violation) {
	if len(name.Name) != 1 || name.Name == "_" || loopIndexes[name.Pos()] {
		return nil
	}

	return &analysis.Violation{
		Position: fileSet.Position(name.Pos()),
		Rule:     rule,
		Message:  message,
	}
}

/* -------------------------------------- Loop Index Names -------------------------------------- */

func loopIndexPositions(file *ast.File) (positions map[token.Pos]bool) {
	positions = make(map[token.Pos]bool)
	ast.Inspect(file, func(node ast.Node) bool {
		switch statement := node.(type) {
		case *ast.ForStmt:
			collectLoopIndexStatement(statement.Init, positions)
		case *ast.RangeStmt:
			collectLoopIndexExpression(statement.Key, positions)
		}

		return true
	})

	return positions
}

func collectLoopIndexStatement(statement ast.Stmt, positions map[token.Pos]bool) {
	switch typed := statement.(type) {
	case *ast.AssignStmt:
		for _, expression := range typed.Lhs {
			collectLoopIndexExpression(expression, positions)
		}
	case *ast.DeclStmt:
		collectLoopIndexDeclaration(typed.Decl, positions)
	}
}

func collectLoopIndexDeclaration(declaration ast.Decl, positions map[token.Pos]bool) {
	genericDeclaration, ok := declaration.(*ast.GenDecl)
	if !ok {
		return
	}

	for _, spec := range genericDeclaration.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, name := range valueSpec.Names {
			collectLoopIndexName(name, positions)
		}
	}
}

func collectLoopIndexExpression(expression ast.Expr, positions map[token.Pos]bool) {
	name, ok := expression.(*ast.Ident)
	if !ok {
		return
	}

	collectLoopIndexName(name, positions)
}

func collectLoopIndexName(name *ast.Ident, positions map[token.Pos]bool) {
	if !isLoopIndexName(name.Name) {
		return
	}

	positions[name.Pos()] = true
}

func isLoopIndexName(name string) (allowed bool) {
	switch name {
	case "i", "j", "k":
		return true
	default:
		return false
	}
}
