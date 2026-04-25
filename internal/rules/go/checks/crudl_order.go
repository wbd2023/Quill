package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"ciphera/tools/internal/rulepack"
)

const (
	crudlUnknown = 0
	crudlCreate  = 1
	crudlRead    = 2
	crudlUpdate  = 3
	crudlDelete  = 4
	crudlList    = 5
)

/* --------------------------------------- Ordering Rules --------------------------------------- */

// CheckCRUDLOrder validates CRUD-L method ordering inside application port interfaces.
func CheckCRUDLOrder(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	classifier PathClassifier,
) (violations []Violation) {
	if !classifier.HasClass(path, rulepack.PathClassApplicationPort) {
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

			lastCategory := crudlUnknown
			lastMethod := ""

			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}
				name := method.Names[0].Name
				category := crudlCategory(name)
				if category == crudlUnknown {
					continue
				}
				if lastCategory != crudlUnknown && category < lastCategory {
					violations = append(violations, Violation{
						Position: fileSet.Position(method.Pos()),
						Rule:     DiagnosticCRUDLOrder,
						Message: fmt.Sprintf(
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

/* -------------------------------------- CRUDL Categories -------------------------------------- */

// crudlCategory classifies method names into CRUD-L categories.
func crudlCategory(name string) (category int) {
	switch {
	case strings.HasPrefix(name, "List"):
		return crudlList

	case strings.HasPrefix(name, "Delete"),
		strings.HasPrefix(name, "Remove"),
		strings.HasPrefix(name, "Consume"):
		return crudlDelete

	case strings.HasPrefix(name, "Update"),
		strings.HasPrefix(name, "Set"),
		strings.HasPrefix(name, "Ack"):
		return crudlUpdate

	case strings.HasPrefix(name, "Read"),
		strings.HasPrefix(name, "Load"),
		strings.HasPrefix(name, "Get"),
		strings.HasPrefix(name, "Fetch"),
		strings.HasPrefix(name, "IdentityExists"),
		strings.HasPrefix(name, "Metadata"),
		strings.HasPrefix(name, "Fingerprint"),
		strings.HasPrefix(name, "Current"):
		return crudlRead

	case strings.HasPrefix(name, "Create"),
		strings.HasPrefix(name, "Save"),
		strings.HasPrefix(name, "Generate"),
		strings.HasPrefix(name, "Register"),
		strings.HasPrefix(name, "Initiate"),
		strings.HasPrefix(name, "Send"):
		return crudlCreate

	default:
		return crudlUnknown
	}
}
