package pack

import (
	"testing"

	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/style"
)

func TestCloneDefinitionReturnsIndependentCopy(t *testing.T) {
	original := Definition{
		ID:      "custom",
		Name:    "Custom",
		ToolIDs: []string{"tool"},
		Rules: []style.RuleDefinition{
			{
				ID: "custom/rule",
				Check: style.FileCommandExecution{
					Arguments: []string{"-w"},
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
	clone.ToolIDs[0] = "changed"
	clone.FileSets[0].Include.Extensions[0] = ".txt"

	execution := clone.Rules[0].Check.(style.FileCommandExecution)
	execution.Arguments[0] = "-changed"

	if got := original.ToolIDs[0]; got != "tool" {
		t.Fatalf("original tool id = %q, want tool", got)
	}

	if got := original.FileSets[0].Include.Extensions[0]; got != ".go" {
		t.Fatalf("original file set extension = %q, want .go", got)
	}

	execution = original.Rules[0].Check.(style.FileCommandExecution)
	if got := execution.Arguments[0]; got != "-w" {
		t.Fatalf("original rule argument = %q, want -w", got)
	}
}
