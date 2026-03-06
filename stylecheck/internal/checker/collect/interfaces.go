package collect

import (
	"go/ast"
	"go/token"
	"strings"

	"stylecheck/internal/checker/support"
)

// CollectInterfaces records interface method order from internal/core/ports.
func CollectInterfaces(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	interfaces map[string]InterfaceDecl,
) {
	if !strings.Contains(path, support.CorePortsPathSegment) {
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

			methods := make([]MethodDecl, 0, len(interfaceType.Methods.List))
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}
				methods = append(methods, MethodDecl{
					Name:     method.Names[0].Name,
					Position: fileSet.Position(method.Pos()),
				})
			}

			interfaces[typeSpec.Name.Name] = InterfaceDecl{
				Name:     typeSpec.Name.Name,
				Methods:  methods,
				Position: fileSet.Position(typeSpec.Pos()),
			}
		}
	}
}

// CollectMockMethods records method order for each mock receiver type.
func CollectMockMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	mocks map[string][]MethodDecl,
) {
	if !strings.Contains(path, support.MocksPathSegment) {
		return
	}

	for _, declaration := range file.Decls {
		funcDecl, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}
		receiver := support.ReceiverTypeName(funcDecl.Recv.List[0].Type)
		if receiver == "" {
			continue
		}
		mocks[receiver] = append(mocks[receiver], MethodDecl{
			Name:     funcDecl.Name.Name,
			Position: fileSet.Position(funcDecl.Name.Pos()),
		})
	}
}
