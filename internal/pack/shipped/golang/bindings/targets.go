package bindings

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wbd2023/quill/internal/execution"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/style"
)

/* -------------------------------------- Target Resolution ------------------------------------- */

// goTargets resolves the named Go targets for the active scope. Targets whose scope does not
// overlap the run scope are skipped, matching how rules narrow to the relevant working set.
func goTargets(
	run execution.RunContext,
	names []string,
	goLanguage string,
) (targets []policy.TargetConfig, err error) {
	for _, name := range names {
		target, err := goTarget(run.Profile, name, goLanguage)
		if err != nil {
			return nil, err
		}

		if !run.Profile.Repository.HasScopeOverlap(run.Scope, target.Scope) {
			continue
		}

		targets = append(targets, target)
	}

	return targets, nil
}

func goTarget(
	config policy.Profile,
	name string,
	goLanguage string,
) (target policy.TargetConfig, err error) {
	target, found := config.Targets.Lookup(name)
	if !found {
		return policy.TargetConfig{}, fmt.Errorf("unknown Go target %q", name)
	}

	if target.Language != goLanguage {
		return policy.TargetConfig{}, fmt.Errorf(
			"target %q is %q, not go",
			name,
			target.Language,
		)
	}

	return target, nil
}

// targetWorkDir resolves a target's working directory against the repository root, treating an
// empty or current-directory path as the root itself.
func targetWorkDir(repoRoot string, target policy.TargetConfig) (workDir string) {
	if target.WorkingDirectory == "" || target.WorkingDirectory == "." {
		return repoRoot
	}

	return filepath.Join(repoRoot, target.WorkingDirectory)
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func decodeGoConfig(
	run execution.RunContext,
	packID string,
) (config gopolicy.Config, err error) {
	pack, found := run.Profile.PackConfigs.Lookup(packID)
	if !found {
		return gopolicy.Config{}, errMissingPackConfig(packID)
	}

	return gopolicy.DecodeConfig(pack)
}

func joinGoLocalImportPrefixes(prefixes []string) (prefix string) {
	return strings.Join(prefixes, ",")
}

// appendDiagnostics converts non-empty tool output into a diagnostic and appends it. Used by the
// Go target commands that run external tools and surface their findings as diagnostics.
func appendDiagnostics(
	diagnostics []style.Diagnostic,
	output string,
	code string,
) (result []style.Diagnostic) {
	output = strings.TrimSpace(output)
	if output == "" {
		return diagnostics
	}

	return append(diagnostics, style.Diagnostic{
		Code:    code,
		Message: output,
	})
}

func errMissingPackConfig(packID string) (err error) {
	return fmt.Errorf("packs.%s must be configured", packID)
}
