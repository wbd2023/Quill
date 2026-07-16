package profile

import "fmt"

func validateEnabledPacks(enabledPacks []string) (err error) {
	if len(enabledPacks) == 0 {
		return fmt.Errorf("packs.enabled must not be empty")
	}

	seen := make(map[string]bool, len(enabledPacks))
	for _, pack := range enabledPacks {
		if isBlank(pack) {
			return fmt.Errorf("packs.enabled contains an empty pack")
		}

		if seen[pack] {
			return fmt.Errorf("packs.enabled contains duplicate pack %q", pack)
		}

		seen[pack] = true
	}

	return nil
}
