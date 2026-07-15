package toolchain

import (
	"strings"
	"testing"
)

func TestDetectVersionRejectsUnknownVersionMethod(t *testing.T) {
	t.Parallel()

	_, err := detectVersion(nil, "/bin/true", unknownVersionMethod{})
	if err == nil {
		t.Fatal("expected unknown version method to fail")
	}

	if !strings.Contains(err.Error(), "unsupported version method") {
		t.Fatalf("unexpected version error: %v", err)
	}
}

type unknownVersionMethod struct{}

func (unknownVersionMethod) versionMethod() {}

func TestParseGoVersionExtractsVersionToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		output string
		want   string
	}{
		{"go version go1.24.5 linux/amd64", "1.24.5"},
		{"go version go1.22.0 darwin/arm64", "1.22.0"},
	}

	for _, test := range tests {
		version, err := parseGoVersion(test.output)
		if err != nil {
			t.Fatalf("parseGoVersion(%q): %v", test.output, err)
		}

		if version != test.want {
			t.Errorf("parseGoVersion(%q) = %q, want %q", test.output, version, test.want)
		}
	}
}

func TestParseGoVersionRejectsUnparseableOutput(t *testing.T) {
	t.Parallel()

	_, err := parseGoVersion("not a go version string")
	if err == nil {
		t.Fatal("expected unparseable output to fail")
	}
}
