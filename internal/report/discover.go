package report

import (
	"fmt"
	"io"
)

/* ------------------------------------------ Selectors ----------------------------------------- */

// List selectors name the discoverable entity kinds reported by `quill list`.
const (
	ListPacks  = "packs"
	ListRules  = "rules"
	ListTools  = "tools"
	ListScopes = "scopes"
)

// IsValidListSelector reports whether selector is one of the supported list selectors.
func IsValidListSelector(selector string) (valid bool) {
	switch selector {
	case ListPacks, ListRules, ListTools, ListScopes:
		return true
	default:
		return false
	}
}

/* -------------------------------------------- List -------------------------------------------- */

// ListResult is the presentation result of one list operation. Only the section matching the
// requested selector is populated; the others remain empty.
type ListResult struct {
	Selector string
	Packs    []ListPack
	Rules    []ListRule
	Tools    []ListTool
	Scopes   []ListScope
}

// ListPack is one rendered Pack row.
type ListPack struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Active   bool   `json:"active"`
	External bool   `json:"external"`
	Rules    int    `json:"rules"`
	Tools    int    `json:"tools"`
}

// ListRule is one rendered Rule row.
type ListRule struct {
	ID          string `json:"id"`
	Pack        string `json:"pack"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Enforcement string `json:"enforcement,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Fix         bool   `json:"fix"`
}

// ListTool is one rendered Tool row.
type ListTool struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Command  string   `json:"command,omitempty"`
	Pin      string   `json:"pin,omitempty"`
	Packs    []string `json:"packs"`
	External bool     `json:"external"`
}

// ListScope is one rendered scope row.
type ListScope struct {
	Name    string   `json:"name"`
	Roots   []string `json:"roots"`
	Default bool     `json:"default"`
}

// WriteList writes a list result in the requested format. In JSON mode it writes the full machine
// envelope tagged with command, carrying only the selected section under result.
func WriteList(
	writer io.Writer,
	command string,
	format OutputFormat,
	result ListResult,
) (err error) {
	switch format {
	case FormatText:
		return writeListText(writer, result)
	case FormatJSON:
		return writeListJSON(writer, command, result)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

/* ------------------------------------------- Explain ------------------------------------------ */

// ExplainResult is the presentation result of one explain operation.
type ExplainResult struct {
	Rule ExplainRule `json:"rule"`
}

// ExplainRule is the rendered explanation of one active Rule.
type ExplainRule struct {
	ID           string            `json:"id"`
	Pack         string            `json:"pack"`
	Name         string            `json:"name"`
	Group        string            `json:"group"`
	External     bool              `json:"external"`
	Enforcement  string            `json:"enforcement"`
	Scope        string            `json:"scope"`
	Requirements []string          `json:"requirements"`
	Check        ExplainExecution  `json:"check"`
	Fix          *ExplainExecution `json:"fix,omitempty"`
	PackConfig   map[string]any    `json:"pack_config,omitempty"`
}

// ExplainExecution is the rendered execution detail for a check or fix.
type ExplainExecution struct {
	Category string   `json:"category"`
	Tools    []string `json:"tools,omitempty"`
	FileSet  string   `json:"file_set,omitempty"`
	Language string   `json:"language,omitempty"`
	Detail   string   `json:"detail,omitempty"`
}

// WriteExplain writes an explain result in the requested format. In JSON mode it writes the full
// machine envelope tagged with command.
func WriteExplain(
	writer io.Writer,
	command string,
	format OutputFormat,
	result ExplainResult,
) (err error) {
	switch format {
	case FormatText:
		return writeExplainText(writer, result)
	case FormatJSON:
		return writeExplainJSON(writer, command, result)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}
