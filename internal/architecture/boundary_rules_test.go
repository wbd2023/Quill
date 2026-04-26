package architecture

func ruleBoundaryCases() (testCases []importBoundaryCase) {
	return []importBoundaryCase{
		{
			name:      "go checks do not import rulepack",
			directory: "internal/rules/golang/checks",
			forbidden: []string{"ciphera/tools/internal/rulepack"},
		},
		{
			name:      "go order checks do not import rulepack",
			directory: "internal/rules/golang/order",
			forbidden: []string{"ciphera/tools/internal/rulepack"},
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
			name:      "go checks do not depend on text helpers",
			directory: "internal/rules/golang",
			recursive: true,
			forbidden: []string{
				"ciphera/tools/internal/profile",
				"ciphera/tools/internal/rulepack",
				"ciphera/tools/internal/runtime",
				"ciphera/tools/internal/rules/text",
			},
		},
		{
			name:      "text checks do not import profile",
			directory: "internal/rules/text",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
		{
			name:      "security checks do not import profile",
			directory: "internal/rules/security",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
		{
			name:      "naming checks do not import profile",
			directory: "internal/rules/naming",
			forbidden: []string{"ciphera/tools/internal/profile"},
		},
	}
}
