package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

func collectImplementationMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	implementations map[string][]methodDecl,
) {
	isPortPath := strings.Contains(path, corePortsPathSegment)
	isMockPath := strings.Contains(path, mocksPathSegment)
	if isPortPath || isMockPath {
		return
	}

	for _, declaration := range file.Decls {
		funcDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDeclaration.Recv == nil || len(funcDeclaration.Recv.List) == 0 {
			continue
		}

		receiverName := receiverTypeName(funcDeclaration.Recv.List[0].Type)
		if receiverName == "" {
			continue
		}

		key := typeDeclKey(path, receiverName)
		implementations[key] = append(implementations[key], methodDecl{
			name:     funcDeclaration.Name.Name,
			position: fileSet.Position(funcDeclaration.Name.Pos()),
		})
	}
}

func collectImplementationBindings(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	bindings *[]implementationBinding,
) {
	if strings.Contains(path, mocksPathSegment) {
		return
	}

	for _, declaration := range file.Decls {
		genDeclaration, ok := declaration.(*ast.GenDecl)
		if !ok || genDeclaration.Tok != token.VAR {
			continue
		}

		for _, spec := range genDeclaration.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			if len(valueSpec.Names) != 1 || valueSpec.Names[0].Name != "_" {
				continue
			}

			interfaceName := typeNameFromExpr(valueSpec.Type)
			if interfaceName == "" || len(valueSpec.Values) != 1 {
				continue
			}

			implementationName := implementationTypeFromAssertion(valueSpec.Values[0])
			if implementationName == "" {
				continue
			}

			*bindings = append(*bindings, implementationBinding{
				interfaceName:      interfaceName,
				implementationName: implementationName,
				implementationKey:  typeDeclKey(path, implementationName),
				position:           fileSet.Position(valueSpec.Pos()),
			})
		}
	}
}

func typeDeclKey(path string, typeName string) (key string) {
	return fmt.Sprintf("%s::%s", filepath.ToSlash(filepath.Dir(path)), typeName)
}
