package profile

import (
	"fmt"

	"ciphera/tools/internal/policy"
)

func validateQualitySurface(surface policy.QualitySurfaceConfig) (err error) {
	switch surface.Driver {
	case "":
		return fmt.Errorf("quality_surface.driver must not be empty")
	case policy.QualitySurfaceDriverMake:
	default:
		return fmt.Errorf("unsupported quality_surface.driver %q", surface.Driver)
	}

	return validateMakeConfig(surface.Make)
}

func validateMakeConfig(makefile policy.MakeConfig) (err error) {
	if makefile.Path == "" {
		return fmt.Errorf("quality_surface.make.path must not be empty")
	}

	if len(makefile.RequiredTargets) == 0 {
		return fmt.Errorf("quality_surface.make.required_targets must not be empty")
	}

	for _, variable := range makefile.RequiredVariables {
		if variable.Name == "" {
			return fmt.Errorf("quality_surface.make.required_variables contains an empty name")
		}
	}

	for _, target := range makefile.RequiredTargets {
		if target.Name == "" {
			return fmt.Errorf("quality_surface.make.required_targets contains an empty name")
		}

		if target.RecipeLine == "" {
			return fmt.Errorf(
				"quality_surface.make.required_targets.%s recipe_line must not be empty",
				target.Name,
			)
		}
	}

	return nil
}
