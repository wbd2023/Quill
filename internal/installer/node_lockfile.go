package installer

import (
	"encoding/json"
	"fmt"
	"os"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func validatePackageLock(
	path string,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var document struct {
		Packages map[string]struct {
			Version string `json:"version"`
		} `json:"packages"`
	}
	if err = json.Unmarshal(contents, &document); err != nil {
		return err
	}

	packageEntry, found := document.Packages["node_modules/"+capability.InstallSource]
	if !found {
		return fmt.Errorf("package lock does not contain %s", capability.InstallSource)
	}

	if packageEntry.Version != tool.PinnedVersion {
		return fmt.Errorf(
			"package lock pins %s@%s, profile pins %s",
			capability.InstallSource,
			packageEntry.Version,
			tool.PinnedVersion,
		)
	}

	return nil
}
