package architecture

func shippedPackModuleBoundaryCases() (testCases []importBoundaryCase) {
	for _, directory := range []string{
		"internal/pack/shipped/bash",
		"internal/pack/shipped/golang",
		"internal/pack/shipped/markdown",
		"internal/pack/shipped/project",
		"internal/pack/shipped/security",
		"internal/pack/shipped/text",
		"internal/pack/shipped/vocabulary",
	} {
		testCases = append(testCases, importBoundaryCase{
			name:      directory + " avoids execution orchestration",
			directory: directory,
			recursive: true,
			forbidden: []string{
				"ciphera/tools/internal/architecture",
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/runner/drivers",
				"ciphera/tools/internal/filewalk",
				"ciphera/tools/internal/installer",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/styleguide",
			},
		})
	}

	testCases = append(testCases, importBoundaryCase{
		name:      "shipped tool capabilities own Tool IDs without Pack imports",
		directory: "internal/pack/shipped/tool",
		recursive: true,
		forbidden: []string{
			"ciphera/tools/internal/architecture",
			"ciphera/tools/internal/cli",
			"ciphera/tools/internal/coverage",
			"ciphera/tools/internal/filewalk",
			"ciphera/tools/internal/installer",
			"ciphera/tools/internal/policy",
			"ciphera/tools/internal/profile",
			"ciphera/tools/internal/report",
			"ciphera/tools/internal/runner",
			"ciphera/tools/internal/runner/drivers",
			"ciphera/tools/internal/pack/shipped/bash",
			"ciphera/tools/internal/pack/shipped/golang",
			"ciphera/tools/internal/pack/shipped/markdown",
			"ciphera/tools/internal/pack/shipped/project",
			"ciphera/tools/internal/pack/shipped/security",
			"ciphera/tools/internal/pack/shipped/text",
			"ciphera/tools/internal/pack/shipped/vocabulary",
			"ciphera/tools/internal/checks",
			"ciphera/tools/internal/runtime",
			"ciphera/tools/internal/styleguide",
		},
	})

	return testCases
}
