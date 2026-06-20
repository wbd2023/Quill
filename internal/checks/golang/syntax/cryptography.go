package syntax

import (
	"fmt"
	"go/ast"
	"go/token"

	"ciphera/tools/internal/checks/golang/analysis"
)

// CheckCryptographySafety check cryptography safety.
func CheckCryptographySafety(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
	isTestFile bool,
	classifier analysis.PathClassifier,
) (violations []analysis.Violation) {
	if !classifier.HasRole(path, analysis.PathRoleGoSource) || isTestFile {
		return nil
	}

	for _, importSpec := range file.Imports {
		importPath, found := literalString(importSpec.Path)
		if !found {
			continue
		}

		switch {
		case importPath == "math/rand":
			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(importSpec.Pos()),
				Rule:     analysis.DiagnosticCryptoRand,
				Message: "production code must not import math/rand; use crypto/rand " +
					"for security-sensitive randomness",
			})

		case isDeprecatedCryptoImport(importPath):
			violations = append(violations, analysis.Violation{
				Position: fileSet.Position(importSpec.Pos()),
				Rule:     analysis.DiagnosticNoDeprecatedCrypto,
				Message: fmt.Sprintf(
					"deprecated cryptographic package %s must not be imported",
					importPath,
				),
			})
		}
	}

	return violations
}

func isDeprecatedCryptoImport(importPath string) (deprecated bool) {
	switch importPath {
	case "crypto/des", "crypto/dsa", "crypto/md5", "crypto/rc4", "crypto/sha1":
		return true
	default:
		return false
	}
}
