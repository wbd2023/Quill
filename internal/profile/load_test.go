package profile_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/profile"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func TestLoadReadsCurrentProfile(t *testing.T) {
	t.Parallel()

	config, err := profile.Load(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if config.SchemaVersion != 1 {
		t.Fatalf("schema version = %d", config.SchemaVersion)
	}
}

func TestLoadWrapsProfilePathOnParseError(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	testutil.WriteFile(t, root, profile.DefaultFilename, "schema_version = 2\n")

	_, err := profile.Load(root)
	requireErrorContains(t, err, profile.DefaultFilename)
}

func TestLoadRejectsMissingConfiguredRootMarker(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	config := profiles.Current(t)
	config.Repository.RootMarkers = []string{"PROJECT.marker"}
	profiles.Write(t, root, config)

	_, err := profile.Load(root)
	requireErrorContains(t, err, "missing marker")
}
