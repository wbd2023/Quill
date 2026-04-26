package checks

import (
	"go/ast"
	"go/token"
)

func collectInterfaceTypeNames(file *ast.File) (names map[string]bool) {
	names = make(map[string]bool)

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
			if _, ok := typeSpec.Type.(*ast.InterfaceType); !ok {
				continue
			}

			names[typeSpec.Name.Name] = true
		}
	}

	return names
}

func collectSliceNames(file *ast.File) (names map[string]bool) {
	names = make(map[string]bool)

	for _, declaration := range file.Decls {
		switch typedDecl := declaration.(type) {
		case *ast.FuncDecl:
			collectSliceNamesFromFieldList(names, typedDecl.Type.Params)
			collectSliceNamesFromFieldList(names, typedDecl.Type.Results)

		case *ast.GenDecl:
			for _, spec := range typedDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok || !isSliceType(valueSpec.Type) {
					continue
				}

				for _, name := range valueSpec.Names {
					names[name.Name] = true
				}
			}
		}
	}

	return names
}

func collectSliceNamesFromFieldList(names map[string]bool, fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, field := range fields.List {
		if !isSliceType(field.Type) {
			continue
		}

		for _, name := range field.Names {
			names[name.Name] = true
		}
	}
}

func isSliceType(expression ast.Expr) (found bool) {
	_, ok := expression.(*ast.ArrayType)
	return ok
}
