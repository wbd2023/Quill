package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

/* ------------------------------------- Lexical containment ------------------------------------ */

func TestValidateRejectsAbsoluteScopeRoot(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Repository.ScopeRoots[profiletest.Scope] = []string{"/etc"}

	requireErrorContains(t, profile.Validate(config), "repository.scope_roots")
}

func TestValidateRejectsTraversalScopeRoot(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Repository.ScopeRoots[profiletest.Scope] = []string{"../etc"}

	requireErrorContains(t, profile.Validate(config), "escapes the repository root")
}

// TestValidateRejectsCrossPlatformTraversal guards validateRepoPath against the Windows
// behaviour of filepath.Clean, which re-introduces backslash separators after normalisation.
// The backslash variants must be rejected on every GOOS.
func TestValidateRejectsCrossPlatformTraversal(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"../etc", "..\\etc", "a/../../b", "a\\..\\..\\b"} {
		config := profiletest.Config()
		config.Repository.ScopeRoots[profiletest.Scope] = []string{value}
		if err := profile.Validate(config); err == nil {
			t.Fatalf("Validate accepted escaping scope root %q", value)
		}
	}
}

func TestValidateRejectsNULByteInScopeRoot(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Repository.ScopeRoots[profiletest.Scope] = []string{"a\x00b"}

	requireErrorContains(t, profile.Validate(config), "NUL")
}

func TestValidateRejectsAbsoluteRootMarker(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Repository.RootMarkers = append(config.Repository.RootMarkers, "/etc/passwd")

	requireErrorContains(t, profile.Validate(config), "repository.root_markers")
}

func TestValidateRejectsAbsoluteTargetWorkingDirectory(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Targets[0].WorkingDirectory = "/etc"

	requireErrorContains(t, profile.Validate(config), "working_directory")
}

func TestValidateRejectsTraversalTargetCheckPath(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Targets[0].CheckPaths = []string{"../etc"}

	err := profile.Validate(config)
	requireErrorContains(t, err, "check_paths")
	requireErrorContains(t, err, "escapes the repository root")
}

func TestValidateRejectsTraversalTargetFormatPath(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Targets[0].FormatPaths = []string{"../../etc"}

	requireErrorContains(t, profile.Validate(config), "format_paths")
}

func TestValidateRejectsAbsoluteFileSetIncludeFile(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.FileSets[0].Include.Files = map[style.Scope][]string{
		profiletest.Scope: {"/etc/passwd"},
	}

	requireErrorContains(t, profile.Validate(config), "include.files")
}

func TestValidateRejectsAbsoluteStyleGuidePath(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.StyleGuide.Path = "/etc/STYLE.md"

	requireErrorContains(t, profile.Validate(config), "style_guide.path")
}

func TestValidateAcceptsRepositoryRelativePaths(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Repository.ScopeRoots[profiletest.Scope] = []string{".", "internal"}
	config.Targets[0].WorkingDirectory = "cmd"
	config.Targets[0].FormatPaths = []string{"cmd", "internal"}
	config.StyleGuide.Path = "docs/STYLE.md"

	if err := profile.Validate(config); err != nil {
		t.Fatalf("Validate rejected contained paths: %v", err)
	}
}

/* ------------------------------------ Physical containment ------------------------------------ */

func TestLoadAcceptsContainedRepository(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	profiles.Write(t, root, profiles.Self(t))

	if _, err := profile.Load(root); err != nil {
		t.Fatalf("Load rejected a contained repository: %v", err)
	}
}

func TestLoadRejectsSymlinkEscapingScopeRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	config := profiles.Self(t)
	config.Repository.ScopeRoots = map[style.Scope][]string{
		style.Scope("all"): {"escape"},
	}
	config.Repository.DefaultScope = style.Scope("all")
	profiles.Write(t, root, config)

	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	_, err := profile.Load(root)
	requireErrorContains(t, err, "scope_roots.all")
	requireErrorContains(t, err, "outside the repository root")
}

func TestLoadRejectsSymlinkEscapingTargetCheckPath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	config := profiles.Self(t)
	if len(config.Targets) == 0 {
		t.Fatal("self profile must define at least one target")
	}
	config.Targets[0].CheckPaths = []string{"escape"}
	profiles.Write(t, root, config)

	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outside, "sneaked"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(
		filepath.Join(outside, "sneaked"), filepath.Join(root, "escape"),
	); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	_, err := profile.Load(root)
	requireErrorContains(t, err, "check_paths")
	requireErrorContains(t, err, "outside the repository root")
}

func TestLoadRejectsSymlinkEscapingStyleGuidePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	config := profiles.Self(t)
	config.StyleGuide.Path = "alias.md"
	profiles.Write(t, root, config)

	outside := t.TempDir()
	target := filepath.Join(outside, "foreign.md")
	if err := os.WriteFile(target, []byte("# foreign\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "alias.md")
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	_, err := profile.Load(root)
	requireErrorContains(t, err, "style_guide.path")
	requireErrorContains(t, err, "outside the repository root")
}
