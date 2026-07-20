package execution

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func TestNewRunContextLoadsCurrentProfileFixture(t *testing.T) {
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
