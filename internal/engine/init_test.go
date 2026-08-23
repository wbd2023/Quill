package engine

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
)

func TestInitCreatesImmediatelyUsableRepository(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	result, err := Init(context.Background(), root, defaultPreset)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if result.Root != root || result.Preset != defaultPreset {
		t.Fatalf("Init result = %#v", result)
	}

	instance, err := New(root)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := instance.Metadata(context.Background()); err != nil {
		t.Fatalf("generated repository must prepare: %v", err)
	}
}

func TestInitRefusesSymlinkedPolicyFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	stylePath := filepath.Join(root, styleFileName)
	missingTarget := filepath.Join(root, "missing")
	if err := os.Symlink(missingTarget, stylePath); err != nil {
		t.Skipf("symlink creation unsupported: %v", err)
	}

	_, err := Init(context.Background(), root, defaultPreset)
	if err == nil {
		t.Fatal("Init succeeded with a symlinked policy file")
	}
	if _, err := os.Lstat(stylePath); err != nil {
		t.Fatalf("STYLE.md symlink must remain untouched: %v", err)
	}
	if _, err := os.Lstat(missingTarget); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Init must not create the symlink target: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(root, profile.DefaultFilename)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("quill.toml must not be written: %v", err)
	}
}

func TestWritePolicyFilesRollsBackAfterSecondWriteFailure(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	stylePath := filepath.Join(root, styleFileName)
	profilePath := filepath.Join(root, "missing", profileFileName)

	if err := writePolicyFiles(context.Background(), stylePath, profilePath, "style", "profile"); err == nil {
		t.Fatal("writePolicyFiles succeeded with a missing profile parent")
	}
	if _, err := os.Lstat(stylePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("STYLE.md must be rolled back: %v", err)
	}
}

func TestInitCancellationBeforeWritesLeavesNoPolicyFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	operationContext, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Init(operationContext, root, defaultPreset)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Init error = %v, want context.Canceled", err)
	}
	for _, name := range []string{styleFileName, profileFileName} {
		if _, statErr := os.Lstat(filepath.Join(root, name)); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("%s must not be written: %v", name, statErr)
		}
	}
}

// TestInitRefusalNamesActuallyOccupiedFile guards the audit fix: Init must name the file that
// prevented creation rather than asserting both policy files exist. Each case pre-creates exactly
// one policy file and checks the error names it and not the other.
func TestInitRefusalNamesActuallyOccupiedFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		create string
		named  string
		other  string
	}{
		{name: "only style guide", create: styleFileName, named: styleFileName, other: profileFileName},
		{name: "only profile", create: profileFileName, named: profileFileName, other: styleFileName},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, test.create), []byte("#"), 0o644); err != nil {
				t.Fatal(err)
			}

			_, err := Init(context.Background(), root, defaultPreset)
			if err == nil {
				t.Fatal("Init succeeded with an existing policy file")
			}
			if !strings.Contains(err.Error(), test.named) {
				t.Fatalf("Init error must name the occupied file %q: %v", test.named, err)
			}
			if strings.Contains(err.Error(), test.other) {
				t.Fatalf("Init error must not name the absent file %q: %v", test.other, err)
			}
		})
	}
}

// TestInitGeneratedProfileDocumentsValidListSelector guards the audit fix: the generated profile
// comment must point at the valid `quill list rules` selector, not the invalid bare `quill list`.
func TestInitGeneratedProfileDocumentsValidListSelector(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if _, err := Init(context.Background(), root, defaultPreset); err != nil {
		t.Fatalf("Init: %v", err)
	}

	body, err := os.ReadFile(filepath.Join(root, profileFileName))
	if err != nil {
		t.Fatalf("read generated profile: %v", err)
	}

	profileText := string(body)
	if !strings.Contains(profileText, "quill list rules") {
		t.Fatalf("generated profile must document the valid `quill list rules` selector:\n%s", profileText)
	}
	if strings.Contains(profileText, "quill list`") {
		t.Fatalf("generated profile must not document the invalid bare `quill list`:\n%s", profileText)
	}
}
