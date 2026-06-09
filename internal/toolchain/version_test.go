package toolchain

import (
	"reflect"
	"strings"
	"testing"

	"ciphera/tools/internal/style"
)

func TestDetectVersionRejectsUnknownVersionKind(t *testing.T) {
	_, err := detectVersion(
		nil,
		style.Tool{ID: "example"},
		Capability{
			ID:          "example",
			VersionKind: "unknown",
		},
		"/bin/true",
		nil,
	)
	if err == nil {
		t.Fatal("expected unknown version kind to fail")
	}

	if !strings.Contains(err.Error(), "unsupported version detector") {
		t.Fatalf("unexpected version error: %v", err)
	}
}

func TestDetectGoVersionUsesCommandRunner(t *testing.T) {
	runner := func(request CommandRequest) (output string, err error) {
		if request.Name != "/bin/go" {
			t.Fatalf("request.Name = %q", request.Name)
		}
		if !reflect.DeepEqual(request.Arguments, []string{"version"}) {
			t.Fatalf("request.Arguments = %#v", request.Arguments)
		}

		return "go version go1.24.5 linux/amd64", nil
	}

	version, err := detectVersion(
		runner,
		style.Tool{ID: "go", PinnedVersion: "1.24.5"},
		Capability{ID: "go", VersionKind: toolVersionGoCommand},
		"/bin/go",
		nil,
	)
	if err != nil {
		t.Fatalf("detectVersion: %v", err)
	}

	if version != "1.24.5" {
		t.Fatalf("version = %q", version)
	}
}
