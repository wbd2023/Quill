package installer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/wbd2023/Quill/internal/ecosystem/golang"
	"github.com/wbd2023/Quill/internal/ecosystem/node"
	"github.com/wbd2023/Quill/internal/lockfile"
	"github.com/wbd2023/Quill/internal/toolchain"
	"github.com/wbd2023/Quill/internal/workspace"
)

const standardPermissions os.FileMode = 0o755

// Install downloads and installs the pinned external tools declared in the profile. All tools are
// attempted; failures from independent tools are collected and returned as a joined error.
func Install(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tools []toolchain.Tool,
	lockfile lockfile.Lockfile,
) (err error) {
	var errs []error
	for _, tool := range tools {
		if installErr := installTool(ctx, layout, writer, tool, lockfile); installErr != nil {
			errs = append(errs, installErr)
		}
	}

	return errors.Join(errs...)
}

func installTool(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	lockfile lockfile.Lockfile,
) (err error) {
	path := layout.BuildPath(node.BinaryDirectory(layout))

	switch install := tool.Install.(type) {

	case toolchain.NoInstall:
		return nil

	case toolchain.GoInstall:
		return golang.Install(ctx, layout, writer, tool, path)

	case toolchain.NpmInstall:
		return node.Install(ctx, layout, writer, tool, path)

	case toolchain.GitHubInstall:
		return installGitHub(ctx, layout, writer, tool, install, lockfile)

	default:
		return fmt.Errorf("unsupported install method %T for tool %s", install, tool.ID)
	}
}
