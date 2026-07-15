package effective_test

import (
	"testing"

	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/profiletest"
)

func TestCompileRejectsDuplicateRuleDefinitions(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules = append(definitions.Rules, definitions.Rules[0])
	_, err := effective.Compile(config, definitions)
	requireErrorContains(t, err, "duplicate rule definition")
}

func TestCompileRejectsBlankRuleDefinitionName(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules[0].Name = " "
	_, err := effective.Compile(config, definitions)
	requireErrorContains(t, err, "empty name")
}

func TestCompileRejectsBlankRuleDefinitionGroup(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules[0].Group = ""
	_, err := effective.Compile(config, definitions)
	requireErrorContains(t, err, "empty group")
}
