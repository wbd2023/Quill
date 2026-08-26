package bash

import (
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/filewalk"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

/* ------------------------------------------ Constants ----------------------------------------- */

// structure constants.
const (
	shellShebangLine = "#!/bin/bash"
	strictModeLine   = "set -euo pipefail"
)

/* --------------------------------------- Structure Rules -------------------------------------- */

// CheckStructure check structure.
func CheckStructure(
	root string,
	repository profile.RepositoryConfig,
	scope style.Scope,
) (result style.ExecutionResult, err error) {
	files, err := filewalk.CollectFiles(
		repository.ResolveScopeRoots(root, scope),
		walkConfig(repository),
		".sh",
	)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	for _, path := range files {
		hasStrictMode := false
		err = filewalk.ScanLines(path, func(line filewalk.Line) error {
			text := line.Text

			if line.Number == 1 && text != shellShebangLine {
				result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
					Code:    "bash/structure/invalid",
					File:    filewalk.DisplayPath(root, path),
					Range:   style.Range{Start: style.Position{Line: 1}},
					Message: fmt.Sprintf("missing %s", shellShebangLine),
				})
			}

			if text == strictModeLine {
				hasStrictMode = true
			}

			if strings.Contains(text, "\r") {
				result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
					Code:    "bash/structure/invalid",
					File:    filewalk.DisplayPath(root, path),
					Range:   style.Range{Start: style.Position{Line: line.Number}},
					Message: "contains CRLF line endings",
				})
			}

			if hasTrailingHorizontalWhitespace(text) {
				result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
					Code:    "bash/structure/invalid",
					File:    filewalk.DisplayPath(root, path),
					Range:   style.Range{Start: style.Position{Line: line.Number}},
					Message: "has trailing whitespace",
				})
			}

			if hasSpaceIndentation(text) {
				result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
					Code:    "bash/structure/invalid",
					File:    filewalk.DisplayPath(root, path),
					Range:   style.Range{Start: style.Position{Line: line.Number}},
					Message: "uses space indentation",
				})
			}

			return nil
		})
		if err != nil {
			return style.ExecutionResult{}, err
		}

		if !hasStrictMode {
			result.Diagnostics = append(result.Diagnostics, style.Diagnostic{
				Code:    "bash/structure/invalid",
				File:    filewalk.DisplayPath(root, path),
				Message: fmt.Sprintf("missing %s", strictModeLine),
			})
		}
	}

	if len(result.Diagnostics) == 0 {
		return style.ExecutionResult{}, nil
	}

	return result, nil
}

func hasSpaceIndentation(line string) (found bool) {
	hasLeadingSpaces := false

	for _, character := range line {
		switch character {
		case ' ':
			hasLeadingSpaces = true
		case '\t':
			return false
		default:
			return hasLeadingSpaces
		}
	}

	return false
}

func hasTrailingHorizontalWhitespace(line string) (found bool) {
	if line == "" {
		return false
	}

	lastByte := line[len(line)-1]
	return lastByte == ' ' || lastByte == '\t'
}
