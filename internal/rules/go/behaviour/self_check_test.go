package behaviour

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	gostyle "ciphera/tools/internal/rules/go"
)

func TestStylePlatformPassesGoStyleChecks(t *testing.T) {
	toolsRoot := fixtures.ToolsRoot(t)
	output, err := gostyle.CheckDirectories(
		filepath.Clean(filepath.Join(toolsRoot, "..")),
		[]string{
			filepath.Join(toolsRoot, "cmd"),
			filepath.Join(toolsRoot, "internal"),
		},
		profiles.Current(t),
	)
	if err != nil {
		t.Fatalf("expected style platform to satisfy Go style checks, output:\n%s", output)
	}
}
