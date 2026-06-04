package scan

import "fmt"

var scanners = newScannerRegistry()

func newScannerRegistry() (scanners map[string]repositoryScanner) {
	scanners = make(map[string]repositoryScanner)
	addRepositoryScanners(scanners, goPackScanners())
	addRepositoryScanners(scanners, textPackScanners())
	addRepositoryScanners(scanners, bashPackScanners())
	addRepositoryScanners(scanners, securityPackScanners())
	addRepositoryScanners(scanners, vocabularyPackScanners())
	return scanners
}

func addRepositoryScanners(
	scanners map[string]repositoryScanner,
	additions map[string]repositoryScanner,
) {
	for scannerID, scanner := range additions {
		if _, exists := scanners[scannerID]; exists {
			panic(fmt.Sprintf("duplicate repository scanner %q", scannerID))
		}

		scanners[scannerID] = scanner
	}
}
