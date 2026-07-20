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
				"internal/architecture",
				"internal/cli",
				"internal/coverage",
				"internal/execution/drivers",
				"internal/filewalk",
				"internal/installer",
				"internal/profile",
				"internal/report",
				"internal/execution",
				"internal/process",
				"internal/workspace",
				"internal/styleguide",
			},
		})
	}

	testCases = append(testCases, importBoundaryCase{
		name:      "shipped tool capabilities own Tool IDs without Pack imports",
		directory: "internal/pack/shipped/tool",
		recursive: true,
		forbidden: []string{
			"internal/architecture",
			"internal/cli",
			"internal/coverage",
			"internal/filewalk",
			"internal/installer",
			"internal/policy",
			"internal/profile",
			"internal/report",
			"internal/execution",
			"internal/execution/drivers",
			"internal/pack/shipped/bash",
			"internal/pack/shipped/golang",
			"internal/pack/shipped/markdown",
			"internal/pack/shipped/project",
			"internal/pack/shipped/security",
			"internal/pack/shipped/text",
			"internal/pack/shipped/vocabulary",
			"internal/checks",
			"internal/process",
			"internal/workspace",
			"internal/styleguide",
		},
	})

	return testCases
}
