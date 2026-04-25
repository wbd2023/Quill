package rulepack

import "ciphera/tools/internal/contract"

/* ----------------------------------------- Shell Pack ----------------------------------------- */

func shellPack() (pack Pack) {
	shfmtRule := fileCommandRule(
		"bash/shfmt",
		"Bash formatting (shfmt)",
		contract.ToolShfmt,
		"shell",
		[]string{"-d"},
	)
	shfmtRule.FixSpec = contract.ExecutionSpec{
		Executor:  contract.ExecutorFileCommand,
		ToolID:    contract.ToolShfmt,
		FileSet:   "shell",
		Arguments: []string{"-w"},
	}

	return Pack{
		ID:   PackShell,
		Name: "Shell",
		Tools: selectTools(
			contract.ToolShellcheck,
			contract.ToolShfmt,
		),
		Rules: []RuleDefinition{
			fileCommandRule(
				"bash/shellcheck",
				"Bash static analysis (shellcheck)",
				contract.ToolShellcheck,
				"shell",
				[]string{"-x"},
			),
			shfmtRule,
			repoScanRule(
				"bash/structure",
				"Bash script structure",
				RepositoryScannerBashStructure,
			),
			repoScanRule(
				"bash/safety",
				"Bash safety and conventions",
				RepositoryScannerBashSafety,
			),
			repoScanRule(
				"bash/test-hygiene",
				"Bash test hygiene",
				RepositoryScannerBashTestHygiene,
			),
			repoScanRule(
				"bash/magic-values",
				"Magic values (Bash)",
				RepositoryScannerBashMagicValues,
			),
		},
	}
}
