package executors

func repositoryScanners() (scanners map[string]repositoryScanner) {
	scanners = make(map[string]repositoryScanner)
	addRepositoryScanners(scanners, goRepositoryScanners())
	addRepositoryScanners(scanners, textRepositoryScanners())
	addRepositoryScanners(scanners, bashRepositoryScanners())
	addRepositoryScanners(scanners, securityRepositoryScanners())
	addRepositoryScanners(scanners, namingRepositoryScanners())
	return scanners
}

func addRepositoryScanners(
	scanners map[string]repositoryScanner,
	additions map[string]repositoryScanner,
) {
	for scannerID, scanner := range additions {
		scanners[scannerID] = scanner
	}
}
