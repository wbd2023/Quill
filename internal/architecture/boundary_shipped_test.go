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
				"internal/profile",
				"internal/report",
				"internal/execution",
				"internal/process",
				"internal/workspace",
				"internal/styleguide",
			},
		})
	}

	// Each Pack's bindings child package is the only place permitted to import the Driver facade.
	// It may import only the facade, its own Pack declarations, and the shared tool catalogue; it
	// must not reach into other Packs, Driver family implementations, or any orchestration.
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

// shippedPackBindingsBoundaryCase builds the import boundary for one Pack's bindings child package.
// The child may import the drivers facade, its own Pack, and the shared tool catalogue; everything
// else (other Packs, Driver family implementations, presentation, persistence) is forbidden.
//
// Note: "internal/execution" is intentionally omitted from the forbidden list because the allowed
// facade lives at "internal/execution/drivers" and the boundary checker matches by path prefix.
// Forbidding the bare "internal/execution" prefix would also block the facade itself. The driver
// family sub-packages are enumerated individually so they remain forbidden while the facade stays
// reachable.
func shippedPackBindingsBoundaryCase(packDirectory string) (testCase importBoundaryCase) {
	forbidden := []string{
		"internal/architecture",
		"internal/cli",
		"internal/coverage",
		"internal/execution/drivers/command",
		"internal/execution/drivers/internal",
		"internal/execution/drivers/profile",
		"internal/execution/drivers/scan",
		"internal/execution/drivers/target",
		"internal/filewalk",
		"internal/installer",
		"internal/policy",
		"internal/profile",
		"internal/report",
		"internal/process",
		"internal/workspace",
		"internal/styleguide",
		"internal/checks",
	}

	for _, other := range shippedPackDeclarationPackages() {
		if other == packDirectory {
			continue
		}
		forbidden = append(forbidden, other)
	}

	return importBoundaryCase{
		name: packDirectory +
			"/bindings imports only driver facade, own Pack, and tool catalogue",
		directory: packDirectory + "/bindings",
		forbidden: forbidden,
	}
}
