package installer

import (
	"errors"
	"fmt"
	"io"
	"os"

	"ciphera/tools/internal/ecosystem/golang"
	"ciphera/tools/internal/ecosystem/node"
	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

const standardPermissions os.FileMode = 0o755

// Install downloads and installs the pinned external tools declared in the profile. All tools are
// attempted; failures from independent tools are collected and returned as a joined error.
func Install(
	layout runtime.Layout,
	writer io.Writer,
	tools []toolchain.Tool,
	lockfile lockfile.Lockfile,
) (err error) {
	var errs []error
	for _, tool := range tools {
		if installErr := installTool(layout, writer, tool, lockfile); installErr != nil {
			errs = append(errs, installErr)
		}
	}

	return errors.Join(errs...)
}

func installTool(
	layout runtime.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	lockfile lockfile.Lockfile,
) (err error) {
	searchPath := runtime.SearchPath(
		layout.ToolBinaryDirectory(),
		node.BinaryDirectory(layout),
	)

	switch install := tool.Install.(type) {

	case toolchain.NoInstall:
		return nil

	case toolchain.GoInstall:
		return golang.Install(layout, writer, tool, install, searchPath)

	case toolchain.NpmInstall:
		return node.Install(layout, writer, tool, install, searchPath)

	case toolchain.GitHubInstall:
		return installGitHub(layout, writer, tool, install, lockfile)

	default:
		return fmt.Errorf("unsupported install method %T for tool %s", install, tool.ID)
	}
}
