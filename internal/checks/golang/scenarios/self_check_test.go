package scenarios

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/quill/internal/checks/golang"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

func TestQuillPassesGoStyleChecks(t *testing.T) {
	root := testutil.RepositoryRoot(t)
	config := profiles.Self(t)

	result, err := golang.CheckDirectories(
		root,
		[]string{
			filepath.Join(root, "cmd"),
			filepath.Join(root, "internal"),
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
