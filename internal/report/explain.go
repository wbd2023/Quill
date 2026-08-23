package report

import (
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/engine"
)

// ExplainResult is the presentation result of one explain operation.
type ExplainResult struct {
	Rule ExplainRule `json:"rule"`
}

// ExplainRule is the rendered explanation of one active Rule.
type ExplainRule struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Group   string            `json:"group"`
	Pack    ExplainPack       `json:"pack"`
	Binding ExplainBinding    `json:"binding"`
	Check   ExplainExecution  `json:"check"`
	Fix     *ExplainExecution `json:"fix,omitempty"`
}

// ExplainPack is the public Pack context for an explained Rule.
type ExplainPack struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Provenance string         `json:"provenance"`
	Policy     map[string]any `json:"policy,omitempty"`
}

// ExplainBinding is the active Profile binding for an explained Rule.
type ExplainBinding struct {
	Enforcement  string   `json:"enforcement"`
	Scope        string   `json:"scope"`
	Requirements []string `json:"requirements"`
}

// ExplainExecution is the rendered execution summary for a check or fix.
type ExplainExecution struct {
	Category string   `json:"category"`
	Tools    []string `json:"tools,omitempty"`
	FileSet  string   `json:"file_set,omitempty"`
	Language string   `json:"language,omitempty"`
}

// WriteExplain writes an explain result in the requested format. In JSON mode it writes the full
// machine envelope identified by metadata.
func WriteExplain(
	writer io.Writer,
	metadata EnvelopeMetadata,
	format OutputFormat,
	result ExplainResult,
) (err error) {
	switch format {
	case FormatText:
		return writeExplainText(writer, result)
	case FormatJSON:
		return writeExplainJSON(writer, metadata, result)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

// NewExplainResult converts an active engine explanation into the explicit explain protocol DTO.
func NewExplainResult(explanation engine.ExplainResult) (result ExplainResult) {
	rule := explanation.Rule
	result.Rule = ExplainRule{
		ID:    rule.ID,
		Name:  rule.Name,
		Group: string(rule.Group),
		Pack: ExplainPack{
			ID:         explanation.Pack.ID,
			Name:       explanation.Pack.Name,
			Provenance: string(explanation.Pack.Provenance),
		},
		Binding: ExplainBinding{
			Enforcement:  string(rule.Enforcement),
			Scope:        string(rule.Scope),
			Requirements: rule.RequirementIDs,
		},
		Check: newExplainExecution(rule.Check),
	}

	if rule.HasFix {
		fix := newExplainExecution(rule.Fix)
		result.Rule.Fix = &fix
	}
	if explanation.PackPolicy != nil {
		result.Rule.Pack.Policy = map[string]any(explanation.PackPolicy)
	}

	return result
}

func newExplainExecution(detail engine.ExecutionDetail) (execution ExplainExecution) {
	return ExplainExecution{
		Category: detail.Category,
		Tools:    detail.ToolIDs,
		FileSet:  detail.FileSet,
		Language: detail.Language,
	}
}
