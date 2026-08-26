package profile

import (
	"fmt"
	"path/filepath"
	"strings"
)

// validateRepoPath enforces lexical repository containment for one Profile path value. It is the
// Profile-boundary mirror of workspace.ValidateRepoRelative: absolute paths, Windows drive
// roots, NUL bytes, and parent-traversal escapes are rejected using only the standard library.
//
// Profile is an architectural trust boundary and cannot import the workspace resolver, so the
// lexical policy is applied here. Physical (symlink) containment for resolvable paths is then
// enforced by validateRepositoryPaths during Load, where the repository root is known.
func validateRepoPath(value string) (err error) {
	if strings.ContainsRune(value, '\x00') {
		return fmt.Errorf("path %q contains a NUL byte", value)
	}

	normalized := strings.ReplaceAll(value, "\\", "/")
	if strings.HasPrefix(normalized, "/") || isWindowsDrivePath(normalized) {
		return fmt.Errorf("path %q must be repository-relative, not absolute", value)
	}

	cleaned := filepath.ToSlash(filepath.Clean(normalized))
	if cleaned == "." {
		return nil
	}

	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("path %q escapes the repository root", value)
	}

	return nil
}

func isWindowsDrivePath(value string) (drive bool) {
	if len(value) < 2 || value[1] != ':' {
		return false
	}

	head := value[0]
	return head >= 'A' && head <= 'Z' || head >= 'a' && head <= 'z'
}
