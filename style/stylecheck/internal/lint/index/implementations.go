package index

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"stylecheck/internal/lint/paths"
	"stylecheck/internal/lint/syntax"
)

func CollectImplementationMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	implementations map[string][]MethodDecl,
) {
	if paths.IsApplicationPortPath(path) || paths.IsTestMockPath(path) {
		return
	}

	for _, declaration := range file.Decls {
		funcDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDeclaration.Recv == nil || len(funcDeclaration.Recv.List) == 0 {
			continue
		}

		receiverName := syntax.ReceiverTypeName(funcDeclaration.Recv.List[0].Type)
		if receiverName == "" {
			continue
		}

		key := TypeDeclKey(path, receiverName)
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
	if paths.IsTestMockPath(path) {
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

			interfaceName := syntax.TypeNameFromExpr(valueSpec.Type)
			if interfaceName == "" || len(valueSpec.Values) != 1 {
				continue
			}

			implementationName := syntax.ImplementationTypeFromAssertion(valueSpec.Values[0])
			if implementationName == "" {
				continue
			}

			*bindings = append(*bindings, ImplementationBinding{
				InterfaceName:      interfaceName,
				ImplementationName: implementationName,
				ImplementationKey:  TypeDeclKey(path, implementationName),
				Position:           fileSet.Position(valueSpec.Pos()),
			})
		}
	}
}

func TypeDeclKey(path string, typeName string) (key string) {
	return fmt.Sprintf("%s::%s", filepath.ToSlash(filepath.Dir(path)), typeName)
}
