package scan

import "ciphera/tools/internal/pack/builtin"

func goPackScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerArchitecture: scanGoArchitecture,
	}
}
