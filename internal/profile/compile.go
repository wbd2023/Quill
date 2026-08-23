package profile

import (
	"github.com/wbd2023/quill/internal/style"
)

// Compile validates config and builds an executable plan from definitions.
func Compile(
	config Profile,
	definitions style.Definitions,
) (plan style.Plan, err error) {
	if err := Validate(config); err != nil {
		return style.Plan{}, err
	}

	return compilePlan(config, definitions)
}
