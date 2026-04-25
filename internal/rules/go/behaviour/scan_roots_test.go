package behaviour

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	gostyle "ciphera/tools/internal/rules/go"
)

func TestStylecheckRejectsMissingScanRoot(t *testing.T) {
	missingRoot := filepath.Join(t.TempDir(), "missing")

	if _, err := gostyle.CheckDirectories(
		fixtures.RepoRoot(t),
		[]string{missingRoot},
		profiles.Current(t),
	); err == nil {
		t.Fatal("expected missing scan root error")
	}
}
