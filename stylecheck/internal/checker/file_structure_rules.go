package checker

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

// checkFileStructureOrder enforces objective top-level declaration ordering (2.9).
// This check intentionally avoids subjective formatting requirements.
func checkFileStructureOrder(fileSet *token.FileSet, file *ast.File) (violations []violation) {
	highestSeenCategory := declUnknown
	seenCategories := map[int]bool{}

	for _, declaration := range file.Decls {
		currentCategory := classifyTopLevelDecl(declaration)
		if currentCategory == declUnknown {
			continue
		}

		seenCategories[currentCategory] = true
		if currentCategory < highestSeenCategory {
			violations = append(violations, violation{
				position: fileSet.Position(declaration.Pos()),
				rule:     "2.9",
				message: fmt.Sprintf(
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

/* ----------------------------------- Classification Helpers ----------------------------------- */

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
