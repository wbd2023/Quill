package installer

import (
	"reflect"
	"testing"
)

func TestNpmInstallArgumentsIncludePackageAndVersion(t *testing.T) {
	expected := []string{
		"install",
		"--save-exact",
		"--ignore-scripts",
		"--no-audit",
		"--no-fund",
		"markdownlint-cli@0.45.0",
	}

	actual := npmArguments("markdownlint-cli", "0.45.0")
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("npm install arguments = %v, want %v", actual, expected)
	}
}
