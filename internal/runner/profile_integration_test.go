package runner

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

func TestNewContextLoadsCurrentProfileFixture(t *testing.T) {
	fixtureRoot := t.TempDir()
	profiles.Write(t, fixtureRoot, profiles.Current(t))

	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "core", "domain", "errors.go"),
		"package domain\n\nvar ErrMissing = error(nil)\n",
	)
	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "client", "application", "port", "messages", "repository.go"),
		"package messages\n\ntype MessageRepository interface { ListMessages() }\n",
	)

	context := testContext(t, fixtureRoot, style.Scope("all"))

	if len(context.Effective.Rules) == 0 {
		t.Fatal("expected effective rules to be loaded for fixture")
	}
}
