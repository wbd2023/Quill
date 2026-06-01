package scenarios

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/rules/golang"
)

func TestStylePlatformPassesGoStyleChecks(t *testing.T) {
	toolsRoot := fixtures.ToolsRoot(t)
	config := profiles.Current(t)

	result, err := golang.CheckDirectories(
		filepath.Clean(filepath.Join(toolsRoot, "..")),
		[]string{
			filepath.Join(toolsRoot, "cmd"),
			filepath.Join(toolsRoot, "internal"),
		},
		config.Repository,
		config.PathRoles,
		goConfigForTest(t, config),
	)
	if err != nil {
		t.Fatalf(
			"expected style platform to satisfy Go style checks, diagnostics: %#v",
			result.Diagnostics,
		)
	}
}
