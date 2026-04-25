package checks

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
)

/* --------------------------------------- Domain ID Rules -------------------------------------- */

// CheckDirectDomainIdentifierCasts enforces parser/constructor usage for key domain IDs.
func CheckDirectDomainIdentifierCasts(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	classifier PathClassifier,
	identifiers profile.GoDomainIdentifierConfig,
) (violations []Violation) {
	if classifier.HasClass(path, rulepack.PathClassDomain) {
		return nil
	}

	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok || len(callExpression.Args) != 1 {
			return true
		}

		selector, ok := callExpression.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		packageIdentifier, ok := selector.X.(*ast.Ident)
		if !ok || packageIdentifier.Name != "domain" {
			return true
		}

		recommendedConstructor, found := recommendedDomainIdentifierConstructor(
			identifiers,
			selector.Sel.Name,
		)
		if !found {
			return true
		}

		violations = append(violations, Violation{
			Position: fileSet.Position(callExpression.Pos()),
			Rule:     DiagnosticNoDirectDomainCasts,
			Message: fmt.Sprintf(
				"direct cast to domain.%s is disallowed; use %s",
				selector.Sel.Name,
				recommendedConstructor,
			),
		})
		return true
	})

	return violations
}

func CollectTypeAwareDomainIdentifierCastViolations(
	rootDirectories []string,
	filePaths []string,
	classifier PathClassifier,
	identifiers profile.GoDomainIdentifierConfig,
) (violations []Violation, ran bool) {
	if len(filePaths) == 0 || len(rootDirectories) == 0 {
		return nil, false
	}

	requestedFilePaths := make(map[string]bool, len(filePaths))
	for _, filePath := range filePaths {
		normalisedPath := normalisePath(filePath)
		requestedFilePaths[normalisedPath] = true
	}

	for _, rootDirectory := range rootDirectories {
		normalisedRoot := normalisePath(rootDirectory)

		packageConfig := &packages.Config{
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedSyntax |
				packages.NeedTypes |
				packages.NeedTypesInfo,
			Dir:   normalisedRoot,
			Tests: true,
		}

		packageList, err := packages.Load(packageConfig, "./...")
		if err != nil || len(packageList) == 0 {
			continue
		}

		ran = true

		for _, packageInfo := range packageList {
			if packageInfo == nil ||
				packageInfo.TypesInfo == nil ||
				packageInfo.Fset == nil {
				continue
			}

			for _, file := range packageInfo.Syntax {
				filePath := normalisePath(packageInfo.Fset.Position(file.Pos()).Filename)

				if !requestedFilePaths[filePath] {
					continue
				}

				if classifier.HasClass(filePath, rulepack.PathClassDomain) {
					continue
				}

				ast.Inspect(file, func(node ast.Node) bool {
					callExpression, ok := node.(*ast.CallExpr)
					if !ok || len(callExpression.Args) != 1 {
						return true
					}

					typeAndValue, ok := packageInfo.TypesInfo.Types[callExpression.Fun]
					if !ok {
						return true
					}

					domainTypeName, found := resolvedDomainIdentifierTypeName(
						typeAndValue.Type,
						classifier,
						identifiers,
					)
					if !found {
						return true
					}

					recommendedConstructor, _ := recommendedDomainIdentifierConstructor(
						identifiers,
						domainTypeName,
					)
					violations = append(violations, Violation{
						Position: packageInfo.Fset.Position(callExpression.Pos()),
						Rule:     DiagnosticNoDirectDomainCasts,
						Message: fmt.Sprintf(
							"direct cast to domain.%s is disallowed; use %s",
							domainTypeName,
							recommendedConstructor,
						),
					})

					return true
				})
			}
		}
	}

	return violations, ran
}

func resolvedDomainIdentifierTypeName(
	targetType types.Type,
	classifier PathClassifier,
	identifiers profile.GoDomainIdentifierConfig,
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
	if !classifier.MatchesImportPath(packagePath, rulepack.PathClassDomain) {
		return "", false
	}

	typeName := typeObject.Name()
	if _, supported := recommendedDomainIdentifierConstructor(identifiers, typeName); !supported {
		return "", false
	}

	return typeName, true
}

func normalisePath(path string) (normalisedPath string) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}

	return filepath.ToSlash(filepath.Clean(absolutePath))
}

func recommendedDomainIdentifierConstructor(
	identifiers profile.GoDomainIdentifierConfig,
	typeName string,
) (constructor string, found bool) {
	constructors := identifiers[typeName]
	if len(constructors) == 0 {
		return "", false
	}

	return strings.Join(constructors, " or "), true
}
