package profile

import (
	"bytes"
	"fmt"
	"strings"

	codec "github.com/BurntSushi/toml"
)

// decodeTOML decodes style profile TOML source with strict unknown-key detection.
func decodeTOML(source string) (config Profile, err error) {
	var schema schemaConfig
	metadata, err := codec.Decode(source, &schema)
	if err != nil {
		return Profile{}, err
	}

	for _, key := range metadata.Undecoded() {
		if strings.HasPrefix(key.String(), "packs.") {
			continue
		}

		return Profile{}, fmt.Errorf("unknown quill.toml key %q", key.String())
	}

	return decodeConfig(schema)
}

// encodeTOML encodes config as canonical style profile TOML.
func encodeTOML(config Profile) (contents string, err error) {
	var buffer bytes.Buffer
	encoder := codec.NewEncoder(&buffer)
	encoder.Indent = ""
	if err = encoder.Encode(encodeConfig(config)); err != nil {
		return "", err
	}

	return formatEncodedTables(buffer.String()), nil
}

func formatEncodedTables(contents string) (formatted string) {
	lines := strings.SplitAfter(contents, "\n")
	kept := make([]string, 0, len(lines))

	for index, line := range lines {
		table, ok := tableHeader(line)
		if ok && hasChildTable(table, lines[index+1:]) {
			continue
		}

		if ok &&
			len(kept) > 0 &&
			strings.TrimSpace(kept[len(kept)-1]) != "" {
			kept = append(kept, "\n")
		}

		kept = append(kept, line)
	}

	return strings.Join(kept, "")
}

func hasChildTable(parent string, lines []string) (found bool) {
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		child, ok := tableHeader(line)
		return ok && strings.HasPrefix(child, parent+".")
	}

	return false
}

func tableHeader(line string) (name string, found bool) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "[[") ||
		!strings.HasPrefix(line, "[") ||
		!strings.HasSuffix(line, "]") {
		return "", false
	}

	return strings.TrimSpace(line[1 : len(line)-1]), true
}

type schemaConfig struct {
	SchemaVersion int `toml:"schema_version"`

	Repository schemaRepositoryConfig `toml:"repository"`
	StyleGuide schemaStyleGuideConfig `toml:"style_guide"`

	PathRoles map[string][]string            `toml:"path_roles"`
	FileSets  map[string]schemaFileSetConfig `toml:"file_sets"`

	Tools   map[string]schemaPinnedTool `toml:"tools"`
	Targets map[string]schemaTarget     `toml:"targets"`

	Packs       map[string]any      `toml:"packs"`
	PackSources []schemaPackSource  `toml:"pack_sources"`
	Rules       []schemaRuleBinding `toml:"rules"`
}

func decodeConfig(schema schemaConfig) (config Profile, err error) {
	enabledPacks, err := decodeEnabledPacks(schema.Packs)
	if err != nil {
		return Profile{}, err
	}

	packPolicies, err := decodePackPolicies(schema.Packs)
	if err != nil {
		return Profile{}, err
	}

	return Profile{
		SchemaVersion: schema.SchemaVersion,

		Repository: decodeRepository(schema.Repository),
		StyleGuide: decodeStyleGuide(schema.StyleGuide),

		PathRoles: decodePathRoles(schema.PathRoles),
		FileSets:  decodeFileSets(schema.FileSets),

		Tools:   decodeTools(schema.Tools),
		Targets: decodeTargets(schema.Targets),

		EnabledPacks: enabledPacks,
		PackPolicies: packPolicies,
		PackSources:  decodePackSources(schema.PackSources),
		Rules:        decodeRules(schema.Rules),
	}, nil
}

func encodeConfig(config Profile) (schema schemaConfig) {
	return schemaConfig{
		SchemaVersion: config.SchemaVersion,

		Repository: encodeRepository(config.Repository),
		StyleGuide: encodeStyleGuide(config.StyleGuide),

		PathRoles: encodePathRoles(config.PathRoles),
		FileSets:  encodeFileSets(config.FileSets),

		Tools:   encodeTools(config.Tools),
		Targets: encodeTargets(config.Targets),

		Packs:       encodePacks(config.EnabledPacks, config.PackPolicies),
		PackSources: encodePackSources(config.PackSources),
		Rules:       encodeRules(config.Rules),
	}
}
