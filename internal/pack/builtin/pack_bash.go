package builtin

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func bashPack() (pack Pack) {
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

	return Pack{
		ID:   PackBash,
		Name: "Bash",
		Tools: selectTools(
			ToolShellcheck,
			ToolShfmt,
		),
		FileSets: bashFileSets(),
		Rules: []RuleDefinition{
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
				ScannerBashStructure,
			),
			scannerRule(
				"bash/safety",
				"Bash safety and conventions",
				ScannerBashSafety,
			),
			scannerRule(
				"bash/test-hygiene",
				"Bash test hygiene",
				ScannerBashTestHygiene,
			),
			scannerRule(
				"bash/magic-values",
				"Magic values (Bash)",
				ScannerBashMagicValues,
			),
		},
	}
}

func bashFileSets() (fileSets policy.FileSets) {
	return append(fileSets, policy.FileSetConfig{
		Name: "bash",
		Include: policy.FileSetInclude{
			Extensions: []string{".sh"},
		},
	})
}
