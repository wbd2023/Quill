package effective_test

import (
	"slices"
	"testing"

	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/profilefixture"
)

func TestCompileInfersTargetsFromRuleScope(t *testing.T) {
	t.Parallel()

	config := profilefixture.Config()
	want := []string{profilefixture.Target, profilefixture.OtherTarget}

	compiled, err := effective.Compile(config, profilefixture.TargetCommandDefinitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	rule := compiled.Rules[0]
	if !slices.Equal(rule.Check.Targets(), want) {
		t.Fatalf("check targets = %v, want %v", rule.Check.Targets(), want)
	}

	if !slices.Equal(rule.Fix.Targets(), want) {
		t.Fatalf("fix targets = %v, want %v", rule.Fix.Targets(), want)
	}
}

func TestCompileRejectsMissingInferredTargets(t *testing.T) {
	t.Parallel()

	config := profilefixture.Config()
	config.Targets = nil

	_, err := effective.Compile(config, profilefixture.TargetCommandDefinitions())
	requireErrorContains(t, err, "has no test targets")
}

func TestCompileRejectsTargetCheckWithoutCheckPaths(t *testing.T) {
	t.Parallel()

	config := profilefixture.Config()
	config.Targets[0].CheckPaths = nil

	_, err := effective.Compile(config, profilefixture.TargetCheckDefinitions())
	requireErrorContains(t, err, "must define check_paths")
}
