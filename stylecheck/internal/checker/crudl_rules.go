package checker

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

const (
	crudUnknown = 0
	crudCreate  = 1
	crudRead    = 2
	crudUpdate  = 3
	crudDelete  = 4
	crudList    = 5
)

/* --------------------------------------- Ordering Rules --------------------------------------- */

// checkCRUDLOrder validates CRUD-L method ordering inside ports interfaces (2.5).
func checkCRUDLOrder(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !strings.Contains(path, corePortsPathSegment) {
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

/* ------------------------------------------- Helpers ------------------------------------------ */

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
