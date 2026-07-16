package runtime

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Layout represents the repository root and state directory from which the engine derives its
// filesystem paths.
type Layout struct {
	RepositoryRoot string
	StateDirectory string
}

// NewLayout creates a Layout for root, defaulting the state directory to .cache/quill/ under the
// root.
func NewLayout(root string) (layout Layout) {
	return Layout{
		RepositoryRoot: root,
		StateDirectory: filepath.Join(root, ".cache", "quill"),
	}
}

// CacheDirectory returns the cache directory under the state directory.
func (layout Layout) CacheDirectory() (directory string) {
	return filepath.Join(layout.StateDirectory, "cache")
}

// BinaryDirectory returns the directory where installed tool binaries live.
func (layout Layout) BinaryDirectory() (directory string) {
	return filepath.Join(layout.StateDirectory, "bin")
}

// BuildPath produces a PATH environment variable value from the layout's binary directory, the
// given directories, and the system PATH. The layout's binaries take priority over system
// equivalents.
func (layout Layout) BuildPath(directories ...string) (path string) {
	combined := slices.Concat(
		[]string{layout.BinaryDirectory()},
		directories,
		[]string{os.Getenv("PATH")},
	)
	return strings.Join(combined, string(os.PathListSeparator))
}
