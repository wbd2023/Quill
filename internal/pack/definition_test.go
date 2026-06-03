package pack

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/toolchain"
)

func TestCloneDefinitionReturnsIndependentCopy(t *testing.T) {
	original := Definition{
		ID:   "custom",
		Name: "Custom",
		Tools: []toolchain.Capability{
			{ID: "tool", Command: "tool"},
		},
		Rules: []contract.RuleDefinition{
			{
				ID: "custom/rule",
				Check: contract.ExecutionSpec{
					Kind: contract.ExecutionFileCommand,
					Detail: contract.FileCommandExecution{
						Arguments: []string{"-w"},
					},
				},
			},
		},
		FileSets: policy.FileSets{
			{
				Name: "source",
				Include: policy.FileSetInclude{
					Extensions: []string{".go"},
				},
			},
		},
	}

	clone := CloneDefinition(original)
	clone.Tools[0].Command = "changed"
	clone.FileSets[0].Include.Extensions[0] = ".txt"

	execution := clone.Rules[0].Check.Detail.(contract.FileCommandExecution)
	execution.Arguments[0] = "-changed"

	if got := original.Tools[0].Command; got != "tool" {
		t.Fatalf("original tool command = %q, want tool", got)
	}

	if got := original.FileSets[0].Include.Extensions[0]; got != ".go" {
		t.Fatalf("original file set extension = %q, want .go", got)
	}

	execution = original.Rules[0].Check.Detail.(contract.FileCommandExecution)
	if got := execution.Arguments[0]; got != "-w" {
		t.Fatalf("original rule argument = %q, want -w", got)
	}
}
