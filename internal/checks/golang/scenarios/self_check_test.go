package scenarios

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/Quill/internal/checks/golang"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func TestQuillPassesGoStyleChecks(t *testing.T) {
	repositoryRoot := testutil.RepositoryRoot(t)
	config := profiles.Current(t)

	result, err := golang.CheckDirectories(
		repositoryRoot,
		[]string{
			filepath.Join(repositoryRoot, "cmd"),
			filepath.Join(repositoryRoot, "internal"),
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
