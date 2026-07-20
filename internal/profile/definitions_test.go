package profile

import (
	"testing"

	"github.com/wbd2023/Quill/internal/profile/internal/profiletest"
)

func TestCompileRejectsDuplicateRuleDefinitions(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules = append(definitions.Rules, definitions.Rules[0])
	_, err := compilePlan(config, definitions)
	requireErrorContainsInternal(t, err, "duplicate rule definition")
}

func TestCompileRejectsBlankRuleDefinitionName(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules[0].Name = " "
	_, err := compilePlan(config, definitions)
	requireErrorContainsInternal(t, err, "empty name")
}

func TestCompileRejectsBlankRuleDefinitionGroup(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules[0].Group = ""
	_, err := compilePlan(config, definitions)
	requireErrorContainsInternal(t, err, "empty group")
}
