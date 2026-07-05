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

/* ------------------------------------------- Layout ------------------------------------------- */

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

/* ---------------------------------------- Derived Paths --------------------------------------- */

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

// NodeDirectory returns the directory where Node tooling state lives.
func (layout Layout) NodeDirectory() (directory string) {
	return filepath.Join(layout.stateDirectory, "npm")
}

// NodeBinaryDirectory returns the directory where Node-installed binaries live.
func (layout Layout) NodeBinaryDirectory() (directory string) {
	return filepath.Join(layout.NodeDirectory(), "node_modules", ".bin")
}

// NpmCache returns the NPM cache directory.
func (layout Layout) NpmCache() (cache string) {
	return filepath.Join(layout.CacheDirectory(), "npm")
}

// GoBuildCache returns the Go build cache directory.
func (layout Layout) GoBuildCache() (cache string) {
	return filepath.Join(layout.CacheDirectory(), "go-build")
}

// GoModuleCache returns the Go module cache directory.
func (layout Layout) GoModuleCache() (cache string) {
	return filepath.Join(layout.CacheDirectory(), "go-mod")
}

// GoPath returns the GOPATH directory.
func (layout Layout) GoPath() (path string) {
	return filepath.Join(layout.CacheDirectory(), "gopath")
}

/* ----------------------------------------- Environment ---------------------------------------- */

// SearchPath returns the PATH value that makes installed tool binaries and Node binaries
// discoverable.
func (layout Layout) SearchPath() (value string) {
	return strings.Join(
		[]string{layout.ToolBinaryDirectory(), layout.NodeBinaryDirectory(), os.Getenv("PATH")},
		string(os.PathListSeparator),
	)
}

// ToolEnvironment returns the environment variables for running non-Go tools (PATH only).
func (layout Layout) ToolEnvironment() (environment map[string]string) {
	return map[string]string{
		"PATH": layout.SearchPath(),
	}
}

// GoEnvironment returns the environment variables for running Go tools (PATH plus Go-specific
// cache and module paths).
func (layout Layout) GoEnvironment() (environment map[string]string) {
	environment = layout.ToolEnvironment()
	environment["GOCACHE"] = layout.GoBuildCache()
	environment["GOMODCACHE"] = layout.GoModuleCache()
	environment["GOPATH"] = layout.GoPath()
	return environment
}
