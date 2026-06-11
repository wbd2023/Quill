package architecture

/* -------------------------------------- Check Boundaries -------------------------------------- */

func checkBoundaryCases() (testCases []importBoundaryCase) {
	return []importBoundaryCase{
		{
			name:      "go check IDs stay independent",
			directory: "internal/checks/golang/check",
			forbidden: []string{"ciphera/tools/internal/"},
		},
		{
			name:      "go policy avoids Check implementations",
			directory: "internal/checks/golang/policy",
			forbidden: []string{
				"ciphera/tools/internal/style",
				"ciphera/tools/internal/pack/shipped",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/checks/golang/analysis",
				"ciphera/tools/internal/checks/golang/architecture",
				"ciphera/tools/internal/checks/golang/check",
				"ciphera/tools/internal/checks/golang/relationships",
				"ciphera/tools/internal/checks/golang/structure",
				"ciphera/tools/internal/checks/golang/syntax",
				"ciphera/tools/internal/checks/golang/test",
			},
		},
		packPolicyBoundaryCase("text"),
		packPolicyBoundaryCase("project"),
		packPolicyBoundaryCase("vocabulary"),
		{
			name:      "go syntax checks do not import Shipped Packs",
			directory: "internal/checks/golang/syntax",
			forbidden: []string{"ciphera/tools/internal/pack/shipped"},
		},
		{
			name:      "go structure checks do not import Shipped Packs",
			directory: "internal/checks/golang/structure",
			forbidden: []string{"ciphera/tools/internal/pack/shipped"},
		},
		{
			name:      "go relationship checks do not import Shipped Packs",
			directory: "internal/checks/golang/relationships",
			forbidden: []string{"ciphera/tools/internal/pack/shipped"},
		},
		{
			name:      "go test checks do not import Shipped Packs",
			directory: "internal/checks/golang/test",
			forbidden: []string{"ciphera/tools/internal/pack/shipped"},
		},
		{
			name:      "go architecture check avoids source-file checks",
			directory: "internal/checks/golang/architecture",
			forbidden: []string{
				"ciphera/tools/internal/pack/shipped",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/checks/golang/analysis",
				"ciphera/tools/internal/checks/golang/check",
				"ciphera/tools/internal/checks/golang/relationships",
				"ciphera/tools/internal/checks/golang/structure",
				"ciphera/tools/internal/checks/golang/syntax",
				"ciphera/tools/internal/checks/golang/test",
			},
		},
		{
			name:      "bash checks use filewalk directly",
			directory: "internal/checks/bash",
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/checks/text",
			},
		},
		{
			name:      "go checks do not depend on orchestration or text helpers",
			directory: "internal/checks/golang",
			recursive: true,
			forbidden: []string{
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/installer",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/pack/shipped",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/checks/text",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "text checks do not import profile",
			directory: "internal/checks/text",
			forbidden: []string{
				"ciphera/tools/internal/profile",
			},
		},
		{
			name:      "security checks do not import profile",
			directory: "internal/checks/security",
			forbidden: []string{
				"ciphera/tools/internal/profile",
			},
		},
		{
			name:      "vocabulary checks do not import profile",
			directory: "internal/checks/vocabulary",
			forbidden: []string{
				"ciphera/tools/internal/profile",
			},
		},
	}
}

func packPolicyBoundaryCase(packID string) (testCase importBoundaryCase) {
	return importBoundaryCase{
		name:      packID + " Pack Policy avoids Check implementations and orchestration",
		directory: "internal/checks/" + packID + "/policy",
		forbidden: []string{
			"ciphera/tools/internal/cli",
			"ciphera/tools/internal/coverage",
			"ciphera/tools/internal/installer",
			"ciphera/tools/internal/pack/shipped",
			"ciphera/tools/internal/profile",
			"ciphera/tools/internal/report",
			"ciphera/tools/internal/runner",
			"ciphera/tools/internal/runtime",
			"ciphera/tools/internal/style",
			"ciphera/tools/internal/styleguide",
			"ciphera/tools/internal/checks/" + packID,
		},
	}
}
