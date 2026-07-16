package runtime

import (
	"os"
	"path/filepath"
	"strings"
)

// Option configures a Layout at construction time.
type Option func(*Layout)

// WithStateDirectory overrides the default state directory (.cache/quill/ under the repository
// root). Use this when the engine's state (downloaded tools, caches) must live elsewhere.
func WithStateDirectory(directory string) (option Option) {
	return func(layout *Layout) {
		layout.stateDirectory = directory
	}
}

// Layout holds the repository root and state directory from which all engine paths are derived.
// Construct one with NewLayout.
type Layout struct {
	repositoryRoot string
	stateDirectory string
}

// NewLayout constructs a Layout for the given repository root, with the engine state directory
// defaulting to .cache/quill/ under the root. Pass options to override defaults.
func NewLayout(repositoryRoot string, options ...Option) (layout Layout) {
	layout = Layout{
		repositoryRoot: repositoryRoot,
		stateDirectory: filepath.Join(repositoryRoot, ".cache", "quill"),
	}
	for _, option := range options {
		option(&layout)
	}
	return layout
}

// RepositoryRoot returns the repository root the engine is operating on.
func (layout Layout) RepositoryRoot() (root string) {
	return layout.repositoryRoot
}

// StateDirectory returns the directory where the engine stores its state (downloaded tools and
// caches).
func (layout Layout) StateDirectory() (directory string) {
	return layout.stateDirectory
}

// CacheDirectory returns the cache directory under the state directory.
func (layout Layout) CacheDirectory() (directory string) {
	return filepath.Join(layout.stateDirectory, "cache")
}

// ToolBinaryDirectory returns the directory where installed tool binaries live.
func (layout Layout) ToolBinaryDirectory() (directory string) {
	return filepath.Join(layout.stateDirectory, "bin")
}

// SearchPath joins the given binary directories and the system PATH, producing a PATH value that
// makes installed tool binaries discoverable. The caller provides the ecosystem-specific binary
// directories (eg the engine's tool binary directory and a Node binary directory).
func SearchPath(binaryDirectories ...string) (value string) {
	directories := append([]string{}, binaryDirectories...)
	directories = append(directories, os.Getenv("PATH"))
	return strings.Join(directories, string(os.PathListSeparator))
}
