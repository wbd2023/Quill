package pack

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
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
				Check: style.FileCommand{
					Arguments: []string{"-w"},
				},
			},
		},
		FileSets: profile.FileSets{
			{
				Name: "source",
				Include: profile.FileSetInclude{
					Extensions: []string{".go"},
				},
			},
		},
	}

	clone := CloneDefinition(original)
	clone.ToolIDs[0] = "changed"
	clone.FileSets[0].Include.Extensions[0] = ".txt"

	execution := clone.Rules[0].Check.(style.FileCommand)
	execution.Arguments[0] = "-changed"

	if got := original.ToolIDs[0]; got != "tool" {
		t.Fatalf("original tool id = %q, want tool", got)
	}

	if got := original.FileSets[0].Include.Extensions[0]; got != ".go" {
		t.Fatalf("original file set extension = %q, want .go", got)
	}

	execution = original.Rules[0].Check.(style.FileCommand)
	if got := execution.Arguments[0]; got != "-w" {
		t.Fatalf("original rule argument = %q, want -w", got)
	}
}
