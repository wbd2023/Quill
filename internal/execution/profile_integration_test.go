package execution

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

func TestNewRunContextLoadsCurrentProfileFixture(t *testing.T) {
	root := t.TempDir()
	profiles.Write(t, root, profiles.Self(t))

	testutil.WriteFile(
		t,
		root,
		filepath.Join("internal", "core", "domain", "errors.go"),
		"package domain\n\nvar ErrMissing = error(nil)\n",
	)
	testutil.WriteFile(
		t,
		root,
		filepath.Join("internal", "client", "application", "port", "messages", "repository.go"),
		"package messages\n\ntype MessageRepository interface { ListMessages() }\n",
	)

	context := testContext(t, root, style.Scope("all"))

	if len(context.Plan.Rules) == 0 {
		t.Fatal("expected compiled rules to be loaded for fixture")
	}
}
