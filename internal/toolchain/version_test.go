package toolchain

import (
	"reflect"
	"strings"
	"testing"
)

func TestDetectVersionRejectsUnknownVersionKind(t *testing.T) {
	_, err := detectVersion(
		nil,
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
		Capability{ID: "go", VersionKind: VersionKindGoCommand},
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
