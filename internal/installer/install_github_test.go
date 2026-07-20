package installer

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/wbd2023/Quill/internal/lockfile"
	"github.com/wbd2023/Quill/internal/toolchain"
	"github.com/wbd2023/Quill/internal/workspace"
)

func TestInstallGitHubRejectsSymlinkBeforeVersionProbe(t *testing.T) {
	t.Parallel()

	layout := workspace.NewLayout(t.TempDir())
	if err := os.MkdirAll(layout.BinaryDirectory(), standardPermissions); err != nil {
		t.Fatalf("create binary directory: %v", err)
	}

	target := filepath.Join(t.TempDir(), "outside-tool")
	if err := os.WriteFile(target, []byte("outside"), standardPermissions); err != nil {
		t.Fatalf("write outside tool: %v", err)
	}
	if err := os.Symlink(target, filepath.Join(layout.BinaryDirectory(), "tool")); err != nil {
		t.Fatalf("create tool symlink: %v", err)
	}

	probeMarker := filepath.Join(layout.RepositoryRoot, "probe-called")
	tool := toolchain.Tool{
		ID:            "tool",
		Name:          "Tool",
		Command:       "tool",
		PinnedVersion: "1.0.0",
		Version: func(
			_ context.Context,
			_ toolchain.CommandRunner,
			_ map[string]string,
			_ string,
		) (version string, err error) {
			if err := os.WriteFile(probeMarker, nil, 0o600); err != nil {
				return "", err
			}
			return "1.0.0", nil
		},
	}

	err := installGitHub(
		t.Context(),
		layout,
		io.Discard,
		tool,
		toolchain.GitHubInstall{},
		lockfile.Lockfile{},
	)
	if err == nil {
		t.Fatal("expected symlink destination to be rejected")
	}

	if _, statErr := os.Stat(probeMarker); !os.IsNotExist(statErr) {
		t.Fatalf("version probe marker exists after rejection: %v", statErr)
	}
}
