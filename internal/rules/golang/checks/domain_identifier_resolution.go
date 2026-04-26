package checks

import (
	"go/types"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/policy"
)

func resolvedDomainIdentifierTypeName(
	targetType types.Type,
	classifier PathClassifier,
	identifiers policy.GoDomainIdentifierConfig,
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
	if !classifier.MatchesImportPath(packagePath, PathClassDomain) {
		return "", false
	}

	typeName := typeObject.Name()
	if _, supported := recommendedDomainIdentifierConstructor(identifiers, typeName); !supported {
		return "", false
	}

	return typeName, true
}

func recommendedDomainIdentifierConstructor(
	identifiers policy.GoDomainIdentifierConfig,
	typeName string,
) (constructor string, found bool) {
	constructors := identifiers[typeName]
	if len(constructors) == 0 {
		return "", false
	}

	return strings.Join(constructors, " or "), true
}

func normalisePath(path string) (normalisedPath string) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolutePath))
}
