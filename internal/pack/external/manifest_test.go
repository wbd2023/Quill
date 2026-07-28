package external_test

import (
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/pack/external"
)

/* ------------------------------------------ Fixtures ------------------------------------------ */

const validManifest = `
schema_version = 1

[pack]
id = "company"
name = "Company Engineering Policy"
version = "0.1.0"
quill_protocol = "quill-pack-v1"

[runtime]
command = "bin/company-quill"
timeout = "30s"

[[rules]]
id = "company/no-direct-database-access"
name = "No direct database access"
check = "no-direct-database-access"
supports_fix = false
`

/* ------------------------------------------ Decoding ------------------------------------------ */

func TestDecodeManifestAcceptsValidManifest(t *testing.T) {
	t.Parallel()

	manifest, err := external.DecodeManifest(validManifest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err = manifest.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if manifest.Pack.ID != "company" {
		t.Fatalf("unexpected pack id: %q", manifest.Pack.ID)
	}
	if manifest.Runtime.Timeout.Seconds() != 30 {
		t.Fatalf("unexpected timeout: %v", manifest.Runtime.Timeout)
	}
	if len(manifest.Rules) != 1 || manifest.Rules[0].Check != "no-direct-database-access" {
		t.Fatalf("unexpected rules: %+v", manifest.Rules)
	}
}

func TestDecodeManifestRejectsUnsupportedSchemaVersion(t *testing.T) {
	t.Parallel()

	source := strings.Replace(validManifest, "schema_version = 1", "schema_version = 2", 1)
	if _, err := external.DecodeManifest(source); err == nil {
		t.Fatal("expected unsupported schema version error")
	}
}

func TestDecodeManifestRejectsUnsupportedProtocol(t *testing.T) {
	t.Parallel()

	source := strings.Replace(
		validManifest,
		`quill_protocol = "quill-pack-v1"`,
		`quill_protocol = "quill-pack-v99"`,
		1,
	)
	if _, err := external.DecodeManifest(source); err == nil {
		t.Fatal("expected unsupported protocol error")
	}
}

func TestDecodeManifestRejectsUnknownKey(t *testing.T) {
	t.Parallel()

	source := validManifest + "\n[unknown_section]\nkey = \"value\"\n"
	if _, err := external.DecodeManifest(source); err == nil {
		t.Fatal("expected unknown key error")
	}
}

func TestDecodeManifestRejectsMalformedTimeout(t *testing.T) {
	t.Parallel()

	source := strings.Replace(validManifest, `timeout = "30s"`, `timeout = "not-a-duration"`, 1)
	if _, err := external.DecodeManifest(source); err == nil {
		t.Fatal("expected malformed timeout error")
	}
}

func TestDecodeManifestRejectsNonPositiveTimeout(t *testing.T) {
	t.Parallel()

	source := strings.Replace(validManifest, `timeout = "30s"`, `timeout = "0s"`, 1)
	if _, err := external.DecodeManifest(source); err == nil {
		t.Fatal("expected non-positive timeout error")
	}
}

func TestDecodeManifestAppliesDefaultTimeoutWhenOmitted(t *testing.T) {
	t.Parallel()

	source := strings.Replace(validManifest, `timeout = "30s"`+"\n", "", 1)
	manifest, err := external.DecodeManifest(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manifest.Runtime.Timeout != external.DefaultTimeout {
		t.Fatalf(
			"expected default timeout %v, got %v",
			external.DefaultTimeout,
			manifest.Runtime.Timeout,
		)
	}
}

/* ----------------------------------------- Validation ----------------------------------------- */

func TestValidateRejectsMissingPackID(t *testing.T) {
	t.Parallel()

	manifest := external.Manifest{
		Pack:    external.PackMeta{Name: "x", QuillProtocol: external.ProtocolVersion},
		Runtime: external.RuntimeMeta{Command: "bin/x"},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected missing pack id error")
	}
}

func TestValidateRejectsMissingRuntimeCommand(t *testing.T) {
	t.Parallel()

	manifest := external.Manifest{
		Pack: external.PackMeta{ID: "x", Name: "x", QuillProtocol: external.ProtocolVersion},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected missing runtime command error")
	}
}

func TestValidateRejectsRuleWithoutCheck(t *testing.T) {
	t.Parallel()

	manifest := external.Manifest{
		Pack:    external.PackMeta{ID: "x", Name: "x", QuillProtocol: external.ProtocolVersion},
		Runtime: external.RuntimeMeta{Command: "bin/x"},
		Rules:   []external.RuleMeta{{ID: "x/r", Name: "r"}},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected rule missing check error")
	}
}

func TestValidateRejectsDuplicateRuleIDs(t *testing.T) {
	t.Parallel()

	manifest := external.Manifest{
		Pack:    external.PackMeta{ID: "x", Name: "x", QuillProtocol: external.ProtocolVersion},
		Runtime: external.RuntimeMeta{Command: "bin/x"},
		Rules: []external.RuleMeta{
			{ID: "x/r", Check: "c"},
			{ID: "x/r", Check: "c"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected duplicate rule id error")
	}
}

func TestValidateRejectsExternalFixSupport(t *testing.T) {
	t.Parallel()

	manifest := external.Manifest{
		Pack:    external.PackMeta{ID: "x", Name: "x", QuillProtocol: external.ProtocolVersion},
		Runtime: external.RuntimeMeta{Command: "bin/x"},
		Rules:   []external.RuleMeta{{ID: "x/r", Check: "c", SupportsFix: true}},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected supports_fix rejection; external fixes are unsupported in this MVP")
	}
}
