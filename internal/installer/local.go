package installer

import (
	"errors"
	"fmt"
	"os"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

func hasPinnedLocalTool(
	tool contract.Tool,
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
	statuses := runtime.InspectToolsWithEnvironment(
		[]contract.Tool{tool},
		map[string]toolchain.Capability{tool.ID: local},
		[]string{tool.ID},
		nil,
	)
	if len(statuses) != 1 {
		return false, fmt.Errorf("inspect local tool %s: missing status", tool.ID)
	}

	return statuses[0].Valid, nil
}
