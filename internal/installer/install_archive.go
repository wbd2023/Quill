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

func installArchive(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	spec := capability.Archive
	if spec == nil {
		return fmt.Errorf("tool %s has no archive spec", tool.ID)
	}

	path := filepath.Join(layout.ToolBinaryDirectory(), capability.Command)
	installed, err := hasPinnedLocalTool(tool, capability, path)
	if err != nil {
		return err
	}

	if installed {
		return nil
	}

	platform, ok := spec.Platforms[goruntime.GOOS+"/"+goruntime.GOARCH]
	if !ok {
		return fmt.Errorf(
			"unsupported platform %s/%s for tool %s",
			goruntime.GOOS,
			goruntime.GOARCH,
			tool.ID,
		)
	}

	url := spec.URL(tool.PinnedVersion, platform)
	hash, err := archiveHashFor(tool.ID, goruntime.GOOS, goruntime.GOARCH)
	if err != nil {
		return err
	}

	dir, err := os.MkdirTemp("", "quill-archive-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	archive := filepath.Join(dir, "archive")
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

	extracted, err := extractBinary(archive, dir, *spec, tool.PinnedVersion)
	if err != nil {
		return err
	}

	return copyExecutable(extracted, path)
}
