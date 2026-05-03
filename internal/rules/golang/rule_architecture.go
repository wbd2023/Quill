package golang

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

/* -------------------------------------- Package Listings -------------------------------------- */

type listedPackage struct {
	ImportPath string   `json:"ImportPath"`
	Imports    []string `json:"Imports"`
}

/* ------------------------------------- Architecture Rules ------------------------------------- */

func CheckArchitecture(
	modulePath string,
	packageList string,
	architecture policy.GoArchitectureConfig,
) (result contract.ExecutionResult, err error) {
	modulePath = strings.TrimSpace(modulePath)
	decoder := json.NewDecoder(strings.NewReader(packageList))
	diagnostics := make([]contract.Diagnostic, 0)

	for decoder.More() {
		var packageEntry listedPackage
		if decodeErr := decoder.Decode(&packageEntry); decodeErr != nil {
			return contract.ExecutionResult{}, decodeErr
		}

		fromLayer := classifyPackage(
			modulePath,
			packageEntry.ImportPath,
			architecture,
		)
		if fromLayer == "" {
			continue
		}

		for _, importPath := range packageEntry.Imports {
			toLayer := classifyPackage(modulePath, importPath, architecture)
			if toLayer == "" {
				continue
			}

			if isAllowedImport(architecture, fromLayer, toLayer) {
				continue
			}

			diagnostics = append(diagnostics, contract.Diagnostic{
				Code: "go/architecture/import-boundary",
				Message: fmt.Sprintf(
					"%s [%s] imports %s [%s]",
					packageEntry.ImportPath,
					fromLayer,
					importPath,
					toLayer,
				),
			})
		}
	}

	if len(diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return contract.ExecutionResult{Diagnostics: diagnostics}, errViolationsFound
}

/* ----------------------------------- Package Classification ----------------------------------- */

func classifyPackage(
	modulePath string,
	importPath string,
	architecture policy.GoArchitectureConfig,
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
	architecture policy.GoArchitectureConfig,
	fromLayer string,
	toLayer string,
) (allowed bool) {
	for _, layer := range architecture.Layers {
		if layer.Name != fromLayer {
			continue
		}

		return slices.Contains(layer.AllowedLayers, toLayer)
	}

	return true
}
