package installer

import (
	"errors"
	"fmt"
	"os"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

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

	local := capability
	local.Command = path
	statuses := toolchain.InspectToolsWithEnvironment(
		[]style.Tool{tool},
		map[string]toolchain.Capability{tool.ID: local},
		[]string{tool.ID},
		nil,
		runtime.RunToolchainCommand,
	)
	if len(statuses) != 1 {
		return false, fmt.Errorf("inspect local tool %s: missing status", tool.ID)
	}

	return statuses[0].Valid, nil
}
