package runtime

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ulikunitz/xz"

	"ciphera/tools/internal/contract"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	shellcheckDownloadRoot  = "https://github.com/koalaman/shellcheck/releases/download"
	shellcheckTempDirPrefix = "style-platform-shellcheck-*"
)

var shellcheckAssets = map[string]shellcheckAsset{
	"darwin/amd64": {
		Name:   "darwin.x86_64",
		SHA256: "ef27684f23279d112d8ad84e0823642e43f838993bbb8c0963db9b58a90464c2",
	},
	"darwin/arm64": {
		Name:   "darwin.aarch64",
		SHA256: "bbd2f14826328eee7679da7221f2bc3afb011f6a928b848c80c321f6046ddf81",
	},
	"linux/amd64": {
		Name:   "linux.x86_64",
		SHA256: "6c881ab0698e4e6ea235245f22832860544f17ba386442fe7e9d629f8cbedf87",
	},
	"linux/arm64": {
		Name:   "linux.aarch64",
		SHA256: "324a7e89de8fa2aed0d0c28f3dab59cf84c6d74264022c00c22af665ed1a09bb",
	},
}

type shellcheckAsset struct {
	Name   string
	SHA256 string
}

/* ----------------------------------------- Shellcheck ----------------------------------------- */

func installShellcheckTool(layout Layout, writer io.Writer, tool contract.Tool) (err error) {
	localPath := filepath.Join(layout.ToolBinDir, tool.Command)
	localVersion, found, err := inspectLocalToolVersion(tool, localPath)
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

/* ------------------------------------------- Archive ------------------------------------------ */

func extractShellcheckBinary(
	archivePath string,
	destination string,
	version string,
) (binaryPath string, err error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close shellcheck archive %q: %w", archivePath, closeErr)
		}
	}()

	xzReader, err := xz.NewReader(file)
	if err != nil {
		return "", err
	}

	expectedName := path.Join(shellcheckArchiveRoot(version), "shellcheck")
	targetPath := filepath.Join(destination, filepath.FromSlash(expectedName))
	foundBinary := false

	tarReader := tar.NewReader(xzReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			if !foundBinary {
				return "", fmt.Errorf("shellcheck archive missing %s", expectedName)
			}

			return targetPath, nil
		}

		if err != nil {
			return "", err
		}

		cleanName, err := validateShellcheckArchiveEntry(header, version)
		if err != nil {
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue

		case tar.TypeReg:
			if cleanName != expectedName {
				continue
			}

			foundBinary = true
			if err = os.MkdirAll(filepath.Dir(targetPath), defaultDirectoryMode); err != nil {
				return "", err
			}
			targetFile, err := os.OpenFile(
				targetPath,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				defaultDirectoryMode,
			)
			if err != nil {
				return "", err
			}

			if _, err = io.Copy(targetFile, tarReader); err != nil {
				if closeErr := targetFile.Close(); closeErr != nil {
					return "", fmt.Errorf(
						"copy shellcheck file %q: %w",
						targetPath,
						errors.Join(err, closeErr),
					)
				}
				return "", err
			}

			if err = targetFile.Close(); err != nil {
				return "", err
			}

		default:
			return "", fmt.Errorf("unsupported shellcheck archive entry %q", header.Name)
		}
	}
}

func validateShellcheckArchiveEntry(
	header *tar.Header,
	version string,
) (name string, err error) {
	switch header.Typeflag {
	case tar.TypeSymlink, tar.TypeLink:
		return "", fmt.Errorf("shellcheck archive contains link entry %q", header.Name)
	}

	rawName := header.Name
	if header.Typeflag == tar.TypeDir {
		rawName = strings.TrimSuffix(rawName, "/")
	}

	name = path.Clean(rawName)
	if name == "." ||
		name != rawName ||
		path.IsAbs(rawName) ||
		strings.HasPrefix(name, "../") ||
		strings.Contains(name, "/../") {
		return "", fmt.Errorf("unsafe shellcheck archive path %q", header.Name)
	}

	root := shellcheckArchiveRoot(version)
	switch name {
	case root,
		path.Join(root, "LICENSE.txt"),
		path.Join(root, "README.txt"),
		path.Join(root, "shellcheck"):
		return name, nil
	default:
		return "", fmt.Errorf("unexpected shellcheck archive entry %q", header.Name)
	}
}

/* -------------------------------------- Platform Mapping -------------------------------------- */

func shellcheckAssetName(goos string, goarch string) (name string, err error) {
	asset, err := shellcheckAssetFor(goos, goarch)
	if err != nil {
		return "", err
	}

	return asset.Name, nil
}

func shellcheckAssetFor(goos string, goarch string) (asset shellcheckAsset, err error) {
	asset, found := shellcheckAssets[goos+"/"+goarch]
	if !found {
		return shellcheckAsset{}, fmt.Errorf("unsupported shellcheck platform: %s/%s", goos, goarch)
	}

	if asset.SHA256 == "" {
		return shellcheckAsset{}, fmt.Errorf("missing shellcheck checksum for %s/%s", goos, goarch)
	}

	return asset, nil
}

func shellcheckArchiveRoot(version string) (root string) {
	return "shellcheck-v" + version
}

/* ------------------------------------- Executable Copying ------------------------------------- */

func copyExecutable(source string, destination string) (err error) {
	input, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	return os.WriteFile(destination, input, defaultDirectoryMode)
}
