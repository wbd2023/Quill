package toolchain

import (
	"reflect"
	"strings"
	"testing"
)

func TestDetectVersionRejectsUnknownVersionSpec(t *testing.T) {
	t.Parallel()

	_, err := detectVersion(
		nil,
		unknownVersionSpec{},
		"/bin/true",
		nil,
	)
	if err == nil {
		t.Fatal("expected unknown version spec to fail")
	}

	if !strings.Contains(err.Error(), "unsupported version spec") {
		t.Fatalf("unexpected version error: %v", err)
	}
}

type unknownVersionSpec struct{}

func (unknownVersionSpec) versionSpec() {}

func TestDetectGoVersionUsesCommandRunner(t *testing.T) {
	t.Parallel()

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
		GoCommandVersion{},
		"/bin/go",
		nil,
	)
	if err != nil {
		t.Fatalf("detectVersion: %v", err)
	}

	if version != "1.24.5" {
		t.Fatalf("version = %q, want 1.24.5", version)
	}
}
