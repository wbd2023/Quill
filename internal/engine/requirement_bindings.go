package engine

import (
	"fmt"
	"sort"
	"strings"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/styleguide"
)

// validateRequirementBindings rejects Profile rule bindings whose requirement ids are not
// documented in the parsed STYLE.md document. Requirement ids have already passed syntactic
// validation during profile loading; this check enforces that each binding references a real
// requirement, so an unknown but syntactically valid id fails before any rule or tool operation.
func validateRequirementBindings(
	config profile.Profile,
	document styleguide.Document,
) (err error) {
	documented := make(map[string]struct{}, len(document.Requirements))
	for _, requirement := range document.Requirements {
		documented[requirement.ID] = struct{}{}
	}

	seen := make(map[string]struct{})
	var unknown []string
	for _, binding := range config.Rules {
		for _, id := range binding.RequirementIDs {
			if _, ok := documented[id]; ok {
				continue
			}
			if _, duplicate := seen[id]; duplicate {
				continue
			}
			seen[id] = struct{}{}
			unknown = append(unknown, id)
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	sort.Strings(unknown)
	quoted := make([]string, len(unknown))
	for index, id := range unknown {
		quoted[index] = fmt.Sprintf("%q", id)
	}

	return fmt.Errorf(
		"style profile binds requirement id(s) %s not documented in %q",
		strings.Join(quoted, ", "),
		config.StyleGuide.Path,
	)
}
