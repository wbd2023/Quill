package syntax

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/wbd2023/quill/internal/checks/golang/analysis"
	"github.com/wbd2023/quill/internal/markers"
)

// weightless-suffix constants.
const (
	weightlessInstanceSuffix = "Instance"
	weightlessImplSuffix     = "Impl"
	weightlessWrapperSuffix  = "Wrapper"
	weightlessSuffixMarker   = "allow-weightless-suffix"
)

/* --------------------------------- Weightless Suffix Detection -------------------------------- */

// CheckWeightlessSuffixes flags declarations whose names end in a weightless suffix (`Instance`,
// `Impl`, or `Wrapper`). A capital letter immediately before the suffix, such as the `ID` in
// `UserIDImpl`, suppresses the finding. Declarations may opt out with an inline
// `style: allow-weightless-suffix` marker on the declaration line.
func CheckWeightlessSuffixes(
	fileSet *token.FileSet,
	file *ast.File,
	lines []string,
) (violations []analysis.Violation) {
	for _, declaration := range file.Decls {
		violations = append(
			violations,
			checkWeightlessTopLevelDeclaration(fileSet, lines, declaration)...,
		)
	}

	violations = append(violations, checkWeightlessMemberNames(fileSet, file, lines)...)

	return violations
}

func checkWeightlessTopLevelDeclaration(
	fileSet *token.FileSet,
	lines []string,
	declaration ast.Decl,
) (violations []analysis.Violation) {
	switch declaration := declaration.(type) {
	case *ast.GenDecl:
		return checkWeightlessSpecs(fileSet, lines, declaration)

	case *ast.FuncDecl:
		kind := "function"
		if declaration.Recv != nil {
			kind = "method"
		}

		return appendWeightlessSuffixViolation(violations, fileSet, lines, declaration.Name, kind)
	}

	return violations
}

func checkWeightlessSpecs(
	fileSet *token.FileSet,
	lines []string,
	declaration *ast.GenDecl,
) (violations []analysis.Violation) {
	for _, spec := range declaration.Specs {
		switch spec := spec.(type) {
		case *ast.TypeSpec:
			violations = appendWeightlessSuffixViolation(
				violations,
				fileSet,
				lines,
				spec.Name,
				"type",
			)

		case *ast.ValueSpec:
			kind := "variable"
			if declaration.Tok == token.CONST {
				kind = "constant"
			}

			for _, name := range spec.Names {
				violations = appendWeightlessSuffixViolation(
					violations,
					fileSet,
					lines,
					name,
					kind,
				)
			}
		}
	}

	return violations
}

func checkWeightlessMemberNames(
	fileSet *token.FileSet,
	file *ast.File,
	lines []string,
) (violations []analysis.Violation) {
	members := collectWeightlessMembers(file)

	violations = make([]analysis.Violation, 0, len(members))
	for _, member := range members {
		violations = appendWeightlessSuffixViolation(
			violations,
			fileSet,
			lines,
			member.identifier,
			member.kind,
		)
	}

	return violations
}

type weightlessMember struct {
	identifier *ast.Ident
	kind       string
}

func collectWeightlessMembers(file *ast.File) (members []weightlessMember) {
	ast.Inspect(file, func(node ast.Node) bool {
		switch layout := node.(type) {
		case *ast.StructType:
			for _, field := range layout.Fields.List {
				for _, name := range field.Names {
					members = append(members, weightlessMember{identifier: name, kind: "field"})
				}
			}

		case *ast.InterfaceType:
			for _, method := range layout.Methods.List {
				for _, name := range method.Names {
					members = append(members, weightlessMember{
						identifier: name,
						kind:       "interface method",
					})
				}
			}
		}

		return true
	})

	return members
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func appendWeightlessSuffixViolation(
	violations []analysis.Violation,
	fileSet *token.FileSet,
	lines []string,
	name *ast.Ident,
	kind string,
) (updated []analysis.Violation) {
	violation := weightlessSuffixViolation(fileSet, lines, name, kind)
	if violation != nil {
		return append(violations, *violation)
	}

	return violations
}

func weightlessSuffixViolation(
	fileSet *token.FileSet,
	lines []string,
	name *ast.Ident,
	kind string,
) (violation *analysis.Violation) {
	suffix, found := weightlessSuffixOf(name.Name)
	if !found ||
		allowsInlineException(lines, fileSet.Position(name.Pos()), weightlessSuffixMarker) {
		return nil
	}

	return &analysis.Violation{
		Position: fileSet.Position(name.Pos()),
		Rule:     analysis.DiagnosticWeightlessSuffix,
		Message:  fmt.Sprintf("weightless suffix %q in %s %q", suffix, kind, name.Name),
	}
}

func weightlessSuffixOf(name string) (suffix string, found bool) {
	candidates := []string{
		weightlessInstanceSuffix,
		weightlessImplSuffix,
		weightlessWrapperSuffix,
	}

	for _, candidate := range candidates {
		if !strings.HasSuffix(name, candidate) {
			continue
		}

		if hasLowercaseBefore(name, len(name)-len(candidate)) {
			return candidate, true
		}
	}

	return "", false
}

func hasLowercaseBefore(name string, index int) (valid bool) {
	if index <= 0 {
		return false
	}

	preceding := name[index-1]
	return 'a' <= preceding && preceding <= 'z'
}

func allowsInlineException(lines []string, position token.Position, rule string) (allowed bool) {
	lineIndex := position.Line - 1
	if lineIndex < 0 || lineIndex >= len(lines) {
		return false
	}

	return markers.HasMarker(lines[lineIndex], rule)
}
