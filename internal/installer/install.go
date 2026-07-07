package installer

import (
	"errors"
	"fmt"
	"io"
	"os"

	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
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
	tools []style.Tool,
	capabilities map[string]toolchain.Capability,
	lockfile lockfile.Lockfile,
) (err error) {
	var errs []error
	for _, tool := range tools {
		capability, found := capabilities[tool.ID]
		if !found {
			errs = append(errs, fmt.Errorf("missing tool capability %q", tool.ID))
			continue
		}

		if installErr := installTool(
			layout, writer, tool, capability, lockfile,
		); installErr != nil {
			errs = append(errs, installErr)
		}
	}

	err = errors.Join(errs...)
	return err
}

func installTool(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
	lockfile lockfile.Lockfile,
) (err error) {
	switch capability.InstallKind {

	case toolchain.InstallKindNone:
		return nil

	case toolchain.InstallKindGoBinary:
		return installGoBinary(layout, writer, tool, capability)

	case toolchain.InstallKindNodePackage:
		return installNodePackage(layout, writer, tool, capability)

	case toolchain.InstallKindArchive:
		return installArchive(layout, writer, tool, capability, lockfile)

	default:
		return fmt.Errorf(
			"unsupported install strategy %q for tool %s",
			capability.InstallKind,
			tool.ID,
		)
	}
}

/* -------------------------------------- Idempotency Probe ------------------------------------- */

// hasPinnedLocalTool reports whether a tool matching the pinned version is already installed at the
// given path.
func hasPinnedLocalTool(
	tool style.Tool,
	capability toolchain.Capability,
	path string,
) (installed bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	probe := capability
	probe.Command = path
	statuses := toolchain.InspectToolsWithEnvironment(
		[]style.Tool{tool},
		map[string]toolchain.Capability{tool.ID: probe},
		[]string{tool.ID},
		nil,
		runtime.RunToolchainCommand,
	)
	if len(statuses) != 1 {
		return false, fmt.Errorf("inspect local tool %s: missing status", tool.ID)
	}

	return statuses[0].Valid, nil
}
