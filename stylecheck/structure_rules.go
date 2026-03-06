package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	crudUnknown = 0
	crudCreate  = 1
	crudRead    = 2
	crudUpdate  = 3
	crudDelete  = 4
	crudList    = 5
)

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

// checkServiceTypeNaming enforces exported type naming in internal/core/services (2.2).
func checkServiceTypeNaming(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !strings.Contains(path, "/internal/core/services/") {
		return nil
	}

	if strings.Contains(path, "/internal/core/services/accountref/") {
		return nil
	}

	for _, declaration := range file.Decls {
		genDecl, ok := declaration.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || !typeSpec.Name.IsExported() {
				continue
			}
			name := typeSpec.Name.Name
			if strings.HasSuffix(name, "Service") ||
				strings.HasSuffix(name, "UseCase") ||
				strings.HasSuffix(name, "Config") {
				continue
			}
			violations = append(violations, violation{
				position: fileSet.Position(typeSpec.Pos()),
				rule:     "2.2",
				message: fmt.Sprintf(
					"exported type %q should end with Service, UseCase, or Config",
					name,
				),
			})
		}
	}
	return violations
}

// checkCRUDLOrder validates CRUD-L method ordering inside ports interfaces (2.5).
func checkCRUDLOrder(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !strings.Contains(path, "/internal/core/ports/") {
		return nil
	}

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
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok || interfaceType.Methods == nil {
				continue
			}

			lastCategory := crudUnknown
			lastMethod := ""

			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}
				name := method.Names[0].Name
				category := crudCategory(name)
				if category == crudUnknown {
					continue
				}
				if lastCategory != crudUnknown && category < lastCategory {
					violations = append(violations, violation{
						position: fileSet.Position(method.Pos()),
						rule:     "2.5",
						message: fmt.Sprintf(
							"method %q in interface %q is out of CRUD-L order (after %q)",
							name,
							typeSpec.Name.Name,
							lastMethod,
						),
					})
				}
				lastCategory = category
				lastMethod = name
			}
		}
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

// crudCategory classifies method names into CRUD-L categories.
func crudCategory(name string) (category int) {
	switch {
	case strings.HasPrefix(name, "List"):
		return crudList
	case strings.HasPrefix(name, "Delete"),
		strings.HasPrefix(name, "Remove"),
		strings.HasPrefix(name, "Consume"):
		return crudDelete
	case strings.HasPrefix(name, "Update"),
		strings.HasPrefix(name, "Set"),
		strings.HasPrefix(name, "Ack"):
		return crudUpdate
	case strings.HasPrefix(name, "Read"),
		strings.HasPrefix(name, "Load"),
		strings.HasPrefix(name, "Get"),
		strings.HasPrefix(name, "Fetch"),
		strings.HasPrefix(name, "IdentityExists"),
		strings.HasPrefix(name, "Metadata"),
		strings.HasPrefix(name, "Fingerprint"),
		strings.HasPrefix(name, "Current"):
		return crudRead
	case strings.HasPrefix(name, "Create"),
		strings.HasPrefix(name, "Save"),
		strings.HasPrefix(name, "Generate"),
		strings.HasPrefix(name, "Register"),
		strings.HasPrefix(name, "Initiate"),
		strings.HasPrefix(name, "Send"):
		return crudCreate
	default:
		return crudUnknown
	}
}
