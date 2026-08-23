package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
)

func TestCompileRejectsDuplicateRuleDefinitions(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules = append(definitions.Rules, definitions.Rules[0])
	_, err := profile.Compile(config, definitions)
	requireErrorContains(t, err, "duplicate rule definition")
}

func TestCompileRejectsBlankRuleDefinitionName(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules[0].Name = " "
	_, err := profile.Compile(config, definitions)
	requireErrorContains(t, err, "empty name")
}

func TestCompileRejectsBlankRuleDefinitionGroup(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	definitions := profiletest.Definitions()
	definitions.Rules[0].Group = ""
	_, err := profile.Compile(config, definitions)
	requireErrorContains(t, err, "empty group")
}
