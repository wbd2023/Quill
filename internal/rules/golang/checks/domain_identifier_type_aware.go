package checks

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"

	"ciphera/tools/internal/policy"
)

/* --------------------------------------- Type-Aware Scan -------------------------------------- */

func CollectTypeAwareDomainIdentifierCastViolations(
	rootDirectories []string,
	filePaths []string,
	classifier PathClassifier,
	constructors policy.GoDomainIdentifierConstructors,
) (violations []Violation, ran bool) {
	if len(filePaths) == 0 || len(rootDirectories) == 0 {
		return nil, false
	}

	requestedFilePaths := requestedGoFiles(filePaths)
	for _, rootDirectory := range rootDirectories {
		rootViolations, rootRan := collectTypeAwareViolationsInRoot(
			rootDirectory,
			requestedFilePaths,
			classifier,
			constructors,
		)
		if !rootRan {
			continue
		}

		ran = true
		violations = append(violations, rootViolations...)
	}

	return violations, ran
}

func collectTypeAwareViolationsInRoot(
	rootDirectory string,
	requestedFilePaths map[string]bool,
	classifier PathClassifier,
	constructors policy.GoDomainIdentifierConstructors,
) (violations []Violation, ran bool) {
	packageList, err := packages.Load(typeAwarePackageConfig(rootDirectory), "./...")
	if err != nil || len(packageList) == 0 {
		return nil, false
	}

	for _, packageInfo := range packageList {
		violations = append(violations, collectTypeAwareViolationsInPackage(
			packageInfo,
			requestedFilePaths,
			classifier,
			constructors,
		)...)
	}

	return violations, true
}

func typeAwarePackageConfig(rootDirectory string) (config *packages.Config) {
	return &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Dir:   normalisePath(rootDirectory),
		Tests: true,
	}
}

func requestedGoFiles(filePaths []string) (requested map[string]bool) {
	requested = make(map[string]bool, len(filePaths))
	for _, filePath := range filePaths {
		requested[normalisePath(filePath)] = true
	}

	return requested
}

/* ---------------------------------------- Package Scan ---------------------------------------- */

func collectTypeAwareViolationsInPackage(
	packageInfo *packages.Package,
	requestedFilePaths map[string]bool,
	classifier PathClassifier,
	constructors policy.GoDomainIdentifierConstructors,
) (violations []Violation) {
	if packageInfo == nil || packageInfo.TypesInfo == nil || packageInfo.Fset == nil {
		return nil
	}

	for _, file := range packageInfo.Syntax {
		filePath := normalisePath(packageInfo.Fset.Position(file.Pos()).Filename)
		if !requestedFilePaths[filePath] || classifier.HasClass(filePath, PathClassDomain) {
			continue
		}

		violations = append(violations, collectTypeAwareViolationsInFile(
			packageInfo,
			file,
			classifier,
			constructors,
		)...)
	}

	return violations
}

func collectTypeAwareViolationsInFile(
	packageInfo *packages.Package,
	file *ast.File,
	classifier PathClassifier,
	constructors policy.GoDomainIdentifierConstructors,
) (violations []Violation) {
	ast.Inspect(file, func(node ast.Node) bool {
		callExpression, ok := node.(*ast.CallExpr)
		if !ok || len(callExpression.Args) != 1 {
			return true
		}

		violation, found := typeAwareDomainIdentifierViolation(
			packageInfo,
			callExpression,
			classifier,
			constructors,
		)
		if found {
			violations = append(violations, violation)
		}

		return true
	})

	return violations
}

func typeAwareDomainIdentifierViolation(
	packageInfo *packages.Package,
	callExpression *ast.CallExpr,
	classifier PathClassifier,
	constructors policy.GoDomainIdentifierConstructors,
) (violation Violation, found bool) {
	functionInfo, ok := packageInfo.TypesInfo.Types[callExpression.Fun]
	if !ok {
		return Violation{}, false
	}

	domainTypeName, found := resolvedDomainIdentifierTypeName(
		functionInfo.Type,
		classifier,
		constructors,
	)
	if !found {
		return Violation{}, false
	}

	recommendedConstructor, _ := recommendedDomainIdentifierConstructor(
		constructors,
		domainTypeName,
	)
	return Violation{
		Position: packageInfo.Fset.Position(callExpression.Pos()),
		Rule:     DiagnosticNoDirectDomainCasts,
		Message: fmt.Sprintf(
			"direct cast to domain.%s is disallowed; use %s",
			domainTypeName,
			recommendedConstructor,
		),
	}, true
}
