package text

import (
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
)

/* ------------------------------------------- Markers ------------------------------------------ */

func CheckMaintenanceMarkers(
	repoRoot string,
	repository policy.RepositoryConfig,
	scope contract.Scope,
) (result contract.ExecutionResult, err error) {
	files, err := filewalk.CollectAllFiles(repoRoot, repository, scope)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	for _, path := range files {
		if !supportsMaintenanceMarkers(path) {
			continue
		}

		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			if !strings.Contains(line.Text, "TODO:") && !strings.Contains(line.Text, "FIXME:") {
				return nil
			}

			if markerHasContent(line.Text) {
				return nil
			}

			result.Diagnostics = append(result.Diagnostics, contract.Diagnostic{
				Code:    "text/maintenance-markers/missing-action",
				File:    filewalk.RelativePath(repoRoot, path),
				Line:    line.Number,
				Message: "TODO/FIXME markers must include actionable text after the colon",
			})
			return nil
		})
		if err != nil {
			return contract.ExecutionResult{}, err
		}
	}

	if len(result.Diagnostics) == 0 {
		return contract.ExecutionResult{}, nil
	}

	return result, contract.ViolationsFound()
}

/* -------------------------------------- Marker Detection -------------------------------------- */

func supportsMaintenanceMarkers(path string) (supported bool) {
	base := filepath.Base(path)

	switch base {
	case "Makefile":
		return true
	}

	switch filepath.Ext(path) {
	case ".go", ".sh", ".md", ".txt", ".yml", ".yaml", ".json", ".toml":
		return true
	default:
		return false
	}
}

func markerHasContent(line string) (found bool) {
	for _, marker := range []string{"TODO:", "FIXME:"} {
		index := strings.Index(line, marker)
		if index < 0 {
			continue
		}

		return strings.TrimSpace(line[index+len(marker):]) != ""
	}

	return true
}
