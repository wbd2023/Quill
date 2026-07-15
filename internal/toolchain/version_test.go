package toolchain

import "testing"

func TestExtractGoToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		output string
		want   string
	}{
		{"go version go1.24.5 linux/amd64", "1.24.5"},
		{"go version go1.22.0 darwin/arm64", "1.22.0"},
	}

	for _, test := range tests {
		version, err := ExtractGoToken(test.output)
		if err != nil {
			t.Fatalf("ExtractGoToken(%q): %v", test.output, err)
		}

		if version != test.want {
			t.Errorf("ExtractGoToken(%q) = %q, want %q", test.output, version, test.want)
		}
	}
}

func TestExtractGoTokenRejectsUnparseableOutput(t *testing.T) {
	t.Parallel()

	_, err := ExtractGoToken("not a go version string")
	if err == nil {
		t.Fatal("expected unparseable output to fail")
	}
}
