package vocabulary

import (
	"fmt"
	"slices"
)

func rejectUnknownFields(
	section map[string]any,
	field string,
	allowed ...string,
) (err error) {
	known := make(map[string]bool, len(allowed))
	for _, name := range allowed {
		known[name] = true
	}

	var unknown []string
	for name := range section {
		if !known[name] {
			unknown = append(unknown, name)
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	slices.Sort(unknown)
	return fmt.Errorf("%s.%s is not supported", field, unknown[0])
}
