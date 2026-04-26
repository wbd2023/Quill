package runtime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

const (
	shellcheckDownloadRoot  = "https://github.com/koalaman/shellcheck/releases/download"
	shellcheckTempDirPrefix = "style-platform-shellcheck-*"
)

func installShellcheckTool(
	layout Layout,
	writer io.Writer,
	tool contract.Tool,
	capability toolchain.Capability,
) (err error) {
	localPath := filepath.Join(layout.ToolBinDir, capability.Command)
	localVersion, found, err := inspectLocalToolVersion(tool, capability, localPath)
	if err != nil {
		return err
	}

	if found && matchesPinnedVersion(localVersion, tool.PinnedVersion) {
		return nil
	}

	asset, err := shellcheckAssetFor(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	archiveName := fmt.Sprintf("shellcheck-v%s.%s.tar.xz", tool.PinnedVersion, asset.Name)
	versionRoot := shellcheckDownloadRoot + "/v" + tool.PinnedVersion
	tempDir, err := os.MkdirTemp("", shellcheckTempDirPrefix)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	archivePath := filepath.Join(tempDir, archiveName)
	if _, err = fmt.Fprintln(writer, "Installing shellcheck from release archive..."); err != nil {
		return err
	}

	if err = downloadFile(versionRoot+"/"+archiveName, archivePath); err != nil {
		return err
	}

	if err = verifyFileChecksum(archivePath, archiveName, asset.SHA256); err != nil {
		return err
	}

	sourcePath, err := extractShellcheckBinary(archivePath, tempDir, tool.PinnedVersion)
	if err != nil {
		return err
	}

	return copyExecutable(sourcePath, filepath.Join(layout.ToolBinDir, "shellcheck"))
}
