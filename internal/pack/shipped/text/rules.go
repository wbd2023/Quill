package text

import (
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/style"
)

const (
	ruleGroupExternal style.RuleGroup = "external_tools"
	ruleGroupText     style.RuleGroup = "text_scanners"
)

/* ----------------------------------------- Rule Lists ----------------------------------------- */

func rules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{
		fileCommandRule(
			"text/spelling",
			"Spelling (non-Go)",
			tool.Misspell,
			"spelling",
			[]string{"-error", "-locale", "UK"},
		),
		lineLengthRule(),
		scannerRule(
			"text/ascii",
			"ASCII-only characters",
			ScannerASCII,
		),
		scannerRule(
			"text/exception-markers",
			"Exception marker syntax",
			ScannerExceptionMarkers,
		),
		scannerRule(
			"text/maintenance-markers",
			"TODO and FIXME marker format",
			ScannerMaintenanceMarkers,
		),
		scanRule(
			"text/section-headers",
			"Section header format",
			ruleGroupText,
			ScannerSectionHeaders,
		),
		scanRule(
			"text/section-header-density",
			"Section header density",
			ruleGroupText,
			ScannerSectionHeaderDensity,
		),
		scanRule(
			"text/section-header-names",
			"Section header naming",
			ruleGroupText,
			ScannerSectionHeaderNames,
		),
	}
}

/* ---------------------------------------- Rule Builders --------------------------------------- */

func fileCommandRule(
	id string,
	name string,
	toolID string,
	fileSet string,
	arguments []string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupExternal,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionFileCommand,
			Detail: style.FileCommandExecution{
				ToolID:    toolID,
				FileSet:   fileSet,
				Arguments: append([]string{}, arguments...),
			},
		},
	}
}

func scannerRule(
	id string,
	name string,
	scanner string,
) (rule style.RuleDefinition) {
	return scanRule(id, name, ruleGroupText, scanner)
}

func scanRule(
	id string,
	name string,
	group style.RuleGroup,
	scanner string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: group,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionRepositoryScan,
			Detail: style.RepositoryScanExecution{
				Scanner: scanner,
			},
		},
	}
}

func lineLengthRule() (rule style.RuleDefinition) {
	rule = scannerRule(
		"text/line-length",
		"Line length",
		ScannerLineLength,
	)
	execution := rule.Check.Detail.(style.RepositoryScanExecution)
	execution.FileSet = "line_length"
	rule.Check.Detail = execution
	return rule
}
