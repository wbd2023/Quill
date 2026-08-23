package profile

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

/* ---------------------------- Physical Repository Path Containment ---------------------------- */

// validateRepositoryPaths enforces physical (symlink) containment for every repository-relative
// path value in config against root. Lexical containment is already applied by Validate. This pass
// resolves each value, or its deepest existing ancestor, and rejects values whose symlinks escape
// the repository. A missing target is accepted only when its resolved existing ancestor remains
// inside the repository.
//
// Profile cannot import the workspace resolver (architectural boundary), so the physical policy
// is mirrored here with standard-library primitives and kept in lock-step with
// workspace.ResolveRepoRelative. Errors name the offending Profile field and value.
func validateRepositoryPaths(config Profile, root string) (err error) {
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve repository root: %w", err)
	}

	for scope, roots := range config.Repository.ScopeRoots {
		for _, scopeRoot := range roots {
			if err = validateResolvedPath(canonical, scopeRoot,
				"repository.scope_roots."+string(scope)); err != nil {
				return err
			}
		}
	}

	for _, marker := range config.Repository.RootMarkers {
		if err = validateResolvedPath(canonical, marker, "repository.root_markers"); err != nil {
			return err
		}
	}

	if err = validateResolvedPath(canonical, config.StyleGuide.Path,
		"style_guide.path"); err != nil {
		return err
	}

	for _, target := range config.Targets {
		if target.WorkingDirectory != "" {
			field := fmt.Sprintf("target %q working_directory", target.Name)
			if err = validateResolvedPath(canonical, target.WorkingDirectory, field); err != nil {
				return err
			}
		}

		for _, path := range target.FormatPaths {
			if err = validateResolvedPath(canonical, path,
				fmt.Sprintf("target %q format_paths", target.Name)); err != nil {
				return err
			}
		}

		for _, path := range target.CheckPaths {
			if err = validateResolvedPath(canonical, path,
				fmt.Sprintf("target %q check_paths", target.Name)); err != nil {
				return err
			}
		}
	}

	for _, fileSet := range config.FileSets {
		for scope, files := range fileSet.Include.Files {
			for _, file := range files {
				field := fmt.Sprintf("file set %q include.files.%s", fileSet.Name, scope)
				if err = validateResolvedPath(canonical, file, field); err != nil {
					return err
				}
			}
		}

		for scope, paths := range fileSet.Include.Paths {
			for _, path := range paths {
				field := fmt.Sprintf("file set %q include.paths.%s", fileSet.Name, scope)
				if err = validateResolvedPath(canonical, path, field); err != nil {
					return err
				}
			}
		}
	}

	for _, source := range config.PackSources {
		if err = validateResolvedPath(canonical, source.Path, "pack_sources.path"); err != nil {
			return err
		}
	}

	return nil
}

// validateResolvedPath joins value under root and rejects a resolved path or its deepest existing
// ancestor when it escapes root. root must already be canonical.
func validateResolvedPath(root string, value string, field string) (err error) {
	joined := filepath.Join(root, filepath.FromSlash(value))
	resolved, err := resolveExistingAncestor(joined)
	if err != nil {
		return fmt.Errorf("%s: resolve path %q: %w", field, value, err)
	}

	if !isWithinRepoRoot(root, resolved) {
		return fmt.Errorf("%s: path %q resolves outside the repository root", field, value)
	}

	return nil
}

// resolveExistingAncestor resolves a path or, when it does not exist, its nearest physical
// ancestor. A dangling symlink is followed even when its target does not yet exist, so its
// resolved ancestor still participates in repository containment.
func resolveExistingAncestor(path string) (resolved string, err error) {
	for {
		resolved, err = filepath.EvalSymlinks(path)
		if err == nil {
			return resolved, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}

		info, statErr := os.Lstat(path)
		if statErr == nil && info.Mode()&fs.ModeSymlink != 0 {
			target, readErr := os.Readlink(path)
			if readErr != nil {
				return "", readErr
			}
			if !filepath.IsAbs(target) {
				parent, parentErr := filepath.EvalSymlinks(filepath.Dir(path))
				if parentErr != nil {
					return "", parentErr
				}
				target = filepath.Join(parent, target)
			}
			path = filepath.Clean(target)
			continue
		}
		if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
			return "", statErr
		}

		parent := filepath.Dir(path)
		if parent == path {
			return "", err
		}
		path = parent
	}
}

func isWithinRepoRoot(root string, target string) (within bool) {
	rel, err := filepath.Rel(root, filepath.Clean(target))
	if err != nil {
		return false
	}

	rel = filepath.ToSlash(rel)
	return rel != ".." && !strings.HasPrefix(rel, "../")
}
