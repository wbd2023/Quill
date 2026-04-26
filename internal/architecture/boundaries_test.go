package architecture

type importBoundaryCase struct {
	name      string
	directory string
	recursive bool
	forbidden []string
}

func importBoundaryCases() (testCases []importBoundaryCase) {
	testCases = append(testCases, platformBoundaryCases()...)
	testCases = append(testCases, ruleBoundaryCases()...)
	return testCases
}
