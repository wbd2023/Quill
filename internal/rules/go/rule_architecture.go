package gostyle

import (
	"encoding/json"
	"fmt"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/runtime"
)

/* -------------------------------------- Package Listings -------------------------------------- */

type listedPackage struct {
	ImportPath string   `json:"ImportPath"`
	Imports    []string `json:"Imports"`
}

/* ------------------------------------- Architecture Rules ------------------------------------- */

func CheckArchitecture(
	repoRoot string,
	scope contract.Scope,
	policy profile.Profile,
) (output string, err error) {
	if scope == contract.ScopeTools {
		return "", nil
	}

	layout := runtime.LayoutForRepository(repoRoot)
	goEnvironment := layout.GoEnvironment()

	modulePathOutput, err := runtime.RunCommand(
		repoRoot,
		goEnvironment,
		"go",
		"list",
		"-m",
		"-f",
		"{{.Path}}",
	)
	if err != nil {
		return modulePathOutput, err
	}
	modulePath := strings.TrimSpace(modulePathOutput)

	listOutput, err := runtime.RunCommand(
		repoRoot,
		goEnvironment,
		"go",
		"list",
		"-json",
		"./...",
	)
	if err != nil {
		return listOutput, err
	}

	decoder := json.NewDecoder(strings.NewReader(listOutput))
	violations := make([]string, 0)

	for decoder.More() {
		var packageEntry listedPackage
		if decodeErr := decoder.Decode(&packageEntry); decodeErr != nil {
			return "", decodeErr
		}

		fromLayer := classifyPackage(
			modulePath,
			packageEntry.ImportPath,
			policy.Architecture,
		)
		if fromLayer == "" {
			continue
		}

		for _, importPath := range packageEntry.Imports {
			toLayer := classifyPackage(modulePath, importPath, policy.Architecture)
			if toLayer == "" {
				continue
			}

			if isAllowedImport(policy.Architecture, fromLayer, toLayer) {
				continue
			}

			violations = append(violations, fmt.Sprintf(
				"%s [%s] imports %s [%s]",
				packageEntry.ImportPath,
				fromLayer,
				importPath,
				toLayer,
			))
		}
	}

	if len(violations) == 0 {
		return "", nil
	}

	return strings.Join(violations, "\n") + "\n", errViolationsFound
}

/* ----------------------------------- Package Classification ----------------------------------- */

func classifyPackage(
	modulePath string,
	importPath string,
	architecture profile.ArchitectureConfig,
) (layerName string) {
	relativePath, found := trimModulePrefix(modulePath, importPath)
	if !found {
		return ""
	}

	for _, layer := range architecture.Layers {
		if matchesPackageRoots(relativePath, layer.PackageRoots) {
			return layer.Name
		}
	}

	return ""
}

func trimModulePrefix(modulePath string, importPath string) (relativePath string, found bool) {
	if importPath == modulePath {
		return "", false
	}

	prefix := modulePath + "/"
	if !strings.HasPrefix(importPath, prefix) {
		return "", false
	}

	return strings.TrimPrefix(importPath, prefix), true
}

func matchesPackageRoots(relativePath string, packageRoots []string) (found bool) {
	for _, root := range packageRoots {
		trimmedRoot := strings.TrimSuffix(root, "/")
		if relativePath == trimmedRoot || strings.HasPrefix(relativePath, trimmedRoot+"/") {
			return true
		}
	}

	return false
}

/* ---------------------------------------- Import Rules ---------------------------------------- */

func isAllowedImport(
	architecture profile.ArchitectureConfig,
	fromLayer string,
	toLayer string,
) (allowed bool) {
	for _, layer := range architecture.Layers {
		if layer.Name != fromLayer {
			continue
		}

		for _, allowedLayer := range layer.MayImport {
			if allowedLayer == toLayer {
				return true
			}
		}

		return false
	}

	return true
}
