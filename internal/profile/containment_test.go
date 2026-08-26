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

func TestLoadRejectsProfileSymlinkEscapingRepository(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	config := profiles.Self(t)
	config.Repository.RootMarkers = []string{"STYLE.md"}
	profiles.Write(t, root, config)

	outside := t.TempDir()
	profiles.Write(t, outside, config)

	path := filepath.Join(root, profile.DefaultFilename)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(
		filepath.Join(outside, profile.DefaultFilename),
		path,
	); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	_, err := profile.Load(root)
	requireErrorContains(t, err, profile.DefaultFilename)
}

func TestLoadRejectsSymlinkEscapingConfiguredPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		config    func(*profile.Profile)
		stage     func(t *testing.T, root string, outside string)
		wantField string
	}{
		{
			name: "scope root",
			config: func(config *profile.Profile) {
				config.Repository.ScopeRoots = map[style.Scope][]string{
					style.Scope("all"): {"escape"},
				}
				config.Repository.DefaultScope = style.Scope("all")
			},
			stage: func(t *testing.T, root string, outside string) {
				if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
					t.Skipf("symlink unsupported: %v", err)
				}
			},
			wantField: "scope_roots.all",
		},
		{
			name: "target check path",
			config: func(config *profile.Profile) {
				if len(config.Targets) == 0 {
					t.Fatal("self profile must define at least one target")
				}
				config.Targets[0].CheckPaths = []string{"escape"}
			},
			stage: func(t *testing.T, root string, outside string) {
				if err := os.MkdirAll(filepath.Join(outside, "sneaked"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(
					filepath.Join(outside, "sneaked"), filepath.Join(root, "escape"),
				); err != nil {
					t.Skipf("symlink unsupported: %v", err)
				}
			},
			wantField: "check_paths",
		},
		{
			name: "style guide path",
			config: func(config *profile.Profile) {
				config.StyleGuide.Path = "alias.md"
			},
			stage: func(t *testing.T, root string, outside string) {
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
			},
			wantField: "style_guide.path",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			config := profiles.Self(t)
			test.config(&config)
			profiles.Write(t, root, config)
			test.stage(t, root, t.TempDir())

			_, err := profile.Load(root)
			requireErrorContains(t, err, test.wantField)
			requireErrorContains(t, err, "outside the repository root")
		})
	}
}

func TestLoadRejectsMissingPathBeneathEscapingSymlink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config func(*profile.Profile)
	}{
		{
			name: "scope root",
			config: func(config *profile.Profile) {
				config.Repository.ScopeRoots = map[style.Scope][]string{
					style.Scope("all"): {"escape/missing"},
				}
				config.Repository.DefaultScope = style.Scope("all")
			},
		},
		{
			name: "root marker",
			config: func(config *profile.Profile) {
				config.Repository.RootMarkers = []string{"escape/missing"}
			},
		},
		{
			name: "style guide",
			config: func(config *profile.Profile) {
				config.StyleGuide.Path = "escape/missing"
			},
		},
		{
			name: "target working directory",
			config: func(config *profile.Profile) {
				config.Targets[0].WorkingDirectory = "escape/missing"
			},
		},
		{
			name: "target format path",
			config: func(config *profile.Profile) {
				config.Targets[0].FormatPaths = []string{"escape/missing"}
			},
		},
		{
			name: "target check path",
			config: func(config *profile.Profile) {
				config.Targets[0].CheckPaths = []string{"escape/missing"}
			},
		},
		{
			name: "file set include file",
			config: func(config *profile.Profile) {
				config.FileSets[0].Include.Files = map[style.Scope][]string{
					config.Repository.DefaultScope: {"escape/missing"},
				}
			},
		},
		{
			name: "file set include path",
			config: func(config *profile.Profile) {
				config.FileSets[0].Include.Paths = map[style.Scope][]string{
					config.Repository.DefaultScope: {"escape/missing"},
				}
			},
		},
		{
			name: "Pack source",
			config: func(config *profile.Profile) {
				config.PackSources = []profile.PackSource{{Path: "escape/missing"}}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			config := profiles.Self(t)
			profiles.Write(t, root, config)

			if err := os.Symlink(
				filepath.Join(t.TempDir(), "missing"),
				filepath.Join(root, "escape"),
			); err != nil {
				t.Skipf("symlink unsupported: %v", err)
			}

			test.config(&config)
			contents := profiles.Format(t, config)
			if err := os.WriteFile(
				filepath.Join(root, profile.DefaultFilename),
				[]byte(contents),
				0o600,
			); err != nil {
				t.Fatal(err)
			}

			_, err := profile.Load(root)
			requireErrorContains(t, err, "outside the repository root")
		})
	}
}

func TestLoadRejectsRelativeDanglingSymlinkBeneathAlias(t *testing.T) {
	t.Parallel()

	parent := t.TempDir()
	root := filepath.Join(parent, "root")

	config := profiles.Self(t)
	profiles.Write(t, root, config)
	real := filepath.Join(root, "real")
	alias := filepath.Join(root, "deep", "nested", "alias")
	if err := os.MkdirAll(filepath.Dir(alias), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(real, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(parent, "outside"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../../real", alias); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	if err := os.Symlink("../../outside/missing", filepath.Join(real, "escape")); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	config.StyleGuide.Path = "deep/nested/alias/escape/future"
	contents := profiles.Format(t, config)
	if err := os.WriteFile(
		filepath.Join(root, profile.DefaultFilename),
		[]byte(contents),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	_, err := profile.Load(root)
	requireErrorContains(t, err, "outside the repository root")
}

func TestLoadAcceptsMissingPathBeneathContainedAncestor(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	config := profiles.Self(t)
	profiles.Write(t, root, config)

	config.Targets[0].CheckPaths = []string{"generated/not-yet-created"}
	contents := profiles.Format(t, config)
	if err := os.WriteFile(
		filepath.Join(root, profile.DefaultFilename),
		[]byte(contents),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	if _, err := profile.Load(root); err != nil {
		t.Fatalf("Load rejected missing contained path: %v", err)
	}
}
