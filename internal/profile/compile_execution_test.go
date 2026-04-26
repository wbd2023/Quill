package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rulepack"
)

func TestCompileRejectsExecutorSpecFieldMismatch(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Rules = []policy.RuleBinding{
		{
			RuleID:         "test/bad-file-command",
			Level:          contract.LevelRequired,
			Scope:          contract.Scope("all"),
			RequirementIDs: []string{"0.4.pin-go"},
		},
	}
	config.Tools = []policy.ToolPin{{ID: rulepack.ToolShfmt, Version: "v3.12.0"}}
	definitions := contract.Definitions{
		Tools: []contract.Tool{{ID: rulepack.ToolShfmt}},
		Rules: []contract.RuleDefinition{
			{
				ID: "test/bad-file-command",
				Spec: contract.ExecutionSpec{
					Kind: rulepack.ExecutorFileCommand,
					Detail: contract.FileCommandExecution{
						ToolID: rulepack.ToolShfmt,
					},
				},
			},
		},
	}

	_, err = Compile(config, definitions)
	if err == nil || !strings.Contains(err.Error(), "file-command spec must define a file set") {
		t.Fatalf("expected file-command shape error, got %v", err)
	}
}
