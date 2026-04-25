package runtime

import (
	"os"
	"path/filepath"
	"strings"
)

// Layout describes the repository-local directories used by the style toolchain.
type Layout struct {
	RepoRoot     string
	ToolsDir     string
	StateDir     string
	CacheDir     string
	ToolBinDir   string
	NodeDir      string
	NodeBinDir   string
	NpmCache     string
	GoBuildCache string
	GoModCache   string
	GoPath       string
}

// LayoutForRepository builds the style-tool layout from the repository root.
func LayoutForRepository(repoRoot string) (layout Layout) {
	return newLayout(repoRoot, filepath.Join(repoRoot, "tools"))
}

// LayoutForToolsDir builds the style-tool layout from the tools directory path.
func LayoutForToolsDir(toolsDir string) (layout Layout) {
	repoRoot := filepath.Clean(filepath.Join(toolsDir, ".."))
	return newLayout(repoRoot, toolsDir)
}

func newLayout(repoRoot string, toolsDir string) (layout Layout) {
	stateDir := filepath.Join(repoRoot, ".cache", "style")
	cacheDir := filepath.Join(stateDir, "cache")
	nodeDir := filepath.Join(stateDir, "npm")

	return Layout{
		RepoRoot:     repoRoot,
		ToolsDir:     toolsDir,
		StateDir:     stateDir,
		CacheDir:     cacheDir,
		ToolBinDir:   filepath.Join(stateDir, "bin"),
		NodeDir:      nodeDir,
		NodeBinDir:   filepath.Join(nodeDir, "node_modules", ".bin"),
		NpmCache:     filepath.Join(cacheDir, "npm"),
		GoBuildCache: filepath.Join(cacheDir, "go-build"),
		GoModCache:   filepath.Join(cacheDir, "go-mod"),
		GoPath:       filepath.Join(cacheDir, "gopath"),
	}
}

// SearchPath returns the PATH value that exposes pinned repository-local tools first.
func (layout Layout) SearchPath() (value string) {
	return strings.Join(
		[]string{layout.ToolBinDir, layout.NodeBinDir, os.Getenv("PATH")},
		string(os.PathListSeparator),
	)
}

// ToolEnvironment returns the base environment for external style tools.
func (layout Layout) ToolEnvironment() (environment map[string]string) {
	return map[string]string{
		"PATH": layout.SearchPath(),
	}
}

// GoEnvironment returns the tool environment with repository-local Go caches.
func (layout Layout) GoEnvironment() (environment map[string]string) {
	environment = layout.ToolEnvironment()
	environment["GOCACHE"] = layout.GoBuildCache
	environment["GOMODCACHE"] = layout.GoModCache
	environment["GOPATH"] = layout.GoPath
	return environment
}
