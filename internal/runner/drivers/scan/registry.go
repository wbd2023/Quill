package scan

import "fmt"

func repositoryScanners() (scanners map[string]repositoryScanner) {
	scanners = make(map[string]repositoryScanner)
	addRepositoryScanners(scanners, goRepositoryScanners())
	addRepositoryScanners(scanners, textRepositoryScanners())
	addRepositoryScanners(scanners, bashRepositoryScanners())
	addRepositoryScanners(scanners, securityRepositoryScanners())
	addRepositoryScanners(scanners, vocabularyRepositoryScanners())
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

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
