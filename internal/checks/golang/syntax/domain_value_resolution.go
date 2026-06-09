package syntax

import (
	"go/types"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/checks/golang/analysis"
	gopolicy "ciphera/tools/internal/checks/golang/policy"
)

func resolvedDomainValueTypeName(
	targetType types.Type,
	classifier analysis.PathClassifier,
	constructors gopolicy.DomainValueConstructors,
) (name string, found bool) {
	namedType, ok := types.Unalias(targetType).(*types.Named)
	if !ok {
		return "", false
	}

	typeObject := namedType.Obj()
	if typeObject == nil || typeObject.Pkg() == nil {
		return "", false
	}

	packagePath := typeObject.Pkg().Path()
	if !classifier.MatchesImportPath(packagePath, analysis.PathRoleDomain) {
		return "", false
	}

	typeName := typeObject.Name()
	if _, supported := recommendedDomainValueConstructor(constructors, typeName); !supported {
		return "", false
	}

	return typeName, true
}

func recommendedDomainValueConstructor(
	constructors gopolicy.DomainValueConstructors,
	typeName string,
) (constructor string, found bool) {
	names := constructors[typeName]
	if len(names) == 0 {
		return "", false
	}

	return strings.Join(names, " or "), true
}

func normalisePath(path string) (normalisedPath string) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolutePath))
}
