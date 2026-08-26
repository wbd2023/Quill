package external

import (
	"fmt"
	"strings"
	"time"

	codec "github.com/BurntSushi/toml"

	"github.com/wbd2023/quill/internal/profile"
)

// SchemaVersion is the manifest schema version this implementation accepts.
const SchemaVersion = 1

// ProtocolVersion is the subprocess protocol version external Packs must speak.
const ProtocolVersion = "quill-pack-v1"

// DefaultTimeout is applied when a manifest omits its runtime timeout.
const DefaultTimeout = 30 * time.Second

// Manifest is the validated, data-only external Pack manifest. It holds only the concepts an
// external author controls; it never mirrors Quill's internal pack.Definition.
type Manifest struct {
	SchemaVersion int

	Pack    PackMeta
	Runtime RuntimeMeta
	Rules   []RuleMeta
}

// PackMeta is the Pack identity declared by the manifest.
type PackMeta struct {
	ID            string
	Name          string
	Version       string
	QuillProtocol string
}

// RuntimeMeta is the executable invocation declared by the manifest.
type RuntimeMeta struct {
	Command string
	Timeout time.Duration
}

// RuleMeta is one rule declared by the manifest.
type RuleMeta struct {
	ID          string
	Name        string
	Group       string
	Check       string
	FileSet     string
	SupportsFix bool
}

/* ------------------------------------------ Decoding ------------------------------------------ */

// tomlManifest is the TOML surface decoded straight from pack.toml. Every field an external author
// may write has an explicit tag; anything else is rejected as an unknown key.
type tomlManifest struct {
	SchemaVersion int `toml:"schema_version"`

	Pack    tomlPack    `toml:"pack"`
	Runtime tomlRuntime `toml:"runtime"`
	Rules   []tomlRule  `toml:"rules"`
}

type tomlPack struct {
	ID            string `toml:"id"`
	Name          string `toml:"name"`
	Version       string `toml:"version"`
	QuillProtocol string `toml:"quill_protocol"`
}

type tomlRuntime struct {
	Command string `toml:"command"`
	Timeout string `toml:"timeout"`
}

type tomlRule struct {
	ID          string `toml:"id"`
	Name        string `toml:"name"`
	Group       string `toml:"group"`
	Check       string `toml:"check"`
	FileSet     string `toml:"file_set"`
	SupportsFix bool   `toml:"supports_fix"`
}

// DecodeManifest strictly decodes pack.toml and validates the resulting Manifest. Unknown keys,
// unsupported versions, and invalid or missing required fields are rejected before the manifest is
// trusted.
func DecodeManifest(source string) (manifest Manifest, err error) {
	var schema tomlManifest
	metadata, err := codec.Decode(source, &schema)
	if err != nil {
		return Manifest{}, fmt.Errorf("decode pack.toml: %w", err)
	}

	if undecoded := metadata.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, len(undecoded))
		for index, key := range undecoded {
			keys[index] = key.String()
		}
		return Manifest{}, fmt.Errorf(
			"pack.toml contains unknown key(s): %s",
			strings.Join(keys, ", "),
		)
	}

	if schema.SchemaVersion != SchemaVersion {
		return Manifest{}, fmt.Errorf(
			"pack.toml schema_version %d is unsupported; expected %d",
			schema.SchemaVersion,
			SchemaVersion,
		)
	}

	if schema.Pack.QuillProtocol != ProtocolVersion {
		return Manifest{}, fmt.Errorf(
			"pack.toml quill_protocol %q is unsupported; expected %q",
			schema.Pack.QuillProtocol,
			ProtocolVersion,
		)
	}

	timeout := DefaultTimeout
	if schema.Runtime.Timeout != "" {
		timeout, err = time.ParseDuration(schema.Runtime.Timeout)
		if err != nil {
			return Manifest{}, fmt.Errorf(
				"pack.toml runtime timeout %q: %w",
				schema.Runtime.Timeout, err,
			)
		}

		if timeout <= 0 {
			return Manifest{}, fmt.Errorf(
				"pack.toml runtime timeout %q must be positive",
				schema.Runtime.Timeout,
			)
		}
	}

	rules := make([]RuleMeta, 0, len(schema.Rules))
	for _, rule := range schema.Rules {
		rules = append(rules, RuleMeta(rule))
	}

	manifest = Manifest{
		SchemaVersion: schema.SchemaVersion,
		Pack: PackMeta{
			ID:            schema.Pack.ID,
			Name:          schema.Pack.Name,
			Version:       schema.Pack.Version,
			QuillProtocol: schema.Pack.QuillProtocol,
		},
		Runtime: RuntimeMeta{
			Command: schema.Runtime.Command,
			Timeout: timeout,
		},
		Rules: rules,
	}
	if err = manifest.Validate(); err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}

// Validate reports whether a decoded manifest is internally consistent. It runs after decode so a
// caller that constructs a Manifest directly cannot bypass field checks. Validation never touches
// the filesystem; executable existence is enforced separately at source load time.
func (manifest Manifest) Validate() (err error) {
	if strings.TrimSpace(manifest.Pack.ID) == "" {
		return errManifestField("pack.id")
	}

	if manifest.Pack.ID == profile.EnabledPacksKey {
		return fmt.Errorf("pack.toml pack.id %q is reserved", manifest.Pack.ID)
	}

	if strings.TrimSpace(manifest.Pack.Name) == "" {
		return errManifestField("pack.name")
	}

	if strings.TrimSpace(manifest.Runtime.Command) == "" {
		return errManifestField("runtime.command")
	}

	seen := make(map[string]struct{}, len(manifest.Rules))
	for index, rule := range manifest.Rules {
		if strings.TrimSpace(rule.ID) == "" {
			return fmt.Errorf("pack.toml rule at index %d is missing id", index)
		}

		if strings.TrimSpace(rule.Check) == "" {
			return fmt.Errorf("pack.toml rule %q is missing check", rule.ID)
		}

		if rule.SupportsFix {
			return fmt.Errorf(
				"pack.toml rule %q declares supports_fix, but external fixes are not supported",
				rule.ID,
			)
		}

		if _, duplicate := seen[rule.ID]; duplicate {
			return fmt.Errorf("pack.toml rule %q is declared more than once", rule.ID)
		}
		seen[rule.ID] = struct{}{}
	}

	return nil
}

func errManifestField(field string) (err error) {
	return fmt.Errorf("pack.toml %s is required", field)
}
