package installer

import (
	"errors"
	"fmt"
	"io"
	"os"

	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const standardPermissions os.FileMode = 0o755

/* ---------------------------------------- Installation ---------------------------------------- */

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
	switch install := tool.Install.(type) {

	case toolchain.NoInstall:
		return nil

	case toolchain.GoInstall:
		return installGo(layout, writer, tool, install)

	case toolchain.NpmInstall:
		return installNpm(layout, writer, tool, install)

	case toolchain.GitHubInstall:
		return installGitHub(layout, writer, tool, install, lockfile)

	default:
		return fmt.Errorf("unsupported install method %T for tool %s", install, tool.ID)
	}
}

/* -------------------------------------- Idempotency Probe ------------------------------------- */

// hasPinnedLocalTool reports whether a tool matching the pinned version is already installed at the
// given path.
func hasPinnedLocalTool(
	tool toolchain.Tool,
	path string,
) (installed bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	probe := tool
	probe.Command = path
	statuses := toolchain.InspectTools(
		map[string]toolchain.Tool{tool.ID: probe},
		nil,
	)
	if len(statuses) != 1 {
		return false, fmt.Errorf("inspect local tool %s: missing status", tool.ID)
	}

	return statuses[0].Valid, nil
}
