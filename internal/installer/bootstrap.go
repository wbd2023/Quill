package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/workspace"
)

/* ----------------------------------- Trusted Bootstrap Path ----------------------------------- */

// bootstrapPath builds the PATH used to resolve the trusted host Go and npm toolchain that
// bootstrap installation. It starts from the ambient host PATH and removes any entry that is
// equal to or beneath the repository root or the Quill state directory, comparing both the
// lexical absolute path and the symlink-resolved physical path. Ambient order is preserved.
//
// This is deliberately separate from Layout.BuildPath, which keeps the cache-first PATH used for
// normal managed-tool execution. Bootstrap tooling must never resolve from repository or state
// directories, where a checked-out or cached executable could otherwise become the bootstrap
// itself. The result is required to be non-empty: a host with no trusted directory cannot
// bootstrap installation and fails closed instead of executing an untrusted candidate.
func bootstrapPath(layout workspace.Layout) (path string, err error) {
	roots := bootstrapExclusions(layout)
	filtered, err := filterBootstrapPath(os.Getenv("PATH"), roots)
	if err != nil {
		return "", err
	}

	if filtered == "" {
		return "", fmt.Errorf(
			"bootstrap PATH has no trusted host directories after excluding %s and %s",
			layout.Root,
			layout.StateDirectory,
		)
	}

	return filtered, nil
}

// bootstrapExclusions returns the absolute directory roots that bootstrap PATH resolution must
// never select from: the repository root and the Quill state directory. The state directory is
// listed separately so a custom state directory outside the repository is excluded even when it
// is not lexically beneath the root.
const bootstrapExclusionCount = 2

func bootstrapExclusions(layout workspace.Layout) (roots []string) {
	roots = make([]string, 0, bootstrapExclusionCount)
	if root, err := filepath.Abs(layout.Root); err == nil {
		roots = append(roots, filepath.Clean(root))
	}
	if state, err := filepath.Abs(layout.StateDirectory); err == nil {
		roots = append(roots, filepath.Clean(state))
	}

	return roots
}

/* ------------------------------------ Path Entry Filtering ------------------------------------ */

// filterBootstrapPath returns canonical physical entries from path that are not equal to or
// beneath any root, in their original order. Every retained entry is absolute and symlink-resolved:
// a relative entry cannot be reinterpreted from the child working directory, and a later executable
// lookup cannot follow a retained directory link into repository or state. Missing entries are
// omitted so they cannot become search locations later; other resolution failures are rejected.
func filterBootstrapPath(path string, roots []string) (filtered string, err error) {
	resolvedRoots := resolveSymlinkRoots(roots)

	var kept []string
	for _, directory := range filepath.SplitList(path) {
		if directory == "" {
			continue
		}

		canonical, excluded, entryErr := canonicalBootstrapEntry(
			directory,
			roots,
			resolvedRoots,
		)
		if entryErr != nil {
			return "", entryErr
		}
		if excluded {
			continue
		}

		kept = append(kept, canonical)
	}

	return strings.Join(kept, string(os.PathListSeparator)), nil
}

// resolveSymlinkRoots returns the symlink-resolved form of each root, falling back to the lexical
// root when resolution fails (for example a root that does not yet exist).
func resolveSymlinkRoots(roots []string) (resolved []string) {
	resolved = make([]string, 0, len(roots))
	for _, root := range roots {
		if physical, err := filepath.EvalSymlinks(root); err == nil {
			resolved = append(resolved, filepath.Clean(physical))
		} else {
			resolved = append(resolved, root)
		}
	}

	return resolved
}

// canonicalBootstrapEntry resolves directory to its physical, absolute form and reports whether
// that form must be excluded from bootstrap PATH.
func canonicalBootstrapEntry(
	directory string,
	roots []string,
	resolvedRoots []string,
) (canonical string, excluded bool, err error) {
	lexical, err := filepath.Abs(directory)
	if err != nil {
		return "", false, fmt.Errorf("resolve bootstrap PATH entry %q: %w", directory, err)
	}
	lexical = filepath.Clean(lexical)

	if isWithinAnyRoot(lexical, roots) {
		return "", true, nil
	}

	canonical, err = filepath.EvalSymlinks(lexical)
	if err != nil {
		if os.IsNotExist(err) {
			return "", true, nil
		}
		return "", false, fmt.Errorf("resolve bootstrap PATH entry %q: %w", directory, err)
	}
	canonical = filepath.Clean(canonical)

	return canonical, isWithinAnyRoot(canonical, resolvedRoots), nil
}

/* ------------------------------------ Executable Resolution ----------------------------------- */

// resolveBootstrap resolves command from the canonical bootstrap PATH and rejects an executable
// whose physical target is inside the repository or Quill state. Resolving the final link closes
// the last provenance gap left by path lookup, which otherwise reports the symlink's pathname.
func resolveBootstrap(
	layout workspace.Layout,
	path string,
	command string,
) (executable string, err error) {
	resolved, err := process.ResolveExecutable(map[string]string{"PATH": path}, command)
	if err != nil {
		return "", fmt.Errorf("resolve bootstrap %s from trusted PATH: %w", command, err)
	}

	executable, err = filepath.EvalSymlinks(resolved)
	if err != nil {
		return "", fmt.Errorf("resolve bootstrap %s target: %w", command, err)
	}
	executable, err = filepath.Abs(executable)
	if err != nil {
		return "", fmt.Errorf("resolve bootstrap %s path: %w", command, err)
	}
	executable = filepath.Clean(executable)

	if isWithinAnyRoot(executable, resolveSymlinkRoots(bootstrapExclusions(layout))) {
		return "", fmt.Errorf(
			"bootstrap %s target %q is within repository or state directories",
			command,
			executable,
		)
	}

	return executable, nil
}

// isWithinAnyRoot reports whether candidate is equal to or beneath any of roots. Containment is
// component-aware through filepath.Rel, so a sibling directory whose name prefixes a root (for
// example "/repo-other" beside "/repo") is not mistaken for a child.
func isWithinAnyRoot(candidate string, roots []string) (within bool) {
	separator := string(os.PathSeparator)
	for _, root := range roots {
		rel, err := filepath.Rel(root, candidate)
		if err != nil {
			continue
		}

		if rel == "." {
			return true
		}

		if rel != ".." && !strings.HasPrefix(rel, ".."+separator) {
			return true
		}
	}

	return false
}
