package architecture

/* -------------------------------------- Check Boundaries -------------------------------------- */

func checkBoundaryCases() (testCases []importBoundaryCase) {
	return []importBoundaryCase{
		{
			name:      "go check IDs stay independent",
			directory: "internal/checks/golang/check",
			forbidden: []string{"internal/"},
		},
		{
			name:      "go policy avoids Check implementations",
			directory: "internal/checks/gopolicy",
			forbidden: []string{
				"internal/style",
				"internal/pack/shipped",
				"internal/profile",
				"internal/execution",
				"internal/checks/golang/analysis",
				"internal/checks/golang/architecture",
				"internal/checks/golang/check",
				"internal/checks/golang/relationships",
				"internal/checks/golang/structure",
				"internal/checks/golang/syntax",
				"internal/checks/golang/test",
			},
		},
		packPolicyBoundaryCase("text"),
		packPolicyBoundaryCase("project"),
		packPolicyBoundaryCase("vocabulary"),
		{
			name:      "go syntax checks do not import Shipped Packs",
			directory: "internal/checks/golang/syntax",
			forbidden: []string{"internal/pack/shipped"},
		},
		{
			name:      "go structure checks do not import Shipped Packs",
			directory: "internal/checks/golang/structure",
			forbidden: []string{"internal/pack/shipped"},
		},
		{
			name:      "go relationship checks do not import Shipped Packs",
			directory: "internal/checks/golang/relationships",
			forbidden: []string{"internal/pack/shipped"},
		},
		{
			name:      "go test checks do not import Shipped Packs",
			directory: "internal/checks/golang/test",
			forbidden: []string{"internal/pack/shipped"},
		},
		{
			name:      "go architecture check avoids source-file checks",
			directory: "internal/checks/golang/architecture",
			forbidden: []string{
				"internal/pack/shipped",
				"internal/profile",
				"internal/checks/golang/analysis",
				"internal/checks/golang/check",
				"internal/checks/golang/relationships",
				"internal/checks/golang/structure",
				"internal/checks/golang/syntax",
				"internal/checks/golang/test",
			},
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
				"internal/pack/shipped",
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
		name:      packID + " Pack Policy avoids Check implementations and orchestration",
		directory: "internal/checks/" + packID + "policy",
		forbidden: []string{
			"internal/cli",
			"internal/coverage",
			"internal/installer",
			"internal/pack/shipped",
			"internal/profile",
			"internal/report",
			"internal/execution",
			"internal/process",
			"internal/workspace",
			"internal/style",
			"internal/styleguide",
			"internal/checks/" + packID,
		},
	}
}
