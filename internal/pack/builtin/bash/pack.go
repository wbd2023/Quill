package bash

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/toolchain"
)

const (
	PackID = "bash"

	ToolShellcheck = "shellcheck"
	ToolShfmt      = "shfmt"
)

const (
	ScannerMagicValues = "bash_magic_values"
	ScannerSafety      = "bash_safety"
	ScannerStructure   = "bash_structure"
	ScannerTestHygiene = "bash_test_hygiene"
)

const (
	ruleGroupExternal contract.RuleGroup = "external_tools"
	ruleGroupText     contract.RuleGroup = "text_scanners"
)

// Pack returns the Bash Shipped Pack definition.
func Pack(tools []toolchain.Capability) (definition pack.Definition) {
	return pack.Definition{
		ID:       PackID,
		Name:     "Bash",
		Tools:    append([]toolchain.Capability{}, tools...),
		FileSets: fileSets(),
		Rules:    rules(),
	}
}

/* ----------------------------------------- Rule Lists ----------------------------------------- */

func rules() (rules []contract.RuleDefinition) {
	shfmtRule := fileCommandRule(
		"bash/shfmt",
		"Bash formatting (shfmt)",
		ToolShfmt,
		"bash",
		[]string{"-d"},
	)
	shfmtRule.Fix = contract.ExecutionSpec{
		Kind: contract.ExecutorFileCommand,
		Detail: contract.FileCommandExecution{
			ToolID:    ToolShfmt,
			FileSet:   "bash",
			Arguments: []string{"-w"},
		},
	}

	return []contract.RuleDefinition{
		fileCommandRule(
			"bash/shellcheck",
			"Bash static analysis (shellcheck)",
			ToolShellcheck,
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

func fileSets() (fileSets policy.FileSets) {
	return append(fileSets, policy.FileSetConfig{
		Name: "bash",
		Include: policy.FileSetInclude{
			Extensions: []string{".sh"},
		},
	})
}

/* ---------------------------------------- Rule Builders --------------------------------------- */

func fileCommandRule(
	id string,
	name string,
	toolID string,
	fileSet string,
	arguments []string,
) (rule contract.RuleDefinition) {
	return contract.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupExternal,
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorFileCommand,
			Detail: contract.FileCommandExecution{
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
) (rule contract.RuleDefinition) {
	return contract.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupText,
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorRepositoryScan,
			Detail: contract.RepositoryScanExecution{
				Scanner: scanner,
			},
		},
	}
}
