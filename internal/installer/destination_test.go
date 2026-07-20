package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyExecutableRejectsSymlinkDestination(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	source := filepath.Join(directory, "source")
	target := filepath.Join(directory, "target")
	destination := filepath.Join(directory, "destination")

	if err := os.WriteFile(source, []byte("replacement"), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.WriteFile(target, []byte("protected"), 0o600); err != nil {
		t.Fatalf("write target: %v", err)
	}
	if err := os.Symlink(target, destination); err != nil {
		t.Fatalf("create destination symlink: %v", err)
	}

	if err := copyExecutable(directory, source, destination); err == nil {
		t.Fatal("expected symlink destination to be rejected")
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(content) != "protected" {
		t.Fatalf("target content = %q, want protected", content)
	}

	info, err := os.Lstat(destination)
	if err != nil {
		t.Fatalf("inspect destination: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("destination mode = %v, want symlink", info.Mode())
	}
}

func TestWriteExecutableRejectsSymlinkedParentDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, ".cache")); err != nil {
		t.Fatalf("create cache symlink: %v", err)
	}

	destination := filepath.Join(root, ".cache", "quill", "bin", "tool")
	if err := writeExecutable(root, destination, strings.NewReader("replacement")); err == nil {
		t.Fatal("expected symlinked parent directory to be rejected")
	}

	outsideDestination := filepath.Join(outside, "quill", "bin", "tool")
	if _, err := os.Stat(outsideDestination); !os.IsNotExist(err) {
		t.Fatalf("outside destination exists after rejection: %v", err)
	}
}

func TestWriteExecutableRejectsDestinationOutsideRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	destination := filepath.Join(t.TempDir(), "tool")
	if err := writeExecutable(root, destination, strings.NewReader("replacement")); err == nil {
		t.Fatal("expected destination outside root to be rejected")
	}

	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("outside destination exists after rejection: %v", err)
	}
}
