package bash

import (
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/style"
)

// rules constants.
const (
	ruleGroupExternal style.RuleGroup = "external_tools"
	ruleGroupText     style.RuleGroup = "text_scanners"
)

/* ----------------------------------------- Rule Lists ----------------------------------------- */

func rules() (rules []style.RuleDefinition) {
	shfmtRule := fileCommandRule(
		"bash/shfmt",
		"Bash formatting (shfmt)",
		tool.Shfmt,
		"bash",
		[]string{"-d"},
	)
	shfmtRule.Fix = style.FileCommandExecution{
		ToolID:    tool.Shfmt,
		FileSet:   "bash",
		Arguments: []string{"-w"},
	}

	return []style.RuleDefinition{
		fileCommandRule(
			"bash/shellcheck",
			"Bash static analysis (shellcheck)",
			tool.Shellcheck,
			"bash",
			[]string{"-x"},
		),
		shfmtRule,
		scannerRule(
			"bash/structure",
			"Bash script structure",
			ScannerStructure,
		),
		scannerRule(
			"bash/safety",
			"Bash safety and conventions",
			ScannerSafety,
		),
		scannerRule(
			"bash/test-hygiene",
			"Bash test hygiene",
			ScannerTestHygiene,
		),
		scannerRule(
			"bash/magic-values",
			"Magic values (Bash)",
			ScannerMagicValues,
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
		Check: style.FileCommandExecution{
			ToolID:    toolID,
			FileSet:   fileSet,
			Arguments: append([]string{}, arguments...),
		},
	}
}

func scannerRule(
	id string,
	name string,
	scanner string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupText,
		Check: style.RepositoryScanExecution{
			Scanner: scanner,
		},
	}
}
