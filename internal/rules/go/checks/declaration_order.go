package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	declUnknown    = 0
	declConstants  = 1
	declErrors     = 2
	declTypes      = 3
	declAssertions = 4
)

const minCategoryKinds = 2

/* --------------------------------------- Ordering Rules --------------------------------------- */

// CheckFileStructureOrder enforces objective top-level declaration ordering.
// This check intentionally avoids subjective formatting requirements.
func CheckFileStructureOrder(
	fileSet *token.FileSet,
	file *ast.File,
) (violations []Violation) {
	highestSeenCategory := declUnknown
	seenCategories := map[int]bool{}

	for _, declaration := range file.Decls {
		currentCategory := classifyTopLevelDecl(declaration)
		if currentCategory == declUnknown {
			continue
		}

		seenCategories[currentCategory] = true
		if currentCategory < highestSeenCategory {
			violations = append(violations, Violation{
				Position: fileSet.Position(declaration.Pos()),
				Rule:     DiagnosticFileOrder,
				Message: fmt.Sprintf(
					"declaration group %q appears after %q",
					declCategoryName(currentCategory),
					declCategoryName(highestSeenCategory),
				),
			})
			continue
		}

		highestSeenCategory = currentCategory
	}

	if len(seenCategories) < minCategoryKinds {
		return nil
	}
	return violations
}

/* --------------------------------- Declaration Classification --------------------------------- */

func classifyTopLevelDecl(declaration ast.Decl) (category int) {
	switch typed := declaration.(type) {
	case *ast.GenDecl:
		switch typed.Tok {
		case token.CONST:
			return declConstants

		case token.TYPE:
			return declTypes

		case token.VAR:
			if isCompileTimeAssertionDecl(typed) {
				return declAssertions
			}

			if isSentinelErrorDecl(typed) {
				return declErrors
			}
		}
	}

	return declUnknown
}

func isCompileTimeAssertionDecl(declaration *ast.GenDecl) (found bool) {
	if declaration.Tok != token.VAR || len(declaration.Specs) == 0 {
		return false
	}

	for _, spec := range declaration.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok || len(valueSpec.Names) != 1 || valueSpec.Names[0].Name != "_" {
			return false
		}

		if valueSpec.Type == nil || len(valueSpec.Values) != 1 {
			return false
		}
	}

	return true
}

func isSentinelErrorDecl(declaration *ast.GenDecl) (found bool) {
	if declaration.Tok != token.VAR || len(declaration.Specs) == 0 {
		return false
	}

	for _, spec := range declaration.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok || len(valueSpec.Names) == 0 {
			return false
		}
		for _, name := range valueSpec.Names {
			if !strings.HasPrefix(name.Name, "Err") {
				return false
			}
		}
	}

	return true
}

func declCategoryName(category int) (name string) {
	switch category {
	case declConstants:
		return "constants"

	case declErrors:
		return "errors"

	case declTypes:
		return "types"

	case declAssertions:
		return "assertions"

	default:
		return "unknown"
	}
}
