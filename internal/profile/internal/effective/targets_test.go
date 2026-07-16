package effective_test

import (
	"slices"
	"testing"

	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/profiletest"
	"ciphera/tools/internal/style"
)

func TestCompileInfersTargetsFromRuleScope(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	want := []string{profiletest.Target, profiletest.OtherTarget}

	compiled, err := effective.Compile(config, profiletest.TargetCommandDefinitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	rule := compiled.Rules[0]
	checkJob := rule.Check.(style.TargetCommandJob)
	if !slices.Equal(checkJob.Targets, want) {
		t.Fatalf("check targets = %v, want %v", checkJob.Targets, want)
	}

	fixJob := rule.Fix.(style.TargetCommandJob)
	if !slices.Equal(fixJob.Targets, want) {
		t.Fatalf("fix targets = %v, want %v", fixJob.Targets, want)
	}
}

func TestCompileRejectsMissingInferredTargets(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Targets = nil

	_, err := effective.Compile(config, profiletest.TargetCommandDefinitions())
	requireErrorContains(t, err, "has no test targets")
}

func TestCompileRejectsTargetCheckWithoutCheckPaths(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Targets[0].CheckPaths = nil

	_, err := effective.Compile(config, profiletest.TargetCheckDefinitions())
	requireErrorContains(t, err, "must define check_paths")
}
