package checker

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

/* ------------------------------------------ Constants ----------------------------------------- */

var directDomainIdentifierConstructors = map[string]string{
	"Username":       "ParseUsername",
	"ConversationID": "ParseConversationID or ConversationIDFromUsername",
	"IdentityID":     "ParseIdentityID",
}

/* --------------------------------------- Domain ID Rules -------------------------------------- */

// checkDirectDomainIdentifierCasts enforces parser/constructor usage for key domain IDs (2.2).
func checkDirectDomainIdentifierCasts(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if strings.Contains(path, domainPathSegment) {
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

		recommendedConstructor, found := directDomainIdentifierConstructors[selector.Sel.Name]
		if !found {
			return true
		}

		violations = append(violations, violation{
			position: fileSet.Position(callExpression.Pos()),
			rule:     "2.2",
			message: fmt.Sprintf(
				"direct cast to domain.%s is disallowed; use %s",
				selector.Sel.Name,
				recommendedConstructor,
			),
		})
		return true
	})

	return violations
}

func collectTypeAwareDomainIdentifierCastViolations(
	rootDirectories []string,
	filePaths []string,
) (violations []violation, ran bool) {
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

		loadedPackages, err := packages.Load(packageConfig, "./...")
		if err != nil || len(loadedPackages) == 0 {
			continue
		}

		ran = true

		for _, loadedPackage := range loadedPackages {
			if loadedPackage == nil ||
				loadedPackage.TypesInfo == nil ||
				loadedPackage.Fset == nil {
				continue
			}

			for _, file := range loadedPackage.Syntax {
				filePath := normalisePath(loadedPackage.Fset.Position(file.Pos()).Filename)

				if !requestedFilePaths[filePath] {
					continue
				}

				if strings.Contains(filePath, domainPathSegment) {
					continue
				}

				ast.Inspect(file, func(node ast.Node) bool {
					callExpression, ok := node.(*ast.CallExpr)
					if !ok || len(callExpression.Args) != 1 {
						return true
					}

					typeAndValue, ok := loadedPackage.TypesInfo.Types[callExpression.Fun]
					if !ok {
						return true
					}

					domainTypeName, found := resolvedDomainIdentifierTypeName(typeAndValue.Type)
					if !found {
						return true
					}

					recommendedConstructor := directDomainIdentifierConstructors[domainTypeName]
					violations = append(violations, violation{
						position: loadedPackage.Fset.Position(callExpression.Pos()),
						rule:     "2.2",
						message: fmt.Sprintf(
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

func resolvedDomainIdentifierTypeName(targetType types.Type) (name string, found bool) {
	namedType, ok := types.Unalias(targetType).(*types.Named)
	if !ok {
		return "", false
	}

	typeObject := namedType.Obj()
	if typeObject == nil || typeObject.Pkg() == nil {
		return "", false
	}

	packagePath := typeObject.Pkg().Path()
	if !isDomainPackagePath(packagePath) {
		return "", false
	}

	typeName := typeObject.Name()
	if _, supported := directDomainIdentifierConstructors[typeName]; !supported {
		return "", false
	}

	return typeName, true
}

func isDomainPackagePath(packagePath string) (found bool) {
	if packagePath == "internal/core/domain" {
		return true
	}

	return strings.HasSuffix(packagePath, domainPathSuffix)
}
