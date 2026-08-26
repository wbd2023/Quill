package report

import (
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/engine"
)

/* ------------------------------------------ Selectors ----------------------------------------- */

// List selectors name the discoverable entity kinds reported by `quill list`.
const (
	ListPacks  = "packs"
	ListRules  = "rules"
	ListTools  = "tools"
	ListScopes = "scopes"
)

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
	ID         string `json:"id"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	Provenance string `json:"provenance"`
	Rules      int    `json:"rules"`
	Tools      int    `json:"tools"`
}

// ListRule is one rendered Rule row.
type ListRule struct {
	ID          string `json:"id"`
	Pack        string `json:"pack"`
	Provenance  string `json:"provenance"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Enforcement string `json:"enforcement,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Fix         bool   `json:"fix"`
}

// ListTool is one rendered Tool row.
type ListTool struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Command string   `json:"command,omitempty"`
	Pin     string   `json:"pin,omitempty"`
	Packs   []string `json:"packs"`
}

// ListScope is one rendered scope row.
type ListScope struct {
	Name    string   `json:"name"`
	Roots   []string `json:"roots"`
	Default bool     `json:"default"`
}

// WriteList writes a list result in the requested format. In JSON mode it writes the full
// machine envelope identified by metadata, carrying only the selected section under result.
func WriteList(
	writer io.Writer,
	metadata EnvelopeMetadata,
	format OutputFormat,
	result ListResult,
) (err error) {
	switch format {
	case FormatText:
		return writeListText(writer, result)
	case FormatJSON:
		return writeListJSON(writer, metadata, result)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

// NewListResult converts the engine's metadata snapshot into the explicit list protocol DTO.
func NewListResult(snapshot engine.MetadataSnapshot, selector string) (result ListResult) {
	result = ListResult{Selector: selector}
	switch selector {
	case ListPacks:
		result.Packs = newListPacks(snapshot.Packs)

	case ListRules:
		result.Rules = newListRules(snapshot.Rules, snapshot.Packs)

	case ListTools:
		result.Tools = newListTools(snapshot.Tools)

	case ListScopes:
		result.Scopes = newListScopes(snapshot.Scopes)
	}

	return result
}

func newListPacks(packs []engine.PackMetadata) (rows []ListPack) {
	rows = make([]ListPack, 0, len(packs))
	for _, pack := range packs {
		rows = append(rows, ListPack{
			ID:         pack.ID,
			Name:       pack.Name,
			Active:     pack.Enabled,
			Provenance: string(pack.Provenance),
			Rules:      len(pack.RuleIDs),
			Tools:      len(pack.ToolIDs),
		})
	}

	return rows
}

func newListRules(
	rules []engine.RuleMetadata,
	packs []engine.PackMetadata,
) (rows []ListRule) {
	provenance := make(map[string]string, len(packs))
	for _, pack := range packs {
		provenance[pack.ID] = string(pack.Provenance)
	}

	rows = make([]ListRule, 0, len(rules))
	for _, rule := range rules {
		row := ListRule{
			ID:         rule.ID,
			Pack:       rule.PackID,
			Provenance: provenance[rule.PackID],
			Name:       rule.Name,
			Active:     rule.Enabled,
			Fix:        rule.HasFix,
		}
		if rule.Enabled {
			row.Enforcement = string(rule.Enforcement)
			row.Scope = string(rule.Scope)
		}
		rows = append(rows, row)
	}

	return rows
}

func newListTools(tools []engine.ToolMetadata) (rows []ListTool) {
	rows = make([]ListTool, 0, len(tools))
	for _, tool := range tools {
		rows = append(rows, ListTool{
			ID:      tool.ID,
			Name:    tool.Name,
			Command: tool.Command,
			Pin:     tool.PinnedVersion,
			Packs:   tool.PackIDs,
		})
	}

	return rows
}

func newListScopes(scopes []engine.ScopeMetadata) (rows []ListScope) {
	rows = make([]ListScope, 0, len(scopes))
	for _, scope := range scopes {
		rows = append(rows, ListScope{
			Name:    string(scope.Name),
			Roots:   scope.Roots,
			Default: scope.IsDefault,
		})
	}

	return rows
}
