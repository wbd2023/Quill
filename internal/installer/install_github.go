package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	goruntime "runtime"

	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

func installGitHub(
	layout runtime.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	install toolchain.GitHubInstall,
	lockfile lockfile.Lockfile,
) (err error) {
	path := filepath.Join(layout.ToolBinaryDirectory(), tool.Command)
	installed, err := toolchain.IsInstalled(tool, path)
	if err != nil {
		return err
	}

	if installed {
		return nil
	}

	platform, ok := install.Platforms[goruntime.GOOS+"/"+goruntime.GOARCH]
	if !ok {
		return fmt.Errorf(
			"unsupported platform %s/%s for tool %s",
			goruntime.GOOS,
			goruntime.GOARCH,
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
	hash, err := lockfile.HashFor(tool.ID, tool.PinnedVersion, goruntime.GOOS, goruntime.GOARCH)
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

	if err = downloadFile(url, archive); err != nil {
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
