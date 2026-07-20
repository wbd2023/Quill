package toml

import (
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

type schemaTarget struct {
	Language         string   `toml:"language"`
	Scope            string   `toml:"scope"`
	WorkingDirectory string   `toml:"working_directory"`
	FormatPaths      []string `toml:"format_paths"`
	CheckPaths       []string `toml:"check_paths"`
}

func decodeTargets(schemas map[string]schemaTarget) (targets policy.TargetConfigs) {
	targets = make(policy.TargetConfigs, 0, len(schemas))
	for _, name := range sortedMapKeys(schemas) {
		target := schemas[name]
		targets = append(targets, policy.TargetConfig{
			Name:             name,
			Language:         target.Language,
			Scope:            style.Scope(target.Scope),
			WorkingDirectory: target.WorkingDirectory,
			FormatPaths:      append([]string{}, target.FormatPaths...),
			CheckPaths:       append([]string{}, target.CheckPaths...),
		})
	}

	return targets
}

func encodeTargets(targets policy.TargetConfigs) (schemas map[string]schemaTarget) {
	if targets == nil {
		return nil
	}

	schemas = make(map[string]schemaTarget, len(targets))
	for _, target := range targets {
		schemas[target.Name] = schemaTarget{
			Language:         target.Language,
			Scope:            string(target.Scope),
			WorkingDirectory: target.WorkingDirectory,
			FormatPaths:      append([]string{}, target.FormatPaths...),
			CheckPaths:       append([]string{}, target.CheckPaths...),
		}
	}

	return schemas
}
