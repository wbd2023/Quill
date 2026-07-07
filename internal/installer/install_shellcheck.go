package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	goruntime "runtime"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

const (
	shellcheckDownloadRoot  = "https://github.com/koalaman/shellcheck/releases/download"
	shellcheckTempDirPrefix = "quill-shellcheck-*"
)

func installShellcheckArchive(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	path := filepath.Join(layout.ToolBinaryDirectory(), capability.Command)
	installed, err := hasPinnedLocalTool(tool, capability, path)
	if err != nil {
		return err
	}

	if installed {
		return nil
	}

	asset, err := shellcheckAssetFor(goruntime.GOOS, goruntime.GOARCH)
	if err != nil {
		return err
	}

	archive := fmt.Sprintf("shellcheck-v%s.%s.tar.xz", tool.PinnedVersion, asset.Name)
	url := shellcheckDownloadRoot + "/v" + tool.PinnedVersion + "/" + archive
	dir, err := os.MkdirTemp("", shellcheckTempDirPrefix)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	downloaded := filepath.Join(dir, archive)
	if _, err = fmt.Fprintf(
		writer,
		"Installing %s@%s...\n",
		tool.Name,
		tool.PinnedVersion,
	); err != nil {
		return err
	}

	if err = downloadFile(url, downloaded); err != nil {
		return err
	}

	if err = verifyChecksum(downloaded, asset.SHA256); err != nil {
		return err
	}

	extracted, err := extractShellcheckBinary(downloaded, dir, tool.PinnedVersion)
	if err != nil {
		return err
	}

	return copyExecutable(extracted, path)
}
