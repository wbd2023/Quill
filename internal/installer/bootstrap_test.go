package installer

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/workspace"
)

/* ----------------------------------------- Test Helper ---------------------------------------- */

// bootstrapExecutableName returns the executable name suffix the process resolver expects on the
// current platform, so the hostile resolution tests are deterministic on POSIX and Windows.
func bootstrapExecutableName(name string) (executableName string) {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}

	return name
}

/* ---------------------------------- Bootstrap Path Filtering ---------------------------------- */

// TestBootstrapPathExcludesRepositoryAndStateEntries is the QUILL-TRUST-006 regression: the
// bootstrap PATH is the ambient host PATH with repository and state entries removed, preserving
// host order. A checked-out or cached go/npm can never be selected from the bootstrap PATH.
func TestBootstrapPathExcludesRepositoryAndStateEntries(t *testing.T) {
	repository := t.TempDir()
	layout := workspace.NewLayout(repository)
	stateBin := layout.BinaryDirectory()
	hostToolDirectory := t.TempDir()
	hostOtherDirectory := t.TempDir()

	for _, directory := range []string{stateBin, layout.CacheDirectory()} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", directory, err)
		}
	}

	pathValue := strings.Join(
		[]string{stateBin, layout.CacheDirectory(), hostToolDirectory, hostOtherDirectory},
		string(os.PathListSeparator),
	)
	t.Setenv("PATH", pathValue)

	filtered, err := bootstrapPath(layout)
	if err != nil {
		t.Fatalf("bootstrapPath: %v", err)
	}

	want := strings.Join(
		[]string{hostToolDirectory, hostOtherDirectory},
		string(os.PathListSeparator),
	)
	if filtered != want {
		t.Fatalf("bootstrap PATH = %q, want %q", filtered, want)
	}
}

// TestBootstrapPathFailsClosedWithoutTrustedHostDirectory pins the explicit non-empty invariant:
// a host whose every PATH entry is repository or state cannot bootstrap installation and fails
// closed rather than executing an untrusted candidate.
func TestBootstrapPathFailsClosedWithoutTrustedHostDirectory(t *testing.T) {
	repository := t.TempDir()
	layout := workspace.NewLayout(repository)
	stateBin := layout.BinaryDirectory()
	if err := os.MkdirAll(stateBin, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", stateBin, err)
	}

	t.Setenv("PATH", stateBin)

	if _, err := bootstrapPath(layout); err == nil {
		t.Fatal("bootstrapPath = nil, want error when no trusted host directory remains")
	}
}

// TestFilterBootstrapPathDropsSymlinkedHostEntry closes the symlink-resolution path of
// QUILL-TRUST-006: a host PATH entry that is a symlink into the repository or state directory is
// removed after EvalSymlinks, not just after a lexical match.
func TestFilterBootstrapPathDropsSymlinkedHostEntry(t *testing.T) {
	repository := t.TempDir()
	layout := workspace.NewLayout(repository)
	stateBin := layout.BinaryDirectory()
	if err := os.MkdirAll(stateBin, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", stateBin, err)
	}

	// A host directory that is actually a symlink back into the state directory.
	aliasedHost := filepath.Join(t.TempDir(), "tools")
	if err := os.Symlink(stateBin, aliasedHost); err != nil {
		t.Fatalf("create aliased host symlink: %v", err)
	}

	trustedHost := t.TempDir()
	pathValue := strings.Join(
		[]string{aliasedHost, trustedHost},
		string(os.PathListSeparator),
	)

	filtered, err := filterBootstrapPath(pathValue, bootstrapExclusions(layout))
	if err != nil {
		t.Fatalf("filterBootstrapPath: %v", err)
	}
	if filtered != trustedHost {
		t.Fatalf("filtered PATH = %q, want only trusted host %q", filtered, trustedHost)
	}
}

/* ---------------------------------- Hostile Resolution Proofs --------------------------------- */

