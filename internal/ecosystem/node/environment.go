package node

import (
	"path/filepath"

	"github.com/wbd2023/Quill/internal/workspace"
)

// InstallDirectory returns the directory where npm install operates.
func InstallDirectory(layout workspace.Layout) (directory string) {
	return filepath.Join(layout.StateDirectory, "npm")
}

// BinaryDirectory returns the directory where Node-installed binaries live.
func BinaryDirectory(layout workspace.Layout) (directory string) {
	return filepath.Join(InstallDirectory(layout), "node_modules", ".bin")
}

// CacheDirectory returns the npm cache directory.
func CacheDirectory(layout workspace.Layout) (cache string) {
	return filepath.Join(layout.CacheDirectory(), "npm")
}

// Environment returns the environment variables for executing Node tooling with isolated caches.
// npm_config_cache is set to a layout-derived path; PATH is the path argument.
func Environment(layout workspace.Layout, path string) (environment map[string]string) {
	return map[string]string{
		"PATH":             path,
		"npm_config_cache": CacheDirectory(layout),
	}
}
