package architecture

/* ----------------------------------- Shipped Pack Boundaries ---------------------------------- */

// shippedPackDeclarationPackages lists every shipped Pack declaration package. These packages own
// Pack identity and profile-visible Rule declarations and must remain independent of the generic
// Driver facade and execution orchestration.
func shippedPackDeclarationPackages() (packages []string) {
	return []string{
		"internal/pack/shipped/bash",
		"internal/pack/shipped/golang",
		"internal/pack/shipped/markdown",
		"internal/pack/shipped/project",
		"internal/pack/shipped/security",
		"internal/pack/shipped/text",
		"internal/pack/shipped/vocabulary",
	}
}

func shippedPackModuleBoundaryCases() (testCases []importBoundaryCase) {
	// Parent Pack declaration packages must not import Drivers or execution orchestration. The
	// check stays recursive so any future sub-package other than the explicitly allowed bindings
	// child remains covered; only the bindings child is excluded because it is the single allowed
	// seam between Pack declarations and the Driver facade.
	for _, directory := range shippedPackDeclarationPackages() {
		testCases = append(testCases, importBoundaryCase{
			name:           directory + " declarations avoid execution orchestration",
			directory:      directory,
			recursive:      true,
			excludeSubdirs: []string{"bindings"},
			forbidden: []string{
				"internal/architecture",
				"internal/cli",
				"internal/coverage",
				"internal/execution/drivers",
				"internal/filewalk",
				"internal/installer",
				"internal/report",
				"internal/execution",
				"internal/process",
				"internal/workspace",
				"internal/styleguide",
			},
		})
	}

	// Each Pack's bindings child is the composition seam. It may import generic
	// Drivers, its Pack Policy, its concrete Checks, process primitives, and tools.
	for _, directory := range shippedPackDeclarationPackages() {
		testCases = append(testCases, shippedPackBindingsBoundaryCase(directory))
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

// shippedPackBindingsBoundaryCase builds the import boundary for one Pack's composition child.
// The child wires only its own declarations, policy, concrete checks, generic Drivers, process
// primitives, and shared tools. It must not reach into other Packs or application orchestration.
func shippedPackBindingsBoundaryCase(packDirectory string) (testCase importBoundaryCase) {
	forbidden := []string{
		"internal/architecture",
		"internal/cli",
		"internal/coverage",
		"internal/filewalk",
		"internal/installer",
		"internal/report",
		"internal/workspace",
		"internal/styleguide",
	}

	for _, other := range shippedPackDeclarationPackages() {
		if other == packDirectory {
			continue
		}
		forbidden = append(forbidden, other)
	}

	return importBoundaryCase{
		name:      packDirectory + "/bindings imports only its Pack runtime dependencies",
		directory: packDirectory + "/bindings",
		forbidden: forbidden,
	}
}
