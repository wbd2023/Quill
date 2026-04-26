package rulepack

import "ciphera/tools/internal/contract"

func shellPack() (pack Pack) {
	shfmtRule := fileCommandRule(
		"bash/shfmt",
		"Bash formatting (shfmt)",
		ToolShfmt,
		"shell",
		[]string{"-d"},
	)
	shfmtRule.FixSpec = contract.ExecutionSpec{
		Kind: ExecutorFileCommand,
		Detail: contract.FileCommandExecution{
			ToolID:    ToolShfmt,
			FileSet:   "shell",
			Arguments: []string{"-w"},
		},
	}

	return Pack{
		ID:   PackShell,
		Name: "Shell",
		Tools: selectTools(
			ToolShellcheck,
			ToolShfmt,
		),
		Rules: []RuleDefinition{
			fileCommandRule(
				"bash/shellcheck",
				"Bash static analysis (shellcheck)",
				ToolShellcheck,
				"shell",
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
