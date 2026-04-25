package repostyle

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
)

/* ------------------------------------------- Markers ------------------------------------------ */

func CheckMaintenanceMarkers(
	repoRoot string,
	repository profile.RepositoryConfig,
	scope contract.Scope,
) (output string, err error) {
	files, err := CollectAllFiles(repoRoot, repository, scope)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	foundViolation := false

	for _, path := range files {
		if !supportsMaintenanceMarkers(path) {
			continue
		}

		file, openErr := os.Open(path)
		if openErr != nil {
			return "", openErr
		}

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()

			if !containsMarker(line, "TODO:") && !containsMarker(line, "FIXME:") {
				continue
			}

			if markerHasContent(line) {
				continue
			}

			foundViolation = true
			builder.WriteString(fmt.Sprintf(
				"%s:%d TODO/FIXME markers must include actionable text after the colon\n",
				RelativePath(repoRoot, path),
				lineNumber,
			))
		}

		if scanErr := scanner.Err(); scanErr != nil {
			return "", closeFile(file, scanErr)
		}

		if closeErr := closeFile(file, nil); closeErr != nil {
			return "", closeErr
		}
	}

	if !foundViolation {
		return "", nil
	}

	return builder.String(), errViolationsFound
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

func containsMarker(line string, marker string) (found bool) {
	return strings.Contains(line, marker)
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
