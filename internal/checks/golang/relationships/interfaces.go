package relationships

import (
	"go/ast"
	"go/token"

	"github.com/wbd2023/Quill/internal/checks/golang/analysis"
)

func collectInterfaces(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	pathClassifier analysis.PathClassifier,
	interfaces map[string]interfaceDeclaration,
) {
	if !pathClassifier.HasRole(path, analysis.PathRoleApplicationPort) {
		return
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

			methods := make([]methodDeclaration, 0, len(interfaceType.Methods.List))
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}
				methods = append(methods, methodDeclaration{
					Name:     method.Names[0].Name,
					Position: fileSet.Position(method.Pos()),
				})
			}

			interfaces[typeSpec.Name.Name] = interfaceDeclaration{
				Name:     typeSpec.Name.Name,
				Methods:  methods,
				Position: fileSet.Position(typeSpec.Pos()),
			}
		}
	}
}

func collectMockMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	pathClassifier analysis.PathClassifier,
	mocks map[string][]methodDeclaration,
) {
	if !pathClassifier.HasRole(path, analysis.PathRoleTestMocks) {
		return
	}

	for _, declaration := range file.Decls {
		funcDecl, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}
		receiver := analysis.ReceiverTypeName(funcDecl.Recv.List[0].Type)
		if receiver == "" {
			continue
		}
		mocks[receiver] = append(mocks[receiver], methodDeclaration{
			Name:     funcDecl.Name.Name,
			Position: fileSet.Position(funcDecl.Name.Pos()),
		})
	}
}
