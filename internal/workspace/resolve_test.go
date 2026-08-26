package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------- Validation Tests -------------------------------------- */

func TestValidateRepoRelativeAcceptsContainedPaths(t *testing.T) {
	t.Parallel()

	values := []string{
		".",
		"",
		"cmd",
		"internal/",
		"STYLE.md",
		"a/b/c.go",
		"./cmd",
		"cmd/.",
		"foo/../bar",
	}
	for _, value := range values {
		if err := ValidateRepoRelative(value); err != nil {
			t.Fatalf("ValidateRepoRelative(%q) = %v, want nil", value, err)
		}
	}
}

func TestValidateRepoRelativeRejectsEscapingPaths(t *testing.T) {
	t.Parallel()

	values := []string{
		"..",
		"../etc",
		"a/../../b",
		"foo/../../bar",
		"/etc",
		"/",
		"C:\\windows",
		"C:/users",
		"c:etc",
		"\\windows",
		"a\x00b",
		"\x00",
	}
	for _, value := range values {
		if err := ValidateRepoRelative(value); err == nil {
			t.Fatalf("ValidateRepoRelative(%q) = nil, want rejection", value)
		}
	}
}

// TestValidateRepoRelativeRejectsCrossPlatformTraversal guards the escape check against the
// Windows behaviour of filepath.Clean, which re-introduces backslash separators after slash
// normalisation. ValidateRepoRelative applies filepath.ToSlash after filepath.Clean so that
// parent-traversal escapes are rejected on every GOOS; the backslash variants below exercise the
// full normalisation pipeline and must be rejected everywhere.
func TestValidateRepoRelativeRejectsCrossPlatformTraversal(t *testing.T) {
	t.Parallel()

	values := []string{
		"../etc",
		"..\\etc",
		"a/../../b",
		"a\\..\\..\\b",
		"foo\\..\\..\\bar",
		"..\\..",
	}
	for _, value := range values {
		if err := ValidateRepoRelative(value); err == nil {
			t.Fatalf("ValidateRepoRelative(%q) = nil, want rejection on every platform", value)
		}
	}
}

func TestValidateRepoRelativeErrorNamesValue(t *testing.T) {
	t.Parallel()

	if err := ValidateRepoRelative("../etc"); err == nil ||
		!strings.Contains(err.Error(), "../etc") {
		t.Fatalf("error must name the offending value %q, got %v", "../etc", err)
	}
}

/* -------------------------------------- Resolution Tests -------------------------------------- */

func TestCanonicalRootResolvesSymlinkedRoot(t *testing.T) {
	t.Parallel()

	target := t.TempDir()
	linkParent := t.TempDir()
	link := filepath.Join(linkParent, "repo-link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	got, err := CanonicalRoot(link)
	if err != nil {
		t.Fatalf("CanonicalRoot: %v", err)
	}

	want, err := CanonicalRoot(target)
	if err != nil {
		t.Fatalf("CanonicalRoot target: %v", err)
	}
	if got != want {
		t.Fatalf("CanonicalRoot(%q) = %q, want %q", link, got, want)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("canonical root %q must be absolute", got)
	}
}

func TestResolveRepoRelativeJoinsContainedExistingPath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "internal", "profile"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveRepoRelative(root, "internal/profile")
	if err != nil {
		t.Fatalf("ResolveRepoRelative: %v", err)
	}

	want := filepath.Join(root, "internal", "profile")
	if got != want {
		t.Fatalf("ResolveRepoRelative = %q, want %q", got, want)
	}
}

func TestResolveRepoRelativeAcceptsNonExistentTarget(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	got, err := ResolveRepoRelative(root, "future/dir")
	if err != nil {
		t.Fatalf("ResolveRepoRelative non-existent: %v", err)
	}

	want := filepath.Join(root, "future", "dir")
	if got != want {
		t.Fatalf("ResolveRepoRelative = %q, want %q", got, want)
	}
}

func TestResolveRepoRelativeRejectsUnsafeTargets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		prepare func(t *testing.T, root string)
	}{
		{
			name: "missing path beneath escaping symlink",
			path: "escape/future",
			prepare: func(t *testing.T, root string) {
				outside := t.TempDir()
				if err := os.Symlink(
					filepath.Join(outside, "missing"), filepath.Join(root, "escape"),
				); err != nil {
					t.Skipf("symlink unsupported: %v", err)
				}
			},
		},
		{
			name: "relative dangling symlink beneath alias",
			path: "deep/nested/alias/escape/future",
			prepare: func(t *testing.T, root string) {
				real := filepath.Join(root, "real")
				alias := filepath.Join(root, "deep", "nested", "alias")
				if err := os.MkdirAll(filepath.Dir(alias), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.MkdirAll(real, 0o755); err != nil {
					t.Fatal(err)
				}
				parentOutside := filepath.Join(filepath.Dir(root), "outside")
				if err := os.MkdirAll(parentOutside, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("../../real", alias); err != nil {
					t.Skipf("symlink unsupported: %v", err)
				}
				missingTarget := filepath.Join(real, "escape")
				if err := os.Symlink("../../outside/missing", missingTarget); err != nil {
					t.Skipf("symlink unsupported: %v", err)
				}
			},
		},
		{
			name: "parent escape",
			path: "../etc",
		},
		{
			name: "absolute path",
			path: "/etc",
		},
		{
			name: "symlink escape",
			path: "escape",
			prepare: func(t *testing.T, root string) {
				outside := t.TempDir()
				secret := filepath.Join(outside, "secret")
				if err := os.WriteFile(secret, []byte("private"), 0o600); err != nil {
					t.Fatal(err)
				}

				if err := os.Symlink(secret, filepath.Join(root, "escape")); err != nil {
					t.Skipf("symlink unsupported: %v", err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			if test.prepare != nil {
				test.prepare(t, root)
			}

			if _, err := ResolveRepoRelative(root, test.path); err == nil {
				t.Fatalf("ResolveRepoRelative accepted unsafe target %q", test.path)
			}
		})
	}
}

func TestResolveRepoRelativeAcceptsInRepoSymlink(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	realDir := filepath.Join(root, "real")
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "alias")
	if err := os.Symlink(realDir, link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	if _, err := ResolveRepoRelative(root, "alias"); err != nil {
		t.Fatalf("ResolveRepoRelative rejected an in-repo symlink: %v", err)
	}
}
