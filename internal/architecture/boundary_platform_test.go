package architecture

func platformBoundaryCases() (testCases []importBoundaryCase) {
	return []importBoundaryCase{
		{
			name:      "contract does not import internal packages",
			directory: "internal/contract",
			forbidden: []string{"ciphera/tools/internal/"},
		},
		{
			name:      "profile depends only on contracts and policy",
			directory: "internal/profile",
			forbidden: []string{
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
			},
		},
		{
			name:      "policy depends only on contracts",
			directory: "internal/policy",
			forbidden: []string{
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/filewalk",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "runner is generic execution machinery",
			directory: "internal/runner",
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runtime",
			},
		},
		{
			name:      "toolchain depends only on contracts",
			directory: "internal/toolchain",
			forbidden: []string{
				"ciphera/tools/internal/architecture",
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/filewalk",
				"ciphera/tools/internal/policy",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "rulepack depends only on contracts and toolchain",
			directory: "internal/rulepack",
			forbidden: []string{
				"ciphera/tools/internal/architecture",
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/executors",
				"ciphera/tools/internal/filewalk",
				"ciphera/tools/internal/policy",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "filewalk does not import profile",
			directory: "internal/filewalk",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
		{
			name:      "styleguide does not import profile or rulepack",
			directory: "internal/styleguide",
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/rulepack",
			},
		},
	}
}
