package order

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"ciphera/tools/internal/rulepack"
	gochecks "ciphera/tools/internal/rules/go/checks"
)

func collectImplementationMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	pathClassifier gochecks.PathClassifier,
	implementations map[string][]methodDeclaration,
) {
	if pathClassifier.HasClass(path, rulepack.PathClassApplicationPort) ||
		pathClassifier.HasClass(path, rulepack.PathClassTestMocks) {
		return
	}

	for _, declaration := range file.Decls {
		funcDeclaration, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDeclaration.Recv == nil || len(funcDeclaration.Recv.List) == 0 {
			continue
		}

		receiverName := gochecks.ReceiverTypeName(funcDeclaration.Recv.List[0].Type)
		if receiverName == "" {
			continue
		}

		key := typeDeclarationKey(path, receiverName)
		implementations[key] = append(implementations[key], methodDeclaration{
			Name:     funcDeclaration.Name.Name,
			Position: fileSet.Position(funcDeclaration.Name.Pos()),
		})
	}
}

func collectImplementationBindings(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	pathClassifier gochecks.PathClassifier,
	bindings *[]implementationBinding,
) {
	if pathClassifier.HasClass(path, rulepack.PathClassTestMocks) {
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

			interfaceName := gochecks.TypeNameFromExpr(valueSpec.Type)
			if interfaceName == "" || len(valueSpec.Values) != 1 {
				continue
			}

			implementationName := gochecks.ImplementationTypeFromAssertion(valueSpec.Values[0])
			if implementationName == "" {
				continue
			}

			*bindings = append(*bindings, implementationBinding{
				InterfaceName:      interfaceName,
				ImplementationName: implementationName,
				ImplementationKey:  typeDeclarationKey(path, implementationName),
				Position:           fileSet.Position(valueSpec.Pos()),
			})
		}
	}
}

func typeDeclarationKey(path string, typeName string) (key string) {
	return fmt.Sprintf("%s::%s", filepath.ToSlash(filepath.Dir(path)), typeName)
}
