package checker

import (
	"go/ast"
	"go/token"
	"strings"
)

// collectInterfaces records interface method order from internal/core/ports.
func collectInterfaces(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	interfaces map[string]interfaceDecl,
) {
	if !strings.Contains(path, corePortsPathSegment) {
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

			methods := make([]methodDecl, 0, len(interfaceType.Methods.List))
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}
				methods = append(methods, methodDecl{
					name:     method.Names[0].Name,
					position: fileSet.Position(method.Pos()),
				})
			}

			interfaces[typeSpec.Name.Name] = interfaceDecl{
				name:     typeSpec.Name.Name,
				methods:  methods,
				position: fileSet.Position(typeSpec.Pos()),
			}
		}
	}
}

// collectMockMethods records method order for each mock receiver type.
func collectMockMethods(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	mocks map[string][]methodDecl,
) {
	if !strings.Contains(path, mocksPathSegment) {
		return
	}

	for _, declaration := range file.Decls {
		funcDecl, ok := declaration.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}
		receiver := receiverTypeName(funcDecl.Recv.List[0].Type)
		if receiver == "" {
			continue
		}
		mocks[receiver] = append(mocks[receiver], methodDecl{
			name:     funcDecl.Name.Name,
			position: fileSet.Position(funcDecl.Name.Pos()),
		})
	}
}
