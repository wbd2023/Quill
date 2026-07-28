package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

/* ----------------------------------- Repository Containment ----------------------------------- */

// Repository containment policy.
//
// A Quill Profile is a trust boundary: its path values are decoded from repository-owned
// quill.toml before any filesystem use. The functions below enforce that every
// repository-relative path value resolves to a location inside the canonical repository root,
// closing absolute, NUL, parent-traversal, and symlink escapes. They use only the standard
// library; there is no virtual filesystem.

// CanonicalRoot returns the absolute, symlink-resolved, cleaned repository root. The engine
// resolves the discovered root once and threads the canonical value through every operation so
// that downstream joins compare against a stable physical root.
func CanonicalRoot(root string) (canonical string, err error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}

	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", fmt.Errorf("resolve repository root %q: %w", absolute, err)
	}

	return filepath.Clean(resolved), nil
}

// ValidateRepoRelative reports whether value is a lexically contained repository-relative path.
// It rejects absolute paths, Windows drive roots, NUL bytes, and any value whose cleaned form
// escapes the repository root via parent traversal. The semantic root value "." is accepted.
// This is the lexical half of the containment policy; it performs no filesystem access.
func ValidateRepoRelative(value string) (err error) {
	if strings.ContainsRune(value, '\x00') {
		return fmt.Errorf("path %q contains a NUL byte", value)
	}

	normalized := strings.ReplaceAll(value, "\\", "/")
	if strings.HasPrefix(normalized, "/") || isWindowsDriveRoot(normalized) {
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

// ResolveRepoRelative joins value under root and enforces lexical and physical containment.
// root must be canonical (see CanonicalRoot). When the joined target exists, its symlinks are
// resolved and the result must remain inside root; a non-existent target is accepted because its
// escape is already ruled out by the lexical check and resolution is not yet possible.
func ResolveRepoRelative(root string, value string) (absolute string, err error) {
	if err = ValidateRepoRelative(value); err != nil {
		return "", err
	}

	joined := filepath.Join(root, filepath.FromSlash(value))
	resolved, err := filepath.EvalSymlinks(joined)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return joined, nil
		}
		return "", fmt.Errorf("resolve repository path %q: %w", value, err)
	}

	if !isWithinRoot(root, resolved) {
		return "", fmt.Errorf("path %q resolves outside the repository root", value)
	}

	return resolved, nil
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// isWithinRoot reports whether target is root itself or a descendant of root after cleaning.
func isWithinRoot(root string, target string) (within bool) {
	cleanedRoot := filepath.Clean(root)
	rel, err := filepath.Rel(cleanedRoot, filepath.Clean(target))
	if err != nil {
		return false
	}

	rel = filepath.ToSlash(rel)
	return rel != ".." && !strings.HasPrefix(rel, "../")
}

func isWindowsDriveRoot(value string) (drive bool) {
	if len(value) < 2 || value[1] != ':' {
		return false
	}

	head := value[0]
	return head >= 'A' && head <= 'Z' || head >= 'a' && head <= 'z'
}