// TestBootstrapResolutionSelectsHostExecutableOverCache is the combined QUILL-TRUST-006 hostile
// proof: with a malicious go (or npm) placed in the state cache and a valid go on the host, the
// bootstrap PATH excludes the cache entry and resolution selects the host executable. The cached
// executable can never run.
func TestBootstrapResolutionSelectsHostExecutableOverCache(t *testing.T) {
	tools := []string{"go", "npm"}

	for _, name := range tools {
		t.Run(name, func(t *testing.T) {
			repository := t.TempDir()
			layout := workspace.NewLayout(repository)

			stateBin := layout.BinaryDirectory()
			if err := os.MkdirAll(stateBin, 0o755); err != nil {
				t.Fatalf("mkdir state bin: %v", err)
			}
			hostBin := t.TempDir()

			// A hostile executable sitting in the repository state directory, plus a valid host
			// executable. Both are executable so resolution would pick whichever PATH lists
			// first.
			testutil.WriteExecutable(t, stateBin, bootstrapExecutableName(name), "cached-hostile")
			testutil.WriteExecutable(t, hostBin, bootstrapExecutableName(name), "trusted-host")

			// The pre-fix cache-first PATH would have selected the state executable.
			pathValue := strings.Join(
				[]string{stateBin, hostBin},
				string(os.PathListSeparator),
			)

			filtered, err := filterBootstrapPath(pathValue, bootstrapExclusions(layout))
			if err != nil {
				t.Fatalf("filterBootstrapPath: %v", err)
			}
			resolved, err := process.ResolveExecutable(
				map[string]string{"PATH": filtered},
				name,
			)
			if err != nil {
				t.Fatalf("resolve %s from bootstrap PATH: %v", name, err)
			}

			resolvedDir := filepath.Dir(filepath.Clean(resolved))
			if resolvedDir != hostBin {
				t.Fatalf("resolved %s = %q (dir %q), want host directory %q",
					name, resolved, resolvedDir, hostBin)
			}
			if strings.HasPrefix(filepath.Clean(resolved), stateBin) {
				t.Fatalf("resolved %s = %q, must not come from state directory %q",
					name, resolved, stateBin)
			}
		})
	}
}

// TestResolveBootstrapRejectsExecutableSymlinkIntoState closes the final bootstrap provenance
// path: a trusted host directory may not select a go/npm link whose target is in Quill state.
func TestResolveBootstrapRejectsExecutableSymlinkIntoState(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating the hostile executable symlink requires Windows link privilege")
	}

	layout := workspace.NewLayout(t.TempDir())
	stateBin := layout.BinaryDirectory()
	if err := os.MkdirAll(stateBin, 0o755); err != nil {
		t.Fatalf("mkdir state bin: %v", err)
	}

	name := bootstrapExecutableName("go")
	target := filepath.Join(stateBin, name)
	testutil.WriteExecutable(t, stateBin, name, "cached-hostile")

	hostBin := t.TempDir()
	if err := os.Symlink(target, filepath.Join(hostBin, name)); err != nil {
		t.Fatalf("create hostile host symlink: %v", err)
	}

	if _, err := resolveBootstrap(layout, hostBin, "go"); err == nil {
		t.Fatal("resolveBootstrap = nil, want state-target rejection")
	}
}

/* ---------------------------------- Cache Only Path Exclusion --------------------------------- */

// TestFilterBootstrapPathDropsStateOnlyEntries proves the cached-executable half of
// QUILL-TRUST-006 deterministically: when the ambient PATH contains only repository or state
// directories, every entry is removed from the bootstrap PATH. A cached go or npm sitting in the
// state directory therefore has no bootstrap entry to run from. (The non-empty invariant itself
// is enforced by bootstrapPath; see TestBootstrapPathFailsClosedWithoutTrustedHostDirectory.)
func TestFilterBootstrapPathDropsStateOnlyEntries(t *testing.T) {
	repository := t.TempDir()
	layout := workspace.NewLayout(repository)
	stateBin := layout.BinaryDirectory()
	if err := os.MkdirAll(stateBin, 0o755); err != nil {
		t.Fatalf("mkdir state bin: %v", err)
	}

	// A hostile executable sitting in the repository state directory would be selected first by
	// the cache-first PATH; the bootstrap PATH must drop its containing entry entirely.
	testutil.WriteExecutable(t, stateBin, bootstrapExecutableName("go"), "cached-hostile")

	filtered, err := filterBootstrapPath(stateBin, bootstrapExclusions(layout))
	if err != nil {
		t.Fatalf("filterBootstrapPath: %v", err)
	}
	if filtered != "" {
		t.Fatalf("bootstrap PATH = %q, want empty after excluding state-only entry %q",
			filtered, stateBin)
	}
}
