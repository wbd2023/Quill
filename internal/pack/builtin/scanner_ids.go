package builtin

import (
	"ciphera/tools/internal/pack/builtin/bash"
	"ciphera/tools/internal/pack/builtin/golang"
	"ciphera/tools/internal/pack/builtin/security"
	"ciphera/tools/internal/pack/builtin/text"
	"ciphera/tools/internal/pack/builtin/vocabulary"
)

const (
	ScannerArchitecture         = golang.ScannerArchitecture
	ScannerASCII                = text.ScannerASCII
	ScannerBashMagicValues      = bash.ScannerMagicValues
	ScannerBashSafety           = bash.ScannerSafety
	ScannerBashStructure        = bash.ScannerStructure
	ScannerBashTestHygiene      = bash.ScannerTestHygiene
	ScannerExceptionMarkers     = text.ScannerExceptionMarkers
	ScannerLineLength           = text.ScannerLineLength
	ScannerMaintenanceMarkers   = text.ScannerMaintenanceMarkers
	ScannerVocabulary           = vocabulary.ScannerVocabulary
	ScannerSecrets              = security.ScannerSecrets
	ScannerSectionHeaderDensity = text.ScannerSectionHeaderDensity
	ScannerSectionHeaderNames   = text.ScannerSectionHeaderNames
	ScannerSectionHeaders       = text.ScannerSectionHeaders
)
