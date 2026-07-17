package installer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/process"
	"ciphera/tools/internal/toolchain"
	"ciphera/tools/internal/workspace"
)

func installGitHub(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	install toolchain.GitHubInstall,
	lockfile lockfile.Lockfile,
) (err error) {
	path := filepath.Join(layout.BinaryDirectory(), tool.Command)
	installed, err := toolchain.IsInstalled(ctx, process.Runner{}, tool, path)
	if err != nil {
		return err
	}

	if installed {
		return nil
	}

	platform, ok := install.Platforms[runtime.GOOS+"/"+runtime.GOARCH]
	if !ok {
		return fmt.Errorf(
			"unsupported platform %s/%s for tool %s",
			runtime.GOOS,
			runtime.GOARCH,
			tool.ID,
		)
	}

	tag := fmt.Sprintf(install.Tag, tool.PinnedVersion)
	asset := fmt.Sprintf(install.Asset, tag, platform)
	url := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s/%s",
		install.Owner,
		install.Repository,
		tag,
		asset,
	)
	hash, err := lockfile.HashFor(tool.ID, tool.PinnedVersion, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	dir, err := os.MkdirTemp("", "quill-github-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	archive := filepath.Join(dir, asset)
	if _, err = fmt.Fprintf(
		writer,
		"Installing %s@%s...\n",
		tool.Name,
		tool.PinnedVersion,
	); err != nil {
		return err
	}

	if err = downloadFile(ctx, url, archive); err != nil {
		return err
	}

	if err = verifyChecksum(archive, hash); err != nil {
		return err
	}

	extracted, err := extractBinary(archive, dir, install, tool.PinnedVersion)
	if err != nil {
		return err
	}

	return copyExecutable(extracted, path)
}
