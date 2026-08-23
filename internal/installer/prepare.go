package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wbd2023/quill/internal/ecosystem/golang"
	"github.com/wbd2023/quill/internal/ecosystem/node"
	"github.com/wbd2023/quill/internal/workspace"
)

// prepareInstallDirectory ensures dir exists as a real directory beneath root, creating missing
// components and rejecting any component that is a symlink or non-directory entry. A destination
// that escapes the repository root (for example a custom state directory configured outside the
// repository) is rejected before any third-party tooling writes into it.
//
// Each existing component is inspected with Lstat, never followed, so a symlinked cache or state
// component cannot redirect the subsequent Go or npm writes outside the repository. The check is
// repeated for every component on the path from the root, closing both a pre-existing link and a
// link planted on a freshly created parent.
func prepareInstallDirectory(root string, dir string) (err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve installation root %q: %w", root, err)
	}
	absRoot = filepath.Clean(absRoot)

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("resolve installation directory %q: %w", dir, err)
	}
	absDir = filepath.Clean(absDir)

	rel, err := filepath.Rel(absRoot, absDir)
	if err != nil {
		return fmt.Errorf("resolve installation directory beneath root: %w", err)
	}

	if rel == "." {
		return nil
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("installation directory %q escapes root %q", absDir, absRoot)
	}

	current := absRoot
	for _, component := range strings.Split(rel, string(os.PathSeparator)) {
		if component == "" || component == "." {
			continue
		}

		current = filepath.Join(current, component)
		info, statErr := os.Lstat(current)
		if os.IsNotExist(statErr) {
			if mkdirErr := os.Mkdir(current, standardPermissions); mkdirErr != nil && !os.IsExist(mkdirErr) {
				return fmt.Errorf("create installation directory %q: %w", current, mkdirErr)
			}
			info, statErr = os.Lstat(current)
		}
		if statErr != nil {
			return fmt.Errorf("inspect installation directory %q: %w", current, statErr)
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("installation directory %q is a symlink", current)
		}

		if !info.IsDir() {
			return fmt.Errorf("installation path %q is not a directory", current)
		}
	}

	return nil
}

// prepareGoInstall prepares every directory the Go toolchain writes to during bootstrap,
// rejecting symlinked or out-of-repository components before go install runs. The directory list
// is sourced from the golang package so the installer stays the single owner of preparation
// policy while the ecosystem package remains the source of its own layout.
func prepareGoInstall(layout workspace.Layout) (err error) {
	return prepareInstallDirectories(layout.RepositoryRoot, []string{
		layout.StateDirectory,
		layout.BinaryDirectory(),
		layout.CacheDirectory(),
		golang.BuildCacheDirectory(layout),
		golang.ModuleCacheDirectory(layout),
		golang.GoPath(layout),
	})
}

// prepareNpmInstall prepares every directory and manifest leaf npm uses during bootstrap,
// rejecting symlinked or out-of-repository components before npm install runs.
func prepareNpmInstall(layout workspace.Layout) (err error) {
	if err = prepareInstallDirectories(layout.RepositoryRoot, []string{
		layout.StateDirectory,
		layout.CacheDirectory(),
		node.CacheDirectory(layout),
		node.InstallDirectory(layout),
		node.BinaryDirectory(layout),
	}); err != nil {
		return err
	}

	return prepareNpmManifests(node.InstallDirectory(layout))
}

func prepareInstallDirectories(root string, directories []string) (err error) {
	for _, directory := range directories {
		if err = prepareInstallDirectory(root, directory); err != nil {
			return err
		}
	}

	return nil
}

// prepareNpmManifests rejects manifest leaves npm can read or rewrite. npm runs with --no-save and
// --package-lock=false, but an existing package file still affects its install plan; following one
// through a symlink would cross the repository-local state boundary.
func prepareNpmManifests(directory string) error {
	for _, name := range []string{"package.json", "package-lock.json", "npm-shrinkwrap.json"} {
		path := filepath.Join(directory, name)
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}

		if err != nil {
			return fmt.Errorf("inspect npm manifest %q: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("npm manifest %q is not a regular file", path)
		}
	}

	return nil
}
