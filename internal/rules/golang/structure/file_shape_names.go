package structure

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"ciphera/tools/internal/rules/golang/analysis"
)

var vagueGoFileNames = map[string]bool{
	"checks.go":  true,
	"syntax.go":  true,
	"helpers.go": true,
	"model.go":   true,
	"types.go":   true,
}

func vagueFileNameViolations(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []analysis.Violation) {
	if !vagueFileNameNeedsReview(file, path) {
		return nil
	}

	return []analysis.Violation{{
		Position: fileSet.Position(file.Package),
		Rule:     analysis.DiagnosticFileShapeVagueName,
		Message: fmt.Sprintf(
			"file name %q is vague; use a role-specific file name "+
				"unless the file is package-wide",
			filepath.Base(path),
		),
	}}
}

func vagueFileNameNeedsReview(file *ast.File, path string) (needsReview bool) {
	name := filepath.Base(path)
	if !vagueGoFileNames[name] {
		return false
	}

	return !isPackageWideShape(file, name)
}

func isPackageWideShape(file *ast.File, name string) (packageWide bool) {
	switch name {
	case "types.go":
		return onlyPackageDataDeclarations(file)
	}

	return false
}

func onlyPackageDataDeclarations(file *ast.File) (onlyData bool) {
	if len(file.Decls) == 0 {
		return false
	}

	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok {
			return false
		}

		switch general.Tok {
		case token.CONST, token.IMPORT, token.TYPE:
		default:
			return false
		}
	}

	return true
}
