package installer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/wbd2023/quill/internal/ecosystem/golang"
	"github.com/wbd2023/quill/internal/ecosystem/node"
	"github.com/wbd2023/quill/internal/lockfile"
	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/toolchain"
	"github.com/wbd2023/quill/internal/workspace"
)

const standardPermissions os.FileMode = 0o755

/* -------------------------------------- Tool Installation ------------------------------------- */

// Install downloads and installs the pinned external tools declared in the profile. Independent
// tool failures are collected, but cancellation stops the operation before another tool starts.
func Install(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tools []toolchain.Tool,
	lockfile lockfile.Lockfile,
) (err error) {
	return installWith(ctx, layout, writer, tools, lockfile, installTool)
}

type toolInstaller func(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	lockfile lockfile.Lockfile,
) error

func installWith(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tools []toolchain.Tool,
	lockfile lockfile.Lockfile,
	install toolInstaller,
) (err error) {
	var errs []error
	for _, tool := range tools {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return errors.Join(append(errs, ctxErr)...)
		}

		if installErr := install(ctx, layout, writer, tool, lockfile); installErr != nil {
			errs = append(errs, installErr)
		}
	}

	return errors.Join(errs...)
}

/* ------------------------------------- Install Strategies ------------------------------------- */

func installTool(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	lockfile lockfile.Lockfile,
) (err error) {
	switch install := tool.Install.(type) {

	case toolchain.NoInstall:
		return nil

	case toolchain.GoInstall:
		installed, err := toolchain.IsInstalled(
			ctx,
			process.Runner{},
			tool,
			golang.BinaryPath(layout, tool.Command),
		)
		if err != nil {
			return err
		}
		if installed {
			return nil
		}

		if err = prepareGoInstall(layout); err != nil {
			return err
		}
		path, err := bootstrapPath(layout)
		if err != nil {
			return err
		}
		executable, err := resolveBootstrap(layout, path, "go")
		if err != nil {
			return err
		}
		return golang.Install(ctx, layout, writer, tool, executable, path, process.Runner{})

	case toolchain.NpmInstall:
		installed, err := toolchain.IsInstalled(
			ctx,
			process.Runner{},
			tool,
			node.BinaryPath(layout, tool.Command),
		)
		if err != nil {
			return err
		}
		if installed {
			return nil
		}

		if err = prepareNpmInstall(layout); err != nil {
			return err
		}
		path, err := bootstrapPath(layout)
		if err != nil {
			return err
		}
		executable, err := resolveBootstrap(layout, path, "npm")
		if err != nil {
			return err
		}
		return node.Install(ctx, layout, writer, tool, executable, path, process.Runner{})

	case toolchain.GitHubInstall:
		return installGitHub(ctx, layout, writer, tool, install, lockfile)

	default:
		return fmt.Errorf("unsupported install method %T for tool %s", install, tool.ID)
	}
}
