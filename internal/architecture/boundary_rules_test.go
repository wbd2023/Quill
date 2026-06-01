package architecture

/* --------------------------------------- Rule Boundaries -------------------------------------- */

func ruleBoundaryCases() (testCases []importBoundaryCase) {
	return []importBoundaryCase{
		{
			name:      "go check IDs stay independent",
			directory: "internal/rules/golang/check",
			forbidden: []string{"ciphera/tools/internal/"},
		},
		{
			name:      "go policy avoids rule implementations",
			directory: "internal/rules/golang/policy",
			forbidden: []string{
				"ciphera/tools/internal/contract",
				"ciphera/tools/internal/pack/builtin",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/rules/golang/analysis",
				"ciphera/tools/internal/rules/golang/architecture",
				"ciphera/tools/internal/rules/golang/check",
				"ciphera/tools/internal/rules/golang/relationships",
				"ciphera/tools/internal/rules/golang/structure",
				"ciphera/tools/internal/rules/golang/syntax",
				"ciphera/tools/internal/rules/golang/test",
			},
		},
		{
			name:      "go syntax checks do not import built-in packs",
			directory: "internal/rules/golang/syntax",
			forbidden: []string{"ciphera/tools/internal/pack/builtin"},
		},
		{
			name:      "go structure checks do not import built-in packs",
			directory: "internal/rules/golang/structure",
			forbidden: []string{"ciphera/tools/internal/pack/builtin"},
		},
		{
			name:      "go relationship checks do not import built-in packs",
			directory: "internal/rules/golang/relationships",
			forbidden: []string{"ciphera/tools/internal/pack/builtin"},
		},
		{
			name:      "go test checks do not import built-in packs",
			directory: "internal/rules/golang/test",
			forbidden: []string{"ciphera/tools/internal/pack/builtin"},
		},
		{
			name:      "go architecture rule avoids source-file checks",
			directory: "internal/rules/golang/architecture",
			forbidden: []string{
				"ciphera/tools/internal/pack/builtin",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/rules/golang/analysis",
				"ciphera/tools/internal/rules/golang/check",
				"ciphera/tools/internal/rules/golang/relationships",
				"ciphera/tools/internal/rules/golang/structure",
				"ciphera/tools/internal/rules/golang/syntax",
				"ciphera/tools/internal/rules/golang/test",
			},
		},
		{
			name:      "bash checks use filewalk directly",
			directory: "internal/rules/bash",
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/rules/text",
			},
		},
		{
			name:      "go checks do not depend on orchestration or text helpers",
			directory: "internal/rules/golang",
			recursive: true,
			forbidden: []string{
				"ciphera/tools/internal/cli",
				"ciphera/tools/internal/coverage",
				"ciphera/tools/internal/installer",
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/pack/builtin",
				"ciphera/tools/internal/report",
				"ciphera/tools/internal/runner",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/rules/text",
				"ciphera/tools/internal/styleguide",
			},
		},
		{
			name:      "text checks do not import profile",
			directory: "internal/rules/text",
			forbidden: []string{
				"ciphera/tools/internal/profile",
			},
		},
		{
			name:      "security checks do not import profile",
			directory: "internal/rules/security",
			forbidden: []string{
				"ciphera/tools/internal/profile",
			},
		},
		{
			name:      "vocabulary checks do not import profile",
			directory: "internal/rules/vocabulary",
			forbidden: []string{
				"ciphera/tools/internal/profile",
			},
		},
	}
}
