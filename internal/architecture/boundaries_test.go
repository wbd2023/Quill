package architecture

type importBoundaryCase struct {
	name           string
	directory      string
	recursive      bool
	excludeSubdirs []string
	allowed        []string
	forbidden      []string
}

func importBoundaryCases() (testCases []importBoundaryCase) {
	testCases = append(testCases, platformBoundaryCases()...)
	testCases = append(testCases, shippedPackModuleBoundaryCases()...)
	testCases = append(testCases, checkBoundaryCases()...)
	return testCases
}
