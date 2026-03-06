package collect

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"stylecheck/internal/checker/support"
)

func CollectImplementationMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	implementations map[string][]MethodDecl,
) {
	isPortPath := strings.Contains(path, support.CorePortsPathSegment)
	isMockPath := strings.Contains(path, support.MocksPathSegment)
	if isPortPath || isMockPath {
		return
	}

	for _, declaration := range file.Decls {
		funcDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDeclaration.Recv == nil || len(funcDeclaration.Recv.List) == 0 {
			continue
		}

		receiverName := support.ReceiverTypeName(funcDeclaration.Recv.List[0].Type)
		if receiverName == "" {
			continue
		}

		key := typeDeclKey(path, receiverName)
		implementations[key] = append(implementations[key], MethodDecl{
			Name:     funcDeclaration.Name.Name,
			Position: fileSet.Position(funcDeclaration.Name.Pos()),
		})
	}
}

func CollectImplementationBindings(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	bindings *[]ImplementationBinding,
) {
	if strings.Contains(path, support.MocksPathSegment) {
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

			interfaceName := support.TypeNameFromExpr(valueSpec.Type)
			if interfaceName == "" || len(valueSpec.Values) != 1 {
				continue
			}

			implementationName := support.ImplementationTypeFromAssertion(valueSpec.Values[0])
			if implementationName == "" {
				continue
			}

			*bindings = append(*bindings, ImplementationBinding{
				InterfaceName:      interfaceName,
				ImplementationName: implementationName,
				ImplementationKey:  typeDeclKey(path, implementationName),
				Position:           fileSet.Position(valueSpec.Pos()),
			})
		}
	}
}

func typeDeclKey(path string, typeName string) (key string) {
	return fmt.Sprintf("%s::%s", filepath.ToSlash(filepath.Dir(path)), typeName)
}
