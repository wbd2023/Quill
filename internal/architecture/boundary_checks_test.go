package architecture

/* -------------------------------------- Check Boundaries -------------------------------------- */

func checkBoundaryCases() (testCases []importBoundaryCase) {
	return []importBoundaryCase{
		packPolicyBoundaryCase("golang"),
		packPolicyBoundaryCase("text"),
		packPolicyBoundaryCase("project"),
		packPolicyBoundaryCase("vocabulary"),
		{
			name:      "go syntax checks use only their own Pack policy",
			directory: "internal/checks/golang/syntax",
			forbidden: goCheckForbiddenPackImports(),
		},
		{
			name:      "go structure checks use only their own Pack policy",
			directory: "internal/checks/golang/structure",
			forbidden: goCheckForbiddenPackImports(),
		},
		{
			name:      "go relationship checks use only their own Pack policy",
			directory: "internal/checks/golang/relationships",
			forbidden: goCheckForbiddenPackImports(),
		},
		{
			name:      "go test checks use only their own Pack policy",
			directory: "internal/checks/golang/test",
			forbidden: goCheckForbiddenPackImports(),
		},
		{
			name:      "go architecture check uses only its own Pack policy",
			directory: "internal/checks/golang/architecture",
			forbidden: append(
				goCheckForbiddenPackImports(),
				"internal/checks/golang/relationships",
				"internal/checks/golang/structure",
				"internal/checks/golang/syntax",
				"internal/checks/golang/test",
			),
		},
		{
			name:      "bash checks use filewalk directly",
			directory: "internal/checks/bash",
			forbidden: []string{
				"internal/profile",
				"internal/checks/text",
			},
		},
		{
			name:      "go checks do not depend on orchestration or text helpers",
			directory: "internal/checks/golang",
			recursive: true,
			forbidden: []string{
				"internal/cli",
				"internal/coverage",
				"internal/installer",
				"internal/profile",
				"internal/pack/shipped/bash",
				"internal/pack/shipped/golang/bindings",
				"internal/pack/shipped/markdown",
				"internal/pack/shipped/project",
				"internal/pack/shipped/security",
				"internal/pack/shipped/text",
				"internal/pack/shipped/vocabulary",
				"internal/report",
				"internal/execution",
				"internal/process",
				"internal/workspace",
				"internal/checks/text",
				"internal/styleguide",
			},
		},
		{
			name:      "text checks do not import profile",
			directory: "internal/checks/text",
			forbidden: []string{
				"internal/profile",
			},
		},
		{
			name:      "security checks do not import profile",
			directory: "internal/checks/security",
			forbidden: []string{
				"internal/profile",
			},
		},
		{
			name:      "vocabulary checks do not import profile",
			directory: "internal/checks/vocabulary",
			forbidden: []string{
				"internal/profile",
			},
		},
	}
}

func packPolicyBoundaryCase(packID string) (testCase importBoundaryCase) {
	return importBoundaryCase{
		name:      packID + " Pack Policy avoids Checks and orchestration",
		directory: "internal/pack/shipped/" + packID + "/policy",
		forbidden: []string{
			"internal/architecture",
			"internal/cli",
			"internal/coverage",
			"internal/execution",
			"internal/filewalk",
			"internal/installer",
			"internal/process",
			"internal/profile",
			"internal/report",
			"internal/style",
			"internal/styleguide",
			"internal/toolchain",
			"internal/workspace",
			"internal/checks",
		},
	}
}

func goCheckForbiddenPackImports() (forbidden []string) {
	return []string{
		"internal/pack/shipped/bash",
		"internal/pack/shipped/golang/bindings",
		"internal/pack/shipped/markdown",
		"internal/pack/shipped/project",
		"internal/pack/shipped/security",
		"internal/pack/shipped/text",
		"internal/pack/shipped/vocabulary",
	}
}
