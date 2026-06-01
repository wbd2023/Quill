package effective_test

import (
	"slices"
	"testing"

	"ciphera/tools/internal/profile/effective"
	"ciphera/tools/internal/profile/internal/fixture"
)

func TestCompileInfersTargetsFromRuleScope(t *testing.T) {
	t.Parallel()

	config := fixture.Config()
	want := []string{fixture.Target, fixture.OtherTarget}

	compiled, err := effective.Compile(config, fixture.TargetCommandDefinitions())
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

	config := fixture.Config()
	config.Targets = nil

	_, err := effective.Compile(config, fixture.TargetCommandDefinitions())
	requireErrorContains(t, err, "has no test targets")
}

func TestCompileRejectsTargetCheckWithoutCheckPaths(t *testing.T) {
	t.Parallel()

	config := fixture.Config()
	config.Targets[0].CheckPaths = nil

	_, err := effective.Compile(config, fixture.TargetCheckDefinitions())
	requireErrorContains(t, err, "must define check_paths")
}
